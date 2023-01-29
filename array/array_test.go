package array_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/qinyuguang/gods/array"
)

func TestArray(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Array Suite")
}

var _ = Describe("Array", func() {
	It("Int slice", func() {
		slice := []int{1, 2, 3}
		a1 := array.NewFrom(slice)
		Expect(a1.Size()).To(Equal(3))
		Expect(a1.Slice()).To(ConsistOf(slice))
		a2 := array.New[int]()
		Expect(a2.Size()).To(BeZero())
		Expect(a2.Slice()).To(BeEmpty())
		Expect(a1.Slice()).To(ConsistOf(slice))
		a3 := array.NewSize[int](1, 3)
		Expect(a3.Size()).To(Equal(1))
		Expect(a3.Slice()).To(ConsistOf([]int{0}))
	})

	It("String slice", func() {
		slice := []string{"a", "b", "c"}
		a1 := array.NewFrom(slice)
		Expect(a1.Size()).To(Equal(3))
		Expect(a1.Slice()).To(ConsistOf(slice))
		a2 := array.New[string]()
		Expect(a2.Size()).To(BeZero())
		Expect(a2.Slice()).To(BeEmpty())
		Expect(a1.Slice()).To(ConsistOf(slice))
		a3 := array.NewSize[string](1, 3)
		Expect(a3.Size()).To(Equal(1))
		Expect(a3.Slice()).To(ConsistOf([]string{""}))
		e1 := a3.Set(0, "a")
		Expect(e1).NotTo(HaveOccurred())
		Expect(a3.Size()).To(Equal(1))
		Expect(a3.Slice()).To(ConsistOf([]string{"a"}))
		e2 := a3.Set(5, "b")
		Expect(e2).To(HaveOccurred())
	})

	DescribeTable("Index",
		func(output int, slice []int, index int) {
			Expect(array.NewFrom(slice).Index(index)).To(Equal(output))
		},
		Entry("out of range left", 0, []int{1, 2, 3}, -1),
		Entry("out of range right", 0, []int{1, 2, 3}, 3),
		Entry("index0", 1, []int{1, 2, 3}, 0),
		Entry("index1", 2, []int{1, 2, 3}, 1),
	)

	It("Set", func() {
		var err error
		a := array.NewFrom([]int{1, 2, 3})
		err = a.Set(0, 3)
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Slice()).To(Equal([]int{3, 2, 3}))
		err = a.Set(-1, 1)
		Expect(err).To(HaveOccurred())
		err = a.Set(3, 1)
		Expect(err).To(HaveOccurred())
	})

	It("Sort", func() {
		a1 := array.NewFrom([]int{1, 3, 2, 1})
		a1.Sort(func(v1, v2 int) bool { return v1 < v2 })
		Expect(a1.Slice()).To(Equal([]int{1, 1, 2, 3}))
		a2 := array.NewFrom([]string{"c", "a", "b", "a"})
		a2.Sort(func(v1, v2 string) bool { return v1 < v2 })
		Expect(a2.Slice()).To(Equal([]string{"a", "a", "b", "c"}))
	})

	It("InsertBefore", func() {
		var err error
		a := array.NewFrom([]int{1, 2, 3})
		err = a.InsertBefore(0, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Slice()).To(Equal([]int{0, 1, 2, 3}))
		err = a.InsertBefore(2, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Slice()).To(Equal([]int{0, 1, 0, 2, 3}))
		err = a.InsertBefore(-1, 1)
		Expect(err).To(HaveOccurred())
		err = a.InsertBefore(99, 1)
		Expect(err).To(HaveOccurred())
	})

	It("InsertAfter", func() {
		var err error
		a := array.NewFrom([]int{1, 2, 3})
		err = a.InsertAfter(0, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Slice()).To(Equal([]int{1, 0, 2, 3}))
		err = a.InsertAfter(2, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Slice()).To(Equal([]int{1, 0, 2, 0, 3}))
		err = a.InsertAfter(-1, 1)
		Expect(err).To(HaveOccurred())
		err = a.InsertAfter(99, 1)
		Expect(err).To(HaveOccurred())
	})

	It("Remove", func() {
		var (
			value int
			found bool
		)
		a := array.NewFrom([]int{1, 2, 3, 4})
		value, found = a.Remove(-1)
		Expect(value).To(BeZero())
		Expect(found).To(BeFalse())
		value, found = a.Remove(9)
		Expect(value).To(BeZero())
		Expect(found).To(BeFalse())
		value, found = a.Remove(0)
		Expect(value).To(Equal(1))
		Expect(found).To(BeTrue())
		value, found = a.Remove(1)
		Expect(value).To(Equal(3))
		Expect(found).To(BeTrue())
		value, found = a.Remove(1)
		Expect(value).To(Equal(4))
		Expect(found).To(BeTrue())
	})

	It("Remove", func() {
		a := array.NewFrom([]int{1, 2, 3})
		Expect(a.RemoveValue(2)).To(BeTrue())
		Expect(a.RemoveValue(0)).To(BeFalse())
	})

	It("PushLeft", func() {
		a := array.NewFrom([]int{1, 2, 3})
		Expect(a.PushLeft(0).Slice()).To(Equal([]int{0, 1, 2, 3}))
		Expect(a.PushLeft(-1).Slice()).To(Equal([]int{-1, 0, 1, 2, 3}))
	})

	It("PushRight|Append", func() {
		a := array.NewFrom([]int{1, 2, 3})
		Expect(a.PushRight(0).Slice()).To(Equal([]int{1, 2, 3, 0}))
		Expect(a.Append(-1).Slice()).To(Equal([]int{1, 2, 3, 0, -1}))
	})

	It("PopLeft", func() {
		var (
			value int
			found bool
		)
		a := array.NewFrom([]int{1, 2})
		value, found = a.PopLeft()
		Expect(value).To(Equal(1))
		Expect(found).To(BeTrue())
		value, found = a.PopLeft()
		Expect(value).To(Equal(2))
		Expect(found).To(BeTrue())
		value, found = a.PopLeft()
		Expect(value).To(BeZero())
		Expect(found).To(BeFalse())
	})

	It("PopRight", func() {
		var (
			value int
			found bool
		)
		a := array.NewFrom([]int{1, 2})
		value, found = a.PopRight()
		Expect(value).To(Equal(2))
		Expect(found).To(BeTrue())
		value, found = a.PopRight()
		Expect(value).To(Equal(1))
		Expect(found).To(BeTrue())
		value, found = a.PopRight()
		Expect(value).To(BeZero())
		Expect(found).To(BeFalse())
	})

	It("SubSlice", func() {
		a := array.NewFrom([]int{1, 2, 3})
		Expect(a.SubSlice(0, 1)).To(Equal([]int{1}))
		Expect(a.SubSlice(4, 1)).To(BeNil())
		Expect(a.SubSlice(3, 1)).To(BeEmpty())
		Expect(a.SubSlice(-1, 1)).To(Equal([]int{3}))
		Expect(a.SubSlice(-4, 1)).To(BeNil())
		Expect(a.SubSlice(1, -1)).To(Equal([]int{1}))
		Expect(a.SubSlice(1, -2)).To(BeNil())
		Expect(a.SubSlice(1, 3)).To(Equal([]int{2, 3}))
	})

	It("Clone", func() {
		a1 := array.NewFrom([]int{1, 2, 3}, true)
		a2 := a1.Clone()
		a3 := array.NewFrom([]int{1, 2, 3}, false)
		Expect(a2).To(Equal(a1))
		Expect(a2).NotTo(Equal(a3))
	})

	It("Clear", func() {
		a := array.NewFrom([]int{1, 2, 3})
		Expect(a.Size()).To(Equal(3))
		a.Clear()
		Expect(a.Size()).To(BeZero())
	})

	It("Contains", func() {
		a := array.NewFrom([]int{1, 2, 3})
		Expect(a.Contains(2)).To(BeTrue())
		Expect(a.Contains(0)).To(BeFalse())
	})

	It("Unique", func() {
		a := array.NewFrom([]int{2, 3, 1, 2, 1, 4})
		Expect(a.Unique().Slice()).To(Equal([]int{2, 3, 1, 4}))
	})

	It("Fill", func() {
		var err error
		a := array.NewFrom([]int{1, 2, 3, 4})
		err = a.Fill(1, 2, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Slice()).To(Equal([]int{1, 0, 0, 4}))
		err = a.Fill(-1, 2, 0)
		Expect(err).To(HaveOccurred())
		err = a.Fill(4, 2, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Slice()).To(Equal([]int{1, 0, 0, 4, 0, 0}))
	})

	It("Chunk", func() {
		a := array.NewFrom([]int{1, 2, 3, 4, 5})
		Expect(a.Chunk(3)).To(Equal([][]int{{1, 2, 3}, {4, 5}}))
		Expect(a.Chunk(100)).To(Equal([][]int{{1, 2, 3, 4, 5}}))
		Expect(a.Chunk(0)).To(BeNil())
	})

	It("Reverse", func() {
		Expect(array.NewFrom([]int{1, 2, 3}).Reverse().Slice()).To(Equal([]int{3, 2, 1}))
		Expect(array.NewFrom([]int{1, 2, 3, 4}).Reverse().Slice()).To(Equal([]int{4, 3, 2, 1}))
	})

	It("Each", func() {
		a := array.NewFrom([]int{1, 2, 3})

		var sum int
		a.Each(func(_, v int) bool {
			sum += v
			return true
		})
		Expect(sum).To(Equal(6))

		var partialSum int
		a.Each(func(k, v int) bool {
			if k == 1 {
				return false
			}
			partialSum += v
			return true
		})
		Expect(partialSum).To(Equal(1))
	})

	It("String", func() {
		Expect(array.NewFrom([]int{1, 2, 3}).String()).To(Equal(`[1 2 3]`))
		Expect(array.NewFrom([]string{"c", "b", "a"}).String()).To(Equal(`[c b a]`))
	})
})
