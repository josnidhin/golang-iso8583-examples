/**
 * @author Jose Nidhin
 */
package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/moov-io/iso8583"
	"github.com/pkg/errors"
)

var (
	ClosedError = errors.New("connection handler closed")
)

var (
	connTimeout     = 30 * time.Second
	connReadTimeout = 5 * time.Second
)

// MessageLengthReader reads message header from the provided reader interface
// and returns message length
type MessageLengthReader func(r io.Reader) (int, error)

// MessageLengthWriter writes message header with encoded length into the
// provided writer interface
type MessageLengthWriter func(w io.Writer, length int) (int, error)

type ConnectionHandler struct {
	id                    uuid.UUID
	conn                  net.Conn
	headerSize            int
	spec                  *iso8583.MessageSpec
	msgLenReader          MessageLengthReader
	msgLenWriter          MessageLengthWriter
	deadlineExceededCount int
	shutdownNotifier      chan struct{}
	reqCh                 chan []byte
	reqMsgCh              chan<- *iso8583.Message
	resMsgCh              <-chan *iso8583.Message
	wg                    *sync.WaitGroup
	isClosingMutex        sync.Mutex
	isClosing             bool
}

func NewConnectionHandler(conn net.Conn,
	headerSize int,
	spec *iso8583.MessageSpec,
	mlReader MessageLengthReader,
	mlWriter MessageLengthWriter,
	reqMsgCh chan<- *iso8583.Message,
	resMsgCh <-chan *iso8583.Message) (*ConnectionHandler, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to unique connection id")
	}

	ch := &ConnectionHandler{
		id:               id,
		conn:             conn,
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

	err := ch.conn.Close()
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

	ch.wg.Add(1)
	defer ch.wg.Done()

	reader := bufio.NewReader(ch.conn)

loop:
	for {
		select {
		case <-ch.shutdownNotifier:
			logger.Printf("%s (%s): shutdown initialized", fnName, ch.id.String())
			close(ch.reqCh)
			break loop
		default:
			ch.conn.SetReadDeadline(time.Now().Add(connReadTimeout))

			msgLen, err = ch.msgLenReader(reader)
			if err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					ch.deadlineExceededCount++
					elapsed := time.Duration(ch.deadlineExceededCount) * connReadTimeout

					if connTimeout < elapsed {
						logger.Printf("%s (%s): connection timeout exceeded", fnName, ch.id.String())
						break loop
					}

					logger.Printf("%s (%s): read dead line exceeded", fnName, ch.id.String())
					continue loop
				}

				logger.Printf("%s (%s): reading msg len failed - %v", fnName, ch.id.String(), err)
				break loop
			}

			ch.deadlineExceededCount = 0

			rawMsg := make([]byte, msgLen)
			_, err = io.ReadFull(reader, rawMsg)
			if err != nil {
				logger.Printf("%s (%s): reading full msg failed - %v", fnName, ch.id.String(), err)
				break loop
			}

			logger.Printf("%s (%s): raw message - %s", fnName, ch.id.String(), string(rawMsg))

			ch.reqCh <- rawMsg
		}
	}

	ch.handleConnectionError(err)
}

// requestListener reads the data from the request channel and invokes the
// requestHandler in a goroutine
func (ch *ConnectionHandler) requestListener() {
	var rawMsg []byte

	for {
		rawMsg = <-ch.reqCh
		go ch.requestHandler(rawMsg)
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
			logger.Printf("%s (%s): shutdown initialized", fnName, ch.id.String())
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
		logger.Printf("%s (%s): connction handler is closing", fnName, ch.id.String())
		return
	}

	ch.isClosingMutex.Unlock()

	packed, err := msg.Pack()
	if err != nil {
		logger.Printf("%s (%s): packing iso8583 message failed - %v", fnName, ch.id.String(), err)
		return
	}

	var buf bytes.Buffer
	_, err = ch.msgLenWriter(&buf, len(packed))
	if err != nil {
		logger.Printf("%s (%s): writing msg header to buffer failed - %v", fnName, ch.id.String(), err)
		return
	}

	_, err = buf.Write(packed)
	if err != nil {
		logger.Printf("%s (%s): writing packed msg to buffer failed - %v", fnName, ch.id.String(), err)
		return
	}

	_, err = ch.conn.Write(buf.Bytes())
	if err != nil {
		logger.Printf("%s (%s): writing message to connection failed - %v", fnName, ch.id.String(), err)
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

	go ch.Close()
}
