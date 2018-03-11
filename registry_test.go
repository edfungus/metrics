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
			Context("When the Key contains a subset of another Key", func() {
				It("Then isSubset should return true", func() {
					k1 := map[string]string{"a": "1", "b": "2", "c": "3"}
					k2 := map[string]string{"a": "1", "b": "2"}
					Expect(isSubset(k2, k1)).To(BeTrue())
				})
			})
			Context("When the Key contains a superset of another Key", func() {
				It("Then isSubset should return false", func() {
					k1 := map[string]string{"a": "1", "b": "2", "c": "3"}
					k2 := map[string]string{"a": "1", "b": "2"}
					Expect(isSubset(k1, k2)).To(BeFalse())
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

					r := []*Entry{
						&Entry{
							Key:   k1,
							Value: 1,
						},
						&Entry{
							Key:   k2,
							Value: 2,
						},
						&Entry{
							Key:   k3,
							Value: v3,
						},
					}

					m, i, err := getEntry(r, k3)
					Expect(err).To(BeNil())
					Expect(i).To(Equal(2))
					Expect(m.Value).To(Equal(v3))
				})
			})
			Context("When finding a Key that does not exist in []simpleRegistryMetric", func() {
				It("Then an error should be thrown", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}
					k3 := map[string]string{"a": "1"}

					r := []*Entry{
						&Entry{
							Key:   k1,
							Value: 1,
						},
						&Entry{
							Key:   k2,
							Value: 2,
						},
					}

					_, _, err := getEntry(r, k3)
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
					r.registry = []*Entry{
						&Entry{
							Key:   k1,
							Value: 1,
						},
						&Entry{
							Key:   k2,
							Value: v2,
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
					r.registry = []*Entry{
						&Entry{
							Key:   k1,
							Value: 1,
						},
						&Entry{
							Key:   k2,
							Value: 2,
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
					r.registry = []*Entry{
						&Entry{
							Key:   k1,
							Value: 1,
						},
						&Entry{
							Key:   k2,
							Value: 2,
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
					r.registry = []*Entry{
						&Entry{
							Key:   k1,
							Value: 1,
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
					r.registry = []*Entry{
						&Entry{
							Key:   k1,
							Value: 1,
						},
						&Entry{
							Key:   k2,
							Value: 2,
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
					r.registry = []*Entry{
						&Entry{
							Key:   k1,
							Value: 1,
						},
						&Entry{
							Key:   k2,
							Value: 2,
						},
					}
					r.Delete(k3)

					Expect(r.registry).To(HaveLen(2))
					Expect(r.Get(k1)).ToNot(BeNil())
					Expect(r.Get(k2)).ToNot(BeNil())
				})
			})
		})
		Describe("Given user wants to FILTER for values", func() {
			Context("When keys exist", func() {
				It("Then the value is deleted", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}
					k3 := map[string]string{"a": "1"}

					r := NewSimpleRegistry()
					r.registry = []*Entry{
						&Entry{
							Key:   k1,
							Value: 1,
						},
						&Entry{
							Key:   k2,
							Value: 2,
						},
					}
					v := r.Filter(k3)

					Expect(v).To(HaveLen(len(r.registry)))
				})
			})
			Context("When key does not exists", func() {
				It("Then the key and value is inserted", func() {
					k1 := map[string]string{"a": "1", "b": "2"}
					k2 := map[string]string{"a": "1", "c": "2"}
					k3 := map[string]string{"a": "2"}

					r := NewSimpleRegistry()
					r.registry = []*Entry{
						&Entry{
							Key:   k1,
							Value: 1,
						},
						&Entry{
							Key:   k2,
							Value: 2,
						},
					}
					v := r.Filter(k3)

					Expect(v).To(HaveLen(0))
				})
			})
		})
	})
	Describe("Better registry", func() {
		// yah.... skipping some unit tests here ... :)
		Describe("Given a user is using the better registry", func() {
			var br *BetterRegistry
			var k1, k2, k3 Key
			Context("Setup better registry", func() {
				br = NewBetterRegistry()
				k1 = map[string]string{"a": "1", "b": "2"}
				k2 = map[string]string{"a": "1", "b": "3"}
				k3 = map[string]string{"a": "1"}
				br.Set(k1, "k1")
				br.Set(k2, "k2")
				br.Set(k3, "k3")
			})
			Context("When users Sets a value", func() {
				It("Then they should be able to Get the value", func() {
					Expect(br.Get(k1)).To(Equal("k1"))
					Expect(br.Get(k2)).To(Equal("k2"))
					Expect(br.Get(k3)).To(Equal("k3"))
				})
			})
			Context("When users Filters a key", func() {
				It("Then they should be able to Get the entries with the filtered key", func() {
					entries := br.Filter(k3)
					Expect(entries).To(HaveLen(3))
				})
			})
			Context("When users Sets an existing value", func() {
				It("Then the value should be updated", func() {
					br.Set(k3, "new k3")
					Expect(br.Get(k1)).To(Equal("k1"))
					Expect(br.Get(k2)).To(Equal("k2"))
					Expect(br.Get(k3)).To(Equal("new k3"))
					br.Set(k3, "k3")
				})
			})
			Context("When users delete a key that does not exist", func() {
				It("Then nothing should have changed", func() {
					k4 := map[string]string{"a": "1", "b": "4"}
					br.Delete(k4)
					Expect(br.Get(k1)).To(Equal("k1"))
					Expect(br.Get(k2)).To(Equal("k2"))
					Expect(br.Get(k3)).To(Equal("k3"))
				})
			})
			Context("When users delete a key that does exist", func() {
				It("Then only that key is deleted", func() {
					br.Delete(k2)
					Expect(br.Get(k1)).To(Equal("k1"))
					Expect(br.Get(k2)).To(BeNil())
					Expect(br.Get(k3)).To(Equal("k3"))
					br.Set(k2, "k2")
					Expect(br.Get(k2)).To(Equal("k2"))
				})
			})
			Context("When users delete a key that has unique key field", func() {
				It("Then that field should be deleted as well", func() {
					k4 := map[string]string{"a": "1", "b": "2", "c": "4"}
					br.Set(k4, "4")
					Expect(br.registry).To(HaveLen(3))
					br.Delete(k4)
					Expect(br.registry).To(HaveLen(2))
				})
			})
		})
	})
})
