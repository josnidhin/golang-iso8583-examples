/**
 * @author Jose Nidhin
 */
package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/moov-io/iso8583"
)

type MockReadWriteCloser struct {
	bytes.Buffer
}

func (rwc *MockReadWriteCloser) Close() error {
	return nil
}

var rwc *MockReadWriteCloser
var conn *Connection

func main() {
	var err error
	fnName := "main"
	wg := &sync.WaitGroup{}
	shutdownNotifier := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()
		signalHandler(shutdownNotifier)
	}()

	rwc = &MockReadWriteCloser{}

	inMsgHandler := func(msg *iso8583.Message) {
		printISOMsg(msg)
	}

	conn, err = NewConnection(rwc, Spec1HeaderSize, Spec1, MsgLenReader, MsgLenWriter, inMsgHandler)
	if err != nil {
		logger.Fatalf("%s: error creating connection - %v", fnName, err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		mockWriter(shutdownNotifier)
	}()

	conn.Done()
	wg.Wait()
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

func signalHandler(shutdownNotifier chan int) {
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
		conn.Close()
		close(shutdownNotifier)
	}
}

func mockWriter(shutdownNotifier <-chan int) {
	fnName := "main.mockWriter"
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	hexStr := "009d49534f303234303030303535303230303532333838303030303841303830303031303831313030393934313830303030303030313030303030333133313032383432343838373539313032363431303331333033313330303030303034303139393131304d4f4e353047415a4f582020204e456469736f6e203132333520202020202020202020204d6f6e746572726579202020204e4c204d58343834"
	rawInput, err := hex.DecodeString(hexStr)
	if err != nil {
		logger.Fatalf("%s: raw input creation error - %v", fnName, err)
	}

	for {
		select {
		case <-shutdownNotifier:
			return
		case <-ticker.C:
			_, err = rwc.Write(rawInput)
			if err != nil {
				logger.Fatalf("%s: error while writing rawInput into rwc - %v", fnName, err)
			}
		}
	}
}
