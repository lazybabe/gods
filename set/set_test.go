package set_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/qinyuguang/gods/set"
)

func TestSet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Set Suite")
}

var _ = Describe("Set", func() {
	It("Each", func() {
		s := set.New[int]()
		s.Add(1, 2, 3)

		var sum int
		s.Each(func(v int) bool {
			sum += v
			return true
		})
		Expect(sum).To(Equal(6))

		var partialSum int
		s.Each(func(v int) bool {
			if v == 1 {
				return false
			}
			partialSum += v
			return true
		})
		Expect(partialSum).To(BeNumerically("<", 6))
	})

	It("New from int slice", func() {
		slice := []int{1, 2, 3}
		s := set.NewFrom(slice)
		Expect(s.Size()).To(Equal(3))
		Expect(s.Slice()).To(ConsistOf(slice))
		s.Add(4)
		Expect(s.Size()).To(Equal(4))
		Expect(s.Contains(4)).To(BeTrue())
		s.Remove(3)
		Expect(s.Size()).To(Equal(3))
		Expect(s.Contains(3)).To(BeFalse())
	})

	It("New from string slice", func() {
		slice := []string{"abc", "def", "ghi"}
		s := set.NewFrom(slice)
		Expect(s.Size()).To(Equal(3))
		Expect(s.Slice()).To(ConsistOf(slice))
		s.Add("jkl")
		Expect(s.Size()).To(Equal(4))
		Expect(s.Contains("jkl")).To(BeTrue())
		s.Remove("ghi")
		Expect(s.Size()).To(Equal(3))
		Expect(s.Contains("ghi")).To(BeFalse())
	})

	It("New from struct slice", func() {
		type S struct {
			ID int
		}
		slice := []S{{ID: 1}, {ID: 2}, {ID: 3}}
		s := set.NewFrom(slice)
		Expect(s.Size()).To(Equal(3))
		Expect(s.Slice()).To(ConsistOf(slice))
		s.Add(S{ID: 4})
		Expect(s.Size()).To(Equal(4))
		Expect(s.Contains(S{ID: 4})).To(BeTrue())
		s.Remove(S{ID: 3})
		Expect(s.Size()).To(Equal(3))
		Expect(s.Contains(S{ID: 3})).To(BeFalse())
	})

	It("Add with manually instance", func() {
		s := &set.Set[int]{}
		slice := []int{1, 2, 3}
		s.Add(slice...)
		Expect(s.Size()).To(Equal(3))
		Expect(s.Slice()).To(ConsistOf(slice))
	})

	It("Clear", func() {
		s := set.NewFrom([]int{1, 2, 3})
		Expect(s.Size()).To(Equal(3))
		s.Clear()
		Expect(s.Size()).To(BeZero())
	})

	It("String", func() {
		s1 := set.NewFrom([]int{3, 2, 1})
		Expect(s1.String()).To(Equal(`[1 2 3]`))
		s2 := set.NewFrom([]string{"c", "b", "a"})
		Expect(s2.String()).To(Equal(`[a b c]`))
	})

	It("Clone", func() {
		s1 := set.NewFrom([]int{1, 2, 3}, true)
		s2 := s1.Clone()
		s3 := set.NewFrom([]int{1, 2, 3}, false)
		Expect(s1 == s2).To(BeFalse())
		Expect(s1).To(Equal(s2))
		Expect(s2).NotTo(Equal(s3))
	})

	It("Equal", func() {
		s1 := set.NewFrom([]int{1, 2, 3})
		s2 := set.NewFrom([]int{3, 2, 1})
		s3 := s1
		s4 := set.NewFrom([]int{1, 2})
		s5 := set.NewFrom([]int{1, 2, 4})
		Expect(s1.Equal(s2)).To(BeTrue())
		Expect(s1.Equal(s3)).To(BeTrue())
		Expect(s1.Equal(s4)).To(BeFalse())
		Expect(s1.Equal(s5)).To(BeFalse())
		Expect(s1.Equal(nil)).To(BeFalse())
	})

	It("IsSubsetOf", func() {
		s1 := set.NewFrom([]int{1, 2, 3})
		s2 := set.NewFrom([]int{1, 2, 3, 4})
		s3 := s1
		s4 := set.NewFrom([]int{1, 2})
		s5 := set.NewFrom([]int{1, 2, 4})
		Expect(s1.IsSubsetOf(s2)).To(BeTrue())
		Expect(s1.IsSubsetOf(s3)).To(BeTrue())
		Expect(s1.IsSubsetOf(s4)).To(BeFalse())
		Expect(s1.IsSubsetOf(s5)).To(BeFalse())
		Expect(s1.IsSubsetOf(nil)).To(BeFalse())
	})

	It("Union", func() {
		s1 := set.NewFrom([]int{1, 2, 3})
		s2 := set.NewFrom([]int{1, 2, 3, 4})
		s3 := s1
		Expect(s1.Union(s2, s3, nil)).To(Equal(set.NewFrom([]int{1, 2, 3, 4})))
	})

	It("Diff", func() {
		s1 := set.NewFrom([]int{1, 2, 4})
		s2 := set.NewFrom([]int{1, 3})
		s3 := set.NewFrom([]int{4, 6})
		Expect(s1.Diff(s2)).To(Equal(set.NewFrom([]int{2, 4})))
		Expect(s1.Diff(s2, s3)).To(Equal(set.NewFrom([]int{2})))
		Expect(s1.Diff(s2, s3, nil)).To(Equal(set.NewFrom([]int{2})))
	})

	It("Intersect", func() {
		s1 := set.NewFrom([]int{1, 2, 4})
		s2 := set.NewFrom([]int{1, 2})
		s3 := set.NewFrom([]int{1, 4})
		Expect(s1.Intersect(s2)).To(Equal(set.NewFrom([]int{1, 2})))
		Expect(s1.Intersect(s2, s3)).To(Equal(set.NewFrom([]int{1})))
		Expect(s1.Intersect(s2, nil)).To(Equal(set.New[int]()))
	})
})
