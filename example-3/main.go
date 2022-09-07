/**
 * @author Jose Nidhin
 */
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	clientMode = "client"
	serverMode = "server"
)

const hexStr = "009d49534f303234303030303535303230303532333838303030303841303830303031303831313030393934313830303030303030313030303030333133313032383432343838373539313032363431303331333033313330303030303034303139393131304d4f4e353047415a4f582020204e456469736f6e203132333520202020202020202020204d6f6e746572726579202020204e4c204d58343834"

var sampleRawInput []byte

func init() {
	var err error
	fnName := "main.init"

	sampleRawInput, err = hex.DecodeString(hexStr)
	if err != nil {
		logger.Panicf("%s: raw input creation failed - %v", fnName, err)
	}
}

func main() {
	var mode string
	flag.StringVar(&mode, "mode", serverMode, "choose the running mode eg: server, client")
	flag.Parse()

	mode = strings.ToLower(mode)

	wg := &sync.WaitGroup{}
	address := ":8080"
	shutdownNotifier := make(chan struct{})

	go func() {
		signalHandler(shutdownNotifier)
	}()

	var err error
	var server *Server

	go func() {
		<-shutdownNotifier

		if server != nil {
			server.Shutdown()
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
		wg.Add(1)
		go func() {
			defer wg.Done()
			startClient(shutdownNotifier, address)
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

func startClient(shutdownNotifier <-chan struct{}, address string) {
	var tcpConn *net.TCPConn
	fnName := "main.startClient"
	network := "tcp"

	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		logger.Printf("%s: address resolve failed - %v", fnName, err)
		return
	}

	ticker := time.NewTicker(1 * time.Second)

loop:
	for {
		select {
		case <-shutdownNotifier:
			break loop
		case <-ticker.C:
			tcpConn, err = net.DialTCP(network, nil, tcpAddr)
			if err != nil {
				logger.Printf("%s: dial failed - %v", fnName, err)
				break
			}

			break loop
		}
	}

	defer tcpConn.Close()
	defer ticker.Stop()

	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(60 * time.Second)

	for {
		select {
		case <-shutdownNotifier:
			return
		case <-ticker.C:
			_, err = tcpConn.Write(sampleRawInput)
			if err != nil {
				logger.Printf("%s: error while writing to tcp connection - %v", fnName, err)
			}
		}
	}
}
