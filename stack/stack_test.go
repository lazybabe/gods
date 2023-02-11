package stack_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lazybabe/gods/stack"
)

func TestStack(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stack Suite")
}

var _ = Describe("Stack", func() {
	It("Push", func() {
		s := stack.NewFrom([]int{1, 2, 3})
		Expect(s.Size()).To(Equal(3))
		Expect(s.Peek()).To(Equal(3))
	})

	It("Pop", func() {
		s := stack.NewFrom([]int{1, 2, 3})
		Expect(s.Pop()).To(Equal(3))
		Expect(s.Pop()).To(Equal(2))
		Expect(s.Pop()).To(Equal(1))
		Expect(s.Pop()).To(BeZero())
	})

	It("Peek", func() {
		s := stack.NewFrom([]int{1, 2, 3})
		Expect(s.Peek()).To(Equal(3))
		Expect(s.Size()).To(Equal(3))
	})

	It("Size", func() {
		s := stack.NewFrom([]int{1})
		Expect(s.Size()).To(Equal(1))
		_ = s.Pop()
		Expect(s.Size()).To(BeZero())
	})

	It("Clone", func() {
		s := stack.NewFrom([]int{1, 2, 3})
		clone := s.Clone()
		Expect(clone.Size()).To(Equal(s.Size()))
		Expect(clone.Peek()).To(Equal(s.Peek()))
	})

	It("IsEmpty", func() {
		s := stack.New[int]()
		Expect(s.IsEmpty()).To(BeTrue())
		s.Push(1)
		Expect(s.IsEmpty()).To(BeFalse())
	})
})
