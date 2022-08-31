/**
 * @author Jose Nidhin
 */
package main

import (
	"bufio"
	"io"

	"github.com/moov-io/iso8583"
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
	rwc io.ReadWriteCloser

	spec *iso8583.MessageSpec

	msgLenReader MessageLengthReader

	msgLenWriter MessageLengthWriter

	inMsgHandler InboundMessageHandler

	reqCh chan iso8583.Message
}

func NewConnection(rwc io.ReadWriteCloser, spec *iso8583.MessageSpec, mlReader MessageLengthReader, mlWriter MessageLengthWriter, inMsgHandler InboundMessageHandler) (*Connection, error) {
	conn := &Connection{
		rwc:          rwc,
		spec:         spec,
		msgLenReader: mlReader,
		msgLenWriter: mlWriter,
		inMsgHandler: inMsgHandler,
	}

	conn.run()

	return conn, nil
}

func (conn *Connection) run() {
	go conn.readLoop()
}

func (conn *Connection) readLoop() {
	var err error
	var msgLen int

	reader := bufio.NewReader(conn.rwc)

	for {
		msgLen, err = conn.msgLenReader(reader)
		if err != nil {
			logger.Fatalf("read loop: reading msg len failed - %v", err)
			break
		}

		rawMessage := make([]byte, msgLen)
		_, err = io.ReadFull(reader, rawMessage)
		if err != nil {
			logger.Fatalf("read loop: reading full msg failed - %v", err)
			break
		}

		logger.Printf("read loop: raw message - %s", string(rawMessage))
	}
}
