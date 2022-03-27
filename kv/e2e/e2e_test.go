package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
)

var _ = Describe("E2e", func() {
	Expect(actual).To(gstruct.MatchAllFields(gstruct.Fields{
		"A": BeNumerically("<", 10),
		"B": BeTrue(),
		"C": Equal("foo"),
	}))
	Expect(actual).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
		"A": BeNumerically("<", 10),
		"B": BeTrue(),
		// 忽略C字段
	}))
	Expect(actual).To(gstruct.MatchFields(gstruct.IgnoreMissing, gstruct.Fields{
		"A": BeNumerically("<", 10),
		"B": BeTrue(),
		"C": Equal("foo"),
		"D": Equal("bar"), // 忽略多余字段
	}))
})

var actual = struct {
	A int
	B bool
	C string
}{5, true, "foo"}

// https://blog.gmem.cc/ginkgo-study-note
