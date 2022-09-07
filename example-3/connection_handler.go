/**
 * @author Jose Nidhin
 */
package main

import (
	"bufio"
	"bytes"
	"io"
	"sync"

	"github.com/moov-io/iso8583"
	"github.com/pkg/errors"
)

var (
	ClosedError = errors.New("connection handler closed")
)

// MessageLengthReader reads message header from the provided reader interface
// and returns message length
type MessageLengthReader func(r io.Reader) (int, error)

// MessageLengthWriter writes message header with encoded length into the
// provided writer interface
type MessageLengthWriter func(w io.Writer, length int) (int, error)

type ConnectionHandler struct {
	rwc              io.ReadWriteCloser
	headerSize       int
	spec             *iso8583.MessageSpec
	msgLenReader     MessageLengthReader
	msgLenWriter     MessageLengthWriter
	shutdownNotifier chan struct{}
	reqCh            chan []byte
	reqMsgCh         chan<- *iso8583.Message
	resMsgCh         <-chan *iso8583.Message
	wg               *sync.WaitGroup
	isClosingMutex   sync.Mutex
	isClosing        bool
}

func NewConnectionHandler(rwc io.ReadWriteCloser,
	headerSize int,
	spec *iso8583.MessageSpec,
	mlReader MessageLengthReader,
	mlWriter MessageLengthWriter,
	reqMsgCh chan<- *iso8583.Message,
	resMsgCh <-chan *iso8583.Message) (*ConnectionHandler, error) {
	ch := &ConnectionHandler{
		rwc:              rwc,
		headerSize:       headerSize,
		spec:             spec,
		msgLenReader:     mlReader,
		msgLenWriter:     mlWriter,
		reqMsgCh:         reqMsgCh,
		resMsgCh:         resMsgCh,
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
	go ch.sendLoop()
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
	var rawMsg []byte
	fnName := "ConnectionHandler.requestListener"

	for {
		select {
		case rawMsg = <-ch.reqCh:
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
	ch.reqMsgCh <- msg
}

func (ch *ConnectionHandler) sendLoop() {
	var msg *iso8583.Message
	fnName := "ConnectionHandler.sendLoop"

	for {
		select {
		case msg = <-ch.resMsgCh:
			ch.sendHandler(msg)
		case <-ch.shutdownNotifier:
			logger.Printf("%s: shutdown initialized", fnName)
			return
		}
	}
}

func (ch *ConnectionHandler) sendHandler(msg *iso8583.Message) {
	fnName := "ConnectionHandler.sendHandler"

	ch.wg.Add(1)
	defer ch.wg.Done()

	ch.isClosingMutex.Lock()
	if ch.isClosing {
		ch.isClosingMutex.Unlock()
		logger.Printf("%s: connction handler is closing", fnName)
		return
	}

	ch.isClosingMutex.Unlock()

	packed, err := msg.Pack()
	if err != nil {
		logger.Printf("%s: packing iso8583 message failed - %v", fnName, err)
		return
	}

	var buf bytes.Buffer
	_, err = ch.msgLenWriter(&buf, len(packed))
	if err != nil {
		logger.Printf("%s: writing msg header to buffer failed - %v", fnName, err)
		return
	}

	_, err = buf.Write(packed)
	if err != nil {
		logger.Printf("%s: writing packed msg to buffer failed - %v", fnName, err)
		return
	}

	_, err = ch.rwc.Write(buf.Bytes())
	if err != nil {
		logger.Printf("%s: writing message to connection failed - %v", fnName, err)
		return
	}
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
