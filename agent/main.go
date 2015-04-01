package main

import (
	//	"fmt"
	//	"github.com/nu7hatch/gouuid"
	"flag"
	"fmt"
	//	"github.com/gonum/matrix/mat64"
	"github.com/owulveryck/flue"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	//	"os"
	//	    "syscall"
	//	"os/signal"
)

func cleanup() {
	log.Println("cleanup")
}

func main() {

	/*
		// Catching the interrupt
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
		    <-c
		    cleanup()
		    os.Exit(1)
		}()
	*/
	//Parsing the dot
	dotFile := flag.String("dot", "", "The dot file")

	flag.Parse()
	if *dotFile == "" {
		log.Println("Server mode")
		proto := "tcp"
		socket := "localhost:4546"
		flue.Rserver(&proto, &socket)
	} else {
		log.Println("Client mode")

		var topologyDot []byte
		topologyDot, err := ioutil.ReadFile(*dotFile)
		if err != nil {
			log.Panic("Cannot read file")
		}

		log.Println("Parsing...")
		taskStructure := flue.ParseTasks(topologyDot)
		// How many tasks

		var wg sync.WaitGroup
		taskStructureChan := make(chan *flue.TaskGraphStructure)
		doneChan := make(chan int)
		// DEBUG
		fmt.Println("Ajacency Matrix")
		flue.PrintAdjacencyMatrix(taskStructure)
		fmt.Println("Degree Matrix")
		flue.PrintDegreeMatrix(taskStructure)

		// END DEBUG
		for taskIndex, _ := range taskStructure.Tasks {
			log.Printf("taskIndex: %v", taskIndex)
			go flue.Runner(taskStructure, taskIndex, taskStructureChan, doneChan, &wg)
			wg.Add(1)
		}
		go flue.Advertize(taskStructureChan, doneChan)
		router := flue.NewRouter(*taskStructure)

		go log.Fatal(http.ListenAndServe(":8080", router))
		wg.Wait()
	}
}
