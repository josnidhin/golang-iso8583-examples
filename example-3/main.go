/**
 * @author Jose Nidhin
 */
package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/moov-io/iso8583"
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
	wg := &sync.WaitGroup{}
	address := ":8080"
	shutdownNotifier := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		signalHandler(shutdownNotifier)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startServer(shutdownNotifier, wg, address)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startClient(shutdownNotifier, address)
	}()

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

func startServer(shutdownNotifier <-chan struct{}, wg *sync.WaitGroup, address string) {
	fnName := "main.startServer"
	network := "tcp"

	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		logger.Printf("%s: address resolve failed - %v", fnName, err)
		return
	}

	tcpListener, err := net.ListenTCP(network, tcpAddr)
	if err != nil {
		logger.Printf("%s: bind listener failed - %v", fnName, err)
		return
	}

	go func() {
		<-shutdownNotifier
		tcpListener.Close()
	}()
	logger.Printf("%s: server listening on address - %s", fnName, tcpAddr)

	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			select {
			case <-shutdownNotifier:
				return
			default:
				logger.Printf("%s: accept connection failed - %v", fnName, err)
				break
			}
		}

		logger.Printf("%s: new connection", fnName)

		connHandler, err := NewConnectionHandler(conn, Spec1HeaderSize, Spec1, MsgLenReader, MsgLenWriter, inMsgHandler)
		connHandler.Start()
		if err != nil {
			logger.Fatalf("%s: error creating connection handler - %v", fnName, err)
		}

		wg.Add(1)
		go func(connHandler *ConnectionHandler) {
			defer wg.Done()
			<-shutdownNotifier

			err := connHandler.Close()
			if err != nil {
				logger.Printf("%s: error closing connection handler - %v", fnName, err)
			}

			connHandler.Done()
		}(connHandler)
	}
}

func inMsgHandler(msg *iso8583.Message) {
	printISOMsg(msg)
}

func printISOMsg(msg *iso8583.Message) {
	tw := tabwriter.NewWriter(os.Stdout, 2, 2, 1, ' ', 0)

	for pos := 0; pos < 128; pos++ {
		value, err := msg.GetString(pos)

		if err != nil {
			continue
		}

		if value == "" {
			continue
		}

		field := msg.GetField(pos)

		fmt.Fprintf(tw, "%3d\t%s\t%s\n", pos, field.Spec().Description, value)
	}
	tw.Flush()

	fmt.Println()
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
