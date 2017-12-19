package stat_test

import (
	"code.cloudfoundry.org/syslog-blackbox/pkg/stat"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Counter", func() {
	It("maintains counts for multiple IDs", func() {
		c := stat.NewCounter()

		c.Add("id-a", 1, 1)
		c.Add("id-b", 1, 2)
		c.Add("id-a", 1, 1)
		c.Add("id-c", 1, 1)

		primers, msgs := c.Counts("id-a")
		Expect(primers).To(Equal(2))
		Expect(msgs).To(Equal(2))

		primers, msgs = c.Counts("id-b")
		Expect(primers).To(Equal(1))
		Expect(msgs).To(Equal(2))

		primers, msgs = c.Counts("id-c")
		Expect(primers).To(Equal(1))
		Expect(msgs).To(Equal(1))
	})
})
