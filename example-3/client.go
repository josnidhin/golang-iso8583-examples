/**
 * @author Jose Nidhin
 */
package main

import (
	"encoding/hex"
	"net"
	"time"

	"github.com/pkg/errors"
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

	fnName := "client.Start"
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
	fnName := "client.Shutdown"
	logger.Printf("%s: graceful shutdown initialised", fnName)

	close(c.shutdownNotifier)
}
