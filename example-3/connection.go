/**
 * @author Jose Nidhin
 */
package main

import (
	"io"

	"github.com/moov-io/iso8583"
)

// MessageLengthReader reads message header from the provided reader interface
// and returns message length
type MessageLengthReader func(r io.Reader) (int, error)

// MessageLengthWriter writes message header with encoded length into the
// provided writer interface
type MessageLengthWriter func(w io.Writer, length int) (int, error)

type Connection struct {
	rwc io.ReadWriteCloser

	spec *iso8583.MessageSpec

	msgLenReader MessageLengthReader

	msgLenWriter MessageLengthWriter
}
