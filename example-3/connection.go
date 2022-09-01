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

type Connection struct {
	rwc              io.ReadWriteCloser
	headerSize       int
	spec             *iso8583.MessageSpec
	msgLenReader     MessageLengthReader
	msgLenWriter     MessageLengthWriter
	inMsgHandler     InboundMessageHandler
	shutdownNotifier chan int
	reqCh            chan []byte
	wg               *sync.WaitGroup
}

func NewConnection(rwc io.ReadWriteCloser,
	headerSize int,
	spec *iso8583.MessageSpec,
	mlReader MessageLengthReader,
	mlWriter MessageLengthWriter,
	inMsgHandler InboundMessageHandler) (*Connection, error) {
	conn := &Connection{
		rwc:              rwc,
		headerSize:       headerSize,
		spec:             spec,
		msgLenReader:     mlReader,
		msgLenWriter:     mlWriter,
		inMsgHandler:     inMsgHandler,
		shutdownNotifier: make(chan int),
		reqCh:            make(chan []byte),
		wg:               &sync.WaitGroup{},
	}

	conn.run()

	return conn, nil
}

func (conn *Connection) Close() error {
	close(conn.shutdownNotifier)

	conn.wg.Wait()

	err := conn.rwc.Close()
	if err != nil {
		return errors.Wrap(err, "connection close error")
	}

	return nil
}

func (conn *Connection) Done() {
	conn.wg.Wait()
	return
}

func (conn *Connection) run() {
	go conn.readLoop()
	go conn.requestListener()
}

func (conn *Connection) readLoop() {
	var err error
	var msgLen int
	fnName := "Connection.readLoop"

	conn.wg.Add(1)
	defer conn.wg.Done()

	reader := bufio.NewReader(conn.rwc)

	for {
		select {
		case <-conn.shutdownNotifier:
			return
		default:
			msgLen, err = conn.msgLenReader(reader)
			if err != nil {
				//logger.Printf("%s: reading msg len failed - %v", fnName, err)
				break
			}

			rawMsg := make([]byte, msgLen)
			_, err = io.ReadFull(reader, rawMsg)
			if err != nil {
				logger.Printf("%s: reading full msg failed - %v", fnName, err)
				break
			}

			logger.Printf("%s: raw message - %s", fnName, string(rawMsg))

			conn.reqCh <- rawMsg
		}
	}
}

func (conn *Connection) requestListener() {
	fnName := "Connection.requestListener"
	for {
		select {
		case rawMsg := <-conn.reqCh:
			go conn.requestHandler(rawMsg)
		case <-conn.shutdownNotifier:
			logger.Printf("%s: shutdown initialized", fnName)
			return
		}
	}
}

func (conn *Connection) requestHandler(rawMsg []byte) {
	msg := iso8583.NewMessage(conn.spec)
	msg.Unpack(rawMsg[conn.headerSize:])
	conn.inMsgHandler(msg)
}
