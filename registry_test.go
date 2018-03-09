package registry

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Registry", func() {
	Describe("Simple registry", func() {
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
		Describe("Given a Key and a []simpleRegistryMetric", func() {
			Context("When finding a Key that exist in []simpleRegistryMetric", func() {
				It("Then the simpleRegistryMetric should be returned", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}
					k3 := map[string]string{"a": "1"}
					v3 := 3

					r := []*simpleRegistryMetric{
						&simpleRegistryMetric{
							key:   k1,
							value: 1,
						},
						&simpleRegistryMetric{
							key:   k2,
							value: 2,
						},
						&simpleRegistryMetric{
							key:   k3,
							value: v3,
						},
					}

					m, i, err := getMetric(r, k3)
					Expect(err).To(BeNil())
					Expect(i).To(Equal(2))
					Expect(m.value).To(Equal(v3))
				})
			})
			Context("When finding a Key that does not exist in []simpleRegistryMetric", func() {
				It("Then an error should be thrown", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}
					k3 := map[string]string{"a": "1"}

					r := []*simpleRegistryMetric{
						&simpleRegistryMetric{
							key:   k1,
							value: 1,
						},
						&simpleRegistryMetric{
							key:   k2,
							value: 2,
						},
					}

					_, _, err := getMetric(r, k3)
					Expect(err).To(Equal(keyNotFound))
				})
			})
		})
		Describe("Given user wants to GET a value", func() {
			Context("When key exists", func() {
				It("Then the value is returned", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}
					v2 := 2

					r := NewSimpleRegistry()
					r.registry = []*simpleRegistryMetric{
						&simpleRegistryMetric{
							key:   k1,
							value: 1,
						},
						&simpleRegistryMetric{
							key:   k2,
							value: v2,
						},
					}
					Expect(r.Get(k2)).To(Equal(v2))
				})
			})
			Context("When key does not exists", func() {
				It("Then the value is returned", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}
					k3 := map[string]string{"a": "1"}

					r := NewSimpleRegistry()
					r.registry = []*simpleRegistryMetric{
						&simpleRegistryMetric{
							key:   k1,
							value: 1,
						},
						&simpleRegistryMetric{
							key:   k2,
							value: 2,
						},
					}
					Expect(r.Get(k3)).To(BeNil())
				})
			})
		})
		Describe("Given user wants to SET a value", func() {
			Context("When key exists", func() {
				It("Then the value is replaced", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}
					v2 := 4

					r := NewSimpleRegistry()
					r.registry = []*simpleRegistryMetric{
						&simpleRegistryMetric{
							key:   k1,
							value: 1,
						},
						&simpleRegistryMetric{
							key:   k2,
							value: 2,
						},
					}
					r.Set(k2, v2)

					Expect(r.Get(k2)).To(Equal(v2))
				})
			})
			Context("When key does not exists", func() {
				It("Then the key and value is inserted", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}
					v2 := 4

					r := NewSimpleRegistry()
					r.registry = []*simpleRegistryMetric{
						&simpleRegistryMetric{
							key:   k1,
							value: 1,
						},
					}
					r.Set(k2, v2)

					Expect(r.Get(k2)).To(Equal(v2))
				})
			})
		})
		Describe("Given user wants to DELETE a value", func() {
			Context("When key exists", func() {
				It("Then the value is deleted", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}

					r := NewSimpleRegistry()
					r.registry = []*simpleRegistryMetric{
						&simpleRegistryMetric{
							key:   k1,
							value: 1,
						},
						&simpleRegistryMetric{
							key:   k2,
							value: 2,
						},
					}
					r.Delete(k2)

					Expect(r.registry).To(HaveLen(1))
					Expect(r.Get(k1)).ToNot(BeNil())
					Expect(r.Get(k2)).To(BeNil())
				})
			})
			Context("When key does not exists", func() {
				It("Then the key and value is inserted", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}
					k3 := map[string]string{"a": "1"}

					r := NewSimpleRegistry()
					r.registry = []*simpleRegistryMetric{
						&simpleRegistryMetric{
							key:   k1,
							value: 1,
						},
						&simpleRegistryMetric{
							key:   k2,
							value: 2,
						},
					}
					r.Delete(k3)

					Expect(r.registry).To(HaveLen(2))
					Expect(r.Get(k1)).ToNot(BeNil())
					Expect(r.Get(k2)).ToNot(BeNil())
				})
			})
		})
	})
})
