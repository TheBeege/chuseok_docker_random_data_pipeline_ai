package main

import (
	"net/http"
	"time"
	"os"
	"io/ioutil"
	"os/signal"
	"syscall"
	"sync"
	"strconv"
	"log"
)

type ResponseData struct {
	Timestamp time.Time
	Value string
}

// http://blog.narenarya.in/concurrent-http-in-go.html
func MakeRequest(url string, dataChannel chan <- *ResponseData, interruptChannel <- chan bool, waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()

	for {
		log.Println("Getting from interruptChannel")
		select {
			case _, ok := <-interruptChannel:
				if !ok {
					log.Print("MakeRequest found interruptChannel closed. Breaking")
					interruptChannel = nil
				}
			default:
		}
		if interruptChannel == nil {
			break
		}
		//log.Println("Passed interruptChannel")
		//log.Println("Getting from url")
		response, err := http.Get(url)
		//log.Println("Finished getting from url")
		if err != nil {
			log.Printf("failed getting url: %s", err)
			continue
		}

		requestValue, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("unable to read response body: %s", err)
			continue
		}
		data := &ResponseData{
			Timestamp: time.Now(),
			Value:     string(requestValue),
		}
		response.Body.Close()
		log.Println("Sending data to dataChannel")
		dataChannel <- data
	}
	log.Println("Finished MakeRequest loop")
}

func HandleData(dataChannel <- chan *ResponseData, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for data := range dataChannel {
		log.Println(data)
	}
	log.Println("HandleData finished ranging over dataChannel")
}

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	<-sigs
	log.Println("Received interrupt...")
}

func main() {
	collectorWaitGroup := &sync.WaitGroup{}
	producerWaitGroup := &sync.WaitGroup{}
	if len(os.Args) != 4 {
		log.Println("This script expects three args in order: collectionUrl, countCollectors, countQueueProducers")
		os.Exit(1)
	}
	url := string(os.Args[1])
	countCollectors, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Printf("countCollectors received an invalid value '%s'. Exiting", os.Args[2])
		os.Exit(1)
	}
	countQueueProducers, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Printf("countQueueProducers received an invalid value '%s'. Exiting", os.Args[3])
		os.Exit(1)
	}
	dataChannel := make(chan *ResponseData)
	interruptChannel := make(chan bool, 1)
	for i := 0; i < countCollectors; i++ {
		collectorWaitGroup.Add(1)
		go MakeRequest(url, dataChannel, interruptChannel, collectorWaitGroup)
	}
	for i := 0; i < countQueueProducers; i++ {
		producerWaitGroup.Add(1)
		go HandleData(dataChannel, producerWaitGroup)
	}

	waitForSignal()
	log.Println("closing interrupt channel")
	close(interruptChannel)
	log.Println("Waiting for collectorWaitGroup")
	collectorWaitGroup.Wait()
	log.Println("Finished waiting for collectorWaitGroup")

	log.Println("Closing data channel")
	close(dataChannel)
	log.Println("Waiting for producerWaitGroup")
	producerWaitGroup.Wait()
	log.Println("Finished waiting for producerWaitGroup")

}