package registry

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Registry", func() {
	Describe("Given two Keys", func() {
		Context("When the two Keys are the same", func() {
			It("Then isEquals should return true", func() {
				k1 := map[string]string{"a": "1", "b": "2"}
				k2 := k1
				Expect(isEquals(k1, k2)).To(BeTrue())
			})
		})
		Context("When the two Keys are different length", func() {
			It("Then isEquals should return false", func() {
				k1 := map[string]string{"a": "1", "b": "2"}
				k2 := map[string]string{"a": "1"}
				Expect(isEquals(k1, k2)).To(BeFalse())
			})
		})
		Context("When the two Keys have different keys", func() {
			It("Then isEquals should return false", func() {
				k1 := map[string]string{"a": "1", "b": "2"}
				k2 := map[string]string{"a": "1", "c": "2"}
				Expect(isEquals(k1, k2)).To(BeFalse())
			})
		})
		Context("When the two Keys have different values", func() {
			It("Then isEquals should return false", func() {
				k1 := map[string]string{"a": "1", "b": "2"}
				k2 := map[string]string{"a": "1", "b": "3"}
				Expect(isEquals(k1, k2)).To(BeFalse())
			})
		})
	})
})
