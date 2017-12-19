package syslog_test

import (
	"fmt"
	"net"
	"sync"
	"time"

	"code.cloudfoundry.org/rfc5424"
	"code.cloudfoundry.org/syslog-blackbox/pkg/syslog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Listener", func() {
	var (
		listener *syslog.Listener
		counter  *spyCounter
	)

	BeforeEach(func() {
		counter = newSpyCounter()
		listener = syslog.NewListener(":0", counter.count)
		listener.Run(false)
	})

	AfterEach(func() {
		_ = listener.Stop()
	})

	It("maintains counts of logs received for a given ID", func() {
		time.Sleep(time.Second)
		conn, err := net.Dial("tcp", listener.Addr())
		Expect(err).ToNot(HaveOccurred())

		msg := message("test-1")
		msg.WriteTo(conn)

		Eventually(counter.ids).Should(ConsistOf("test-1"))
		Eventually(counter.primers).Should(ConsistOf(2))
		Eventually(counter.messages).Should(ConsistOf(1))
	})
})

func message(id string) rfc5424.Message {
	return rfc5424.Message{
		Message: []byte(fmt.Sprintf(messageTemplate, id)),
	}
}

type spyCounter struct {
	mu        sync.Mutex
	_ids      []string
	_primers  []int
	_messages []int
}

func newSpyCounter() *spyCounter {
	return &spyCounter{}
}

func (s *spyCounter) count(id string, primers, msgs int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s._ids = append(s._ids, id)
	s._primers = append(s._primers, primers)
	s._messages = append(s._messages, msgs)
}

func (s *spyCounter) ids() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s._ids
}

func (s *spyCounter) primers() []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s._primers
}

func (s *spyCounter) messages() []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s._messages
}

var messageTemplate = `{
	"id": "%s",
	"cycles": "20",
	"delay": "20ns",
	"msgCount": 1,
	"primeCount": 2,
	"iteration": "20"
}`
