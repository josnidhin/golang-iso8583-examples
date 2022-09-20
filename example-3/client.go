/**
 * @author Jose Nidhin
 */
package main

import (
	"encoding/hex"
	"net"
	"time"

	"github.com/moov-io/iso8583"
	"github.com/pkg/errors"
)

const (
	emHexStr = "004349534F30323131303030353530383030383232303030303030303030303030303034303030303030303030303030303030383231303833323136303135373935333031"
	fmHexStr = "009d49534f303234303030303535303230303532333838303030303841303830303031303831313030393934313830303030303030313030303030333133313032383432343838373539313032363431303331333033313330303030303034303139393131304d4f4e353047415a4f582020204e456469736f6e203132333520202020202020202020204d6f6e746572726579202020204e4c204d58343834"
)

var hexStr = emHexStr

var sampleRawInput []byte

func init() {
	var err error
	fnName := "client.init"

	sampleRawInput, err = hex.DecodeString(hexStr)
	if err != nil {
		logger.Panicf("%s: raw input creation failed - %v", fnName, err)
	}
}

type Client struct {
	network          string
	tcpAddr          *net.TCPAddr
	shutdownNotifier chan struct{}
}

func NewClient(address string) (*Client, error) {
	network := "tcp"

	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, errors.Wrapf(err, "address resolve failed")
	}

	client := &Client{
		network:          network,
		tcpAddr:          tcpAddr,
		shutdownNotifier: make(chan struct{}),
	}

	return client, nil
}

func (c *Client) Start() {
	var tcpConn *net.TCPConn
	var err error

	fnName := "Client.Start"
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-c.shutdownNotifier:
			break loop
		case <-ticker.C:
			tcpConn, err = net.DialTCP(c.network, nil, c.tcpAddr)
			if err != nil {
				logger.Printf("%s: dial failed - %v", fnName, err)
				break
			}

			break loop
		}

		err = nil
	}

	if tcpConn == nil {
		return
	}

	defer tcpConn.Close()

	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(60 * time.Second)

	go c.readResp(tcpConn)

	err = nil
	for {
		select {
		case <-c.shutdownNotifier:
			return
		case <-ticker.C:
			_, err = tcpConn.Write(sampleRawInput)
			if err != nil {
				logger.Printf("%s: error while writing to tcp connection - %v", fnName, err)
				return
			}
		}

		err = nil
	}
}

func (c *Client) Shutdown() {
	fnName := "Client.Shutdown"
	logger.Printf("%s: graceful shutdown initialised", fnName)

	close(c.shutdownNotifier)
}

func (c *Client) readResp(tcpConn *net.TCPConn) {
	fnName := "Client.readResp"

	reqCh := make(chan *iso8583.Message)
	go func() {
		for {
			// discard
			<-reqCh
		}
	}()

	resCh := make(chan *iso8583.Message)

	connHandler, err := NewConnectionHandler(tcpConn, Spec1HeaderSize, Spec1, MsgLenReader, MsgLenWriter, reqCh, resCh)
	if err != nil {
		logger.Printf("%s: error creating connection handler - %v", fnName, err)
	}

	connHandler.Start()

	<-c.shutdownNotifier

	connHandler.Close()
}
