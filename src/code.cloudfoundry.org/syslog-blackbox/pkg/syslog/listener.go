package syslog

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"time"

	"code.cloudfoundry.org/rfc5424"
)

// Counter is a function that is used by the Listener to report information
// about messages it receives.
type Counter func(id string, primers, msgs int)

// Listener will listen for syslog messages and keep counts for received
// messages. It expects message to follow the format emitted by the jsonspinner
// from loggregator tools:
//   https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/jsonspinner
type Listener struct {
	addr    string
	lis     net.Listener
	counter Counter
}

// NewListener initializes and returns a Listener.
func NewListener(addr string, c Counter) *Listener {
	return &Listener{addr: addr, counter: c}
}

// Run starts the Listener. The Listener will listen on the
// configured port and start accepting connections. This method takes blocking
// as a bool. If set to true this method will block, otherwise it will run in a
// goroutine.
func (sl *Listener) Run(blocking bool) {
	l, err := net.Listen("tcp", sl.addr)
	if err != nil {
		log.Fatal(err)
	}

	sl.lis = l
	log.Printf("Listening at %s", sl.Addr())

	if !blocking {
		go sl.acceptConns()
		return
	}

	sl.acceptConns()
}

// Addr returns the address the Listener is listening on.
func (sl *Listener) Addr() string {
	return sl.lis.Addr().String()
}

// Stop stops the Listener.
func (sl *Listener) Stop() error {
	return sl.lis.Close()
}

func (sl *Listener) acceptConns() {
	for {
		conn, err := sl.lis.Accept()
		if err != nil {
			log.Printf("Error accepting: %s", err)
			time.Sleep(20 * time.Millisecond)
			continue
		}
		log.Printf("Accepted connection")

		go sl.handle(conn)
	}
}

func (sl *Listener) handle(conn net.Conn) {
	defer conn.Close()

	var msg rfc5424.Message
	for {
		_, err := msg.ReadFrom(conn)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("ReadFrom err: %s", err)
			return
		}

		var smsg syslogMessage
		err = json.Unmarshal(msg.Message, &smsg)
		if err != nil {
			log.Printf("failed to unmarshal message payload: %s", err)
			continue
		}

		sl.counter(smsg.ID, smsg.PrimeCount, smsg.MsgCount)
	}
}

type syslogMessage struct {
	ID         string `json:"id"`
	MsgCount   int    `json:"msgCount"`
	PrimeCount int    `json:"primeCount"`
}
