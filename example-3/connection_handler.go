/**
 * @author Jose Nidhin
 */
package main

import (
	"bufio"
	"io"
	"sync"

	"github.com/moov-io/iso8583"
	"github.com/pkg/errors"
)

// MessageLengthReader reads message header from the provided reader interface
// and returns message length
type MessageLengthReader func(r io.Reader) (int, error)

// MessageLengthWriter writes message header with encoded length into the
// provided writer interface
type MessageLengthWriter func(w io.Writer, length int) (int, error)

// InboundMessageHandler will be called whenever a message is received
type InboundMessageHandler func(*iso8583.Message)

type ConnectionHandler struct {
	rwc              io.ReadWriteCloser
	headerSize       int
	spec             *iso8583.MessageSpec
	msgLenReader     MessageLengthReader
	msgLenWriter     MessageLengthWriter
	inMsgHandler     InboundMessageHandler
	shutdownNotifier chan struct{}
	reqCh            chan []byte
	wg               *sync.WaitGroup
	isClosingMutex   sync.Mutex
	isClosing        bool
}

func NewConnectionHandler(rwc io.ReadWriteCloser,
	headerSize int,
	spec *iso8583.MessageSpec,
	mlReader MessageLengthReader,
	mlWriter MessageLengthWriter,
	inMsgHandler InboundMessageHandler) (*ConnectionHandler, error) {
	ch := &ConnectionHandler{
		rwc:              rwc,
		headerSize:       headerSize,
		spec:             spec,
		msgLenReader:     mlReader,
		msgLenWriter:     mlWriter,
		inMsgHandler:     inMsgHandler,
		shutdownNotifier: make(chan struct{}),
		reqCh:            make(chan []byte),
		wg:               &sync.WaitGroup{},
	}

	return ch, nil
}

func (ch *ConnectionHandler) Start() {
	ch.run()
}

func (ch *ConnectionHandler) Close() error {
	ch.isClosingMutex.Lock()
	if ch.isClosing {
		ch.isClosingMutex.Unlock()
		return nil
	}

	ch.isClosing = true
	ch.isClosingMutex.Unlock()

	close(ch.shutdownNotifier)

	ch.wg.Wait()

	err := ch.rwc.Close()
	if err != nil {
		return errors.Wrap(err, "connection close error")
	}

	return nil
}

func (ch *ConnectionHandler) Done() {
	ch.wg.Wait()
	return
}

func (ch *ConnectionHandler) run() {
	go ch.readLoop()
	go ch.requestListener()
}

// readLoop reads the data from the connection and sends it on the request
// channel for further processing
func (ch *ConnectionHandler) readLoop() {
	var err error
	var msgLen int
	fnName := "ConnectionHandler.readLoop"

	reader := bufio.NewReader(ch.rwc)

	for {
		msgLen, err = ch.msgLenReader(reader)
		if err != nil {
			logger.Printf("%s: reading msg len failed - %v", fnName, err)
			break
		}

		rawMsg := make([]byte, msgLen)
		_, err = io.ReadFull(reader, rawMsg)
		if err != nil {
			logger.Printf("%s: reading full msg failed - %v", fnName, err)
			break
		}

		logger.Printf("%s: raw message - %s", fnName, string(rawMsg))

		ch.reqCh <- rawMsg
	}

	ch.handleConnectionError(err)
}

// requestListener reads the data from the request channel and invokes the
// requestHandler in a goroutine
func (ch *ConnectionHandler) requestListener() {
	fnName := "ConnectionHandler.requestListener"
	for {
		select {
		case rawMsg := <-ch.reqCh:
			go ch.requestHandler(rawMsg)
		case <-ch.shutdownNotifier:
			logger.Printf("%s: shutdown initialized", fnName)
			return
		}
	}
}

func (ch *ConnectionHandler) requestHandler(rawMsg []byte) {
	msg := iso8583.NewMessage(ch.spec)
	msg.Unpack(rawMsg[ch.headerSize:])
	ch.inMsgHandler(msg)
}

func (ch *ConnectionHandler) handleConnectionError(err error) {
	ch.isClosingMutex.Lock()
	if err == nil || ch.isClosing {
		ch.isClosingMutex.Unlock()
		return
	}

	ch.isClosingMutex.Unlock()

	ch.Close()
}
