package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/serhatyavuzyigit/magazi/handlers"
)

var (
	intervalTime = flag.Int("time", 1, "interval time as minutes")
	listenPort   = flag.Int("port", 9242, "port to listen")
	fileName     = flag.String("file", "backup-data", "name of the backup file")
)

func main() {
	// parse the flag variables
	flag.Parse()

	l := log.New(os.Stdout, "magazi-api ", log.LstdFlags)
	bindAddress := fmt.Sprintf("%s%d", ":", *listenPort)

	// create the magazi handler
	mh := handlers.NewMagazi(l, *fileName)
	// create the serve mux and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", mh)

	// create a new server
	s := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// define time ticker for interval job
	// create goroutine for interval
	ticker := time.NewTicker(time.Duration(*intervalTime) * time.Second)
	intervalDone := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// update file with current store
				mh.UpdateFile()
			case <-intervalDone:
				ticker.Stop()
				return
			}
		}
	}()

	// create goroutine for the server
	go func() {
		l.Printf("Starting server on port %d", *listenPort)

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)
	close(intervalDone)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
