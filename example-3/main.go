/**
 * @author Jose Nidhin
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

const (
	clientMode = "client"
	serverMode = "server"
)

func main() {
	var address, mode string
	flag.StringVar(&address, "address", ":8080", "set the server address")
	flag.StringVar(&mode, "mode", serverMode, "choose the running mode eg: server, client")
	flag.Parse()

	mode = strings.ToLower(mode)

	wg := &sync.WaitGroup{}
	shutdownNotifier := make(chan struct{})

	go func() {
		signalHandler(shutdownNotifier)
	}()

	var err error
	var server *Server
	var client *Client

	go func() {
		<-shutdownNotifier

		if server != nil {
			server.Shutdown()
			wg.Done()
		}

		if client != nil {
			client.Shutdown()
			wg.Done()
		}
	}()

	switch mode {
	case serverMode:
		server, err = NewServer(address)
		if err != nil {
			logger.Fatalf("%v", err)
		}

		wg.Add(1)
		server.Start()

	case clientMode:
		client, err = NewClient(address)
		if err != nil {
			logger.Fatalf("%v", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			client.Start()
		}()

	default:
		fmt.Printf("Unkown mode - %s\n", mode)
		os.Exit(1)
	}

	wg.Wait()
}

func signalHandler(shutdownNotifier chan struct{}) {
	fnName := "main.signalHandler"

	defer func() {
		logger.Printf("%s: graceful shutdown initialised", fnName)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-shutdownNotifier:
		return
	case <-sigChan:
		close(shutdownNotifier)
	}
}
