/**
 * @author Jose Nidhin
 */
package main

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgLenReader(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		Hex    string
		Length int
	}{
		{
			Hex:    "0113",
			Length: 275,
		},
		{
			Hex:    "009d",
			Length: 157,
		},
	}

	for i, c := range cases {
		caseNo := i + 1

		data, err := hex.DecodeString(c.Hex)
		assert.NoError(err, "Case %d - Invalid Hex provided", caseNo)

		length, err := MsgLenReader(bytes.NewReader(data))
		assert.NoError(err, "Case %d - Expected MsgLenReader to succeed without error", caseNo)

		assert.Equal(c.Length, length, "Case %d - Expected length to be equal", caseNo)
	}
}

func TestMsgLenWriter(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		Hex    string
		Length int
	}{
		{
			Hex:    "0113",
			Length: 275,
		},
		{
			Hex:    "009d",
			Length: 157,
		},
	}

	for i, c := range cases {
		caseNo := i + 1
		var buf bytes.Buffer

		wrote, err := MsgLenWriter(&buf, c.Length)
		assert.NoError(err, "Case %d - Expected MsgLenWriter to succeed without error", caseNo)

		assert.Equal(c.Hex, hex.EncodeToString(buf.Bytes()), "Case %d - Expected header to be equal", caseNo)
		assert.Equal(2, wrote, "Expected wrote to be 2", caseNo)
	}
}
