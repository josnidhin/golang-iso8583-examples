/**
 * @author Jose Nidhin
 */
package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/field"
	"github.com/pkg/errors"
)

var sampleFMR = FinancialMessageResponse{
	MTI:                                 field.NewStringValue("0210"),
	PrimaryAccountNumber:                field.NewNumericValue(8110099418),
	ProcessingCode:                      field.NewStringValue("000000"),
	TransactionAmount:                   field.NewNumericValue(10000),
	TransmissionDateTime:                field.NewStringValue("0313102842"),
	STAN:                                field.NewNumericValue(488759),
	LocalTransactionTime:                field.NewStringValue("102641"),
	LocalTransactionDate:                field.NewStringValue("0313"),
	CaptureDate:                         field.NewStringValue("0313"),
	RetrievalReferenceNumber:            field.NewStringValue("000000401991"),
	AuthorizationIdentificationResponse: field.NewStringValue("123456"),
	ResponseCode:                        field.NewStringValue("00"),
	CardAcceptorTerminalIdentification:  field.NewStringValue("10MON50GAZOX   N"),
	TransactionCurrencyCode:             field.NewStringValue("484"),
}

type Server struct {
	tcpAddr          *net.TCPAddr
	tcpListener      *net.TCPListener
	wg               *sync.WaitGroup
	shutdownNotifier chan struct{}
	reqMsgCh         chan *iso8583.Message
	resMsgCh         chan *iso8583.Message
}

func NewServer(address string) (*Server, error) {
	fnName := "server.NewServer"
	network := "tcp"

	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, errors.Wrapf(err, "address resolve failed")
	}

	tcpListener, err := net.ListenTCP(network, tcpAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "bind listener failed")
	}

	server := &Server{
		tcpAddr:          tcpAddr,
		tcpListener:      tcpListener,
		wg:               &sync.WaitGroup{},
		shutdownNotifier: make(chan struct{}),
		reqMsgCh:         make(chan *iso8583.Message),
		resMsgCh:         make(chan *iso8583.Message),
	}

	logger.Printf("%s: server listening on address - %s", fnName, tcpAddr)

	return server, nil
}

func (s *Server) Start() {
	go s.connListenLoop()
	go s.reqMsgReadLoop()
}

func (s *Server) Shutdown() {
	fnName := "server.Shutdown"
	logger.Printf("%s: graceful shutdown initialised", fnName)

	close(s.shutdownNotifier)
	s.tcpListener.Close()
	s.wg.Wait()

	close(s.reqMsgCh)
	close(s.resMsgCh)
}

func (s *Server) connListenLoop() {
	fnName := "server.connListenLoop"

	for {
		conn, err := s.tcpListener.AcceptTCP()
		if err != nil {
			select {
			case <-s.shutdownNotifier:
				return
			default:
				logger.Printf("%s: accept connection failed - %v", fnName, err)
				break
			}
		}

		logger.Printf("%s: new connection", fnName)

		// aggresive keepalive on server to detect connection loss
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(10 * time.Second)

		connHandler, err := NewConnectionHandler(conn, Spec1HeaderSize, Spec1, MsgLenReader, MsgLenWriter, s.reqMsgCh, s.resMsgCh)
		if err != nil {
			logger.Fatalf("%s: error creating connection handler - %v", fnName, err)
		}

		connHandler.Start()

		s.wg.Add(1)
		go func(connHandler *ConnectionHandler) {
			<-s.shutdownNotifier
			defer s.wg.Done()

			err := connHandler.Close()
			if err != nil {
				logger.Printf("%s: error closing connection handler - %v", fnName, err)
			}

			connHandler.Done()
		}(connHandler)
	}
}

func (s *Server) reqMsgReadLoop() {
	for {
		select {
		case <-s.shutdownNotifier:
			return
		case msg := <-s.reqMsgCh:
			s.reqMsgHandler(msg)
		}
	}
}

func (s *Server) reqMsgHandler(msg *iso8583.Message) {
	fnName := "Server.reqMsgHandler"

	s.wg.Add(1)
	defer s.wg.Done()

	s.printISOMsg(msg)

	resMsg := iso8583.NewMessage(Spec1)
	err := resMsg.Marshal(&sampleFMR)

	if err != nil {
		logger.Printf("%s: sample response creation failed - %v", fnName, err)
		return
	}

	s.resMsgCh <- resMsg
}

func (s *Server) printISOMsg(msg *iso8583.Message) {
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
