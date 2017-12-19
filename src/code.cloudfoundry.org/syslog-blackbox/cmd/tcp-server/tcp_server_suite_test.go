package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTcpServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TcpServer Suite")
}
