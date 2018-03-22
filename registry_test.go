// +build all unit
package registry

import (
	"encoding/base64"
	"strconv"

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
	Describe("Cached registry", func() {
		Describe("Given a user wants to get a value", func() {
			Context("When the key exist", func() {
				var r *CachedRegistry
				var c *SimpleCache
				It("Then the value should be returned", func() {
					r = NewCacheRegistry(1)
					c = NewSimpleCache(1)
					r.getCache = c
					k := map[string]string{"k1": "v1", "k2": "v2"}
					he := &hashEntry{
						keys:  k,
						value: 1,
					}
					r.registry.registry["k1"] = map[string]hashEntries{}
					r.registry.registry["k1"]["v1"] = hashEntries{he}
					r.registry.registry["k2"] = map[string]hashEntries{}
					r.registry.registry["k2"]["v2"] = hashEntries{he}

					Expect(r.Get(k)).To(Equal(1))
				})
				It("Then it should be placed on the cache", func() {
					k := map[string]string{"k1": "v1", "k2": "v3"}
					he := &hashEntry{
						keys:  k,
						value: 2,
					}
					r.registry.registry["k1"]["v1"] = append(r.registry.registry["k1"]["v1"], he)
					r.registry.registry["k2"]["v3"] = append(r.registry.registry["k2"]["v3"], he)
					Expect(r.Get(k)).To(Equal(2))
					Expect(c.recentKeys).To(HaveLen(1))
					Expect(c.cache).To(HaveLen(1))
					Expect(c.cache[c.recentKeys[0]].(hashEntries)[0].value).To(Equal(2))
				})
			})
			Context("When the key does not exist", func() {
				It("Then nil should be returned", func() {
					r := NewCacheRegistry(1)
					k := map[string]string{"k1": "v1", "k2": "v2"}
					he := &hashEntry{
						keys:  k,
						value: 1,
					}
					r.registry.registry["k1"] = map[string]hashEntries{}
					r.registry.registry["k1"]["v1"] = hashEntries{he}
					r.registry.registry["k2"] = map[string]hashEntries{}
					r.registry.registry["k2"]["v2"] = hashEntries{he}

					newK := map[string]string{"k1": "v1", "k2": "v3"}
					Expect(r.Get(newK)).To(BeNil())
				})
			})
			Context("When the key has already been retrieved", func() {
				var r *CachedRegistry
				var c *SimpleCache
				It("Then the value should be returned from cache", func() {
					r = NewCacheRegistry(1)
					c = NewSimpleCache(1)
					r.getCache = c
					k := map[string]string{"k1": "v1", "k2": "v2"}
					he1 := &hashEntry{
						keys:  k,
						value: 1,
					}
					he2 := &hashEntry{
						keys:  k,
						value: 3,
					}
					r.registry.registry["k1"] = map[string]hashEntries{}
					r.registry.registry["k1"]["v1"] = hashEntries{he1}
					r.registry.registry["k2"] = map[string]hashEntries{}
					r.registry.registry["k2"]["v2"] = hashEntries{he1}

					Expect(r.Get(k)).To(Equal(1))

					he1.value = 2 // set new value
					// hack registry so that, if this value is returned, we know the cache was not used
					r.registry.registry["k1"]["v1"] = hashEntries{he2}
					r.registry.registry["k2"]["v2"] = hashEntries{he2}
					Expect(r.Get(k)).To(Equal(2))

					// force clear cache so we get value from registry
					c.RemoveWithHash(toHashString(k))
					Expect(r.Get(k)).To(Equal(3))
				})
				It("Then cache should not have grown", func() {
					Expect(c.recentKeys).To(HaveLen(1))
					Expect(c.cache).To(HaveLen(1))
				})
			})
		})
		Describe("Given a user wants to filter a value", func() {
			Context("When the key exist", func() {
				var r *CachedRegistry
				var c *SimpleCache
				It("Then the value should be returned", func() {
					r = NewCacheRegistry(1)
					c = NewSimpleCache(1)
					r.filterCache = c
					k1 := map[string]string{"k1": "v1", "k2": "v2"}
					he1 := &hashEntry{
						keys:  k1,
						value: 1,
					}
					k2 := map[string]string{"k1": "v1", "k2": "v3"}
					he2 := &hashEntry{
						keys:  k2,
						value: 1,
					}
					k3 := map[string]string{"k1": "v0", "k2": "v2"}
					he3 := &hashEntry{
						keys:  k3,
						value: 3,
					}

					r.registry.registry["k1"] = map[string]hashEntries{}
					r.registry.registry["k1"]["v1"] = hashEntries{he1, he2}
					r.registry.registry["k1"]["v0"] = hashEntries{he3}
					r.registry.registry["k2"] = map[string]hashEntries{}
					r.registry.registry["k2"]["v2"] = hashEntries{he1, he3}
					r.registry.registry["k2"]["v3"] = hashEntries{he2}

					k := map[string]string{"k1": "v1"}
					entries := r.Filter(k)
					Expect(entries).To(HaveLen(2))
					Expect(entries).To(ContainElement(Entry{
						Key:   he1.keys,
						Value: he1.value,
					}))
					Expect(entries).To(ContainElement(Entry{
						Key:   he2.keys,
						Value: he2.value,
					}))
					Expect(he1.filterCacheKeys).To(HaveLen(1))
					Expect(he2.filterCacheKeys).To(HaveLen(1))
				})
				It("Then it should be placed on the cache", func() {
					Expect(c.recentKeys).To(HaveLen(1))
					Expect(c.cache).To(HaveLen(1))
					Expect(c.cache[c.recentKeys[0]].(hashEntries)[0].value).To(Equal(1))
				})
			})
			Context("When the key does not exist", func() {
				It("Then nil should be returned", func() {
					r := NewCacheRegistry(1)
					k := map[string]string{"k1": "v1", "k2": "v2"}
					he := &hashEntry{
						keys:  k,
						value: 1,
					}
					r.registry.registry["k1"] = map[string]hashEntries{}
					r.registry.registry["k1"]["v1"] = hashEntries{he}
					r.registry.registry["k2"] = map[string]hashEntries{}
					r.registry.registry["k2"]["v2"] = hashEntries{he}

					newK := map[string]string{"k1": "v0"}
					Expect(r.Filter(newK)).To(HaveLen(0))
				})
			})
			Context("When the key has already been retrieved", func() {
				var r *CachedRegistry
				var c *SimpleCache
				It("Then the value should be returned from cache", func() {
					r = NewCacheRegistry(1)
					c = NewSimpleCache(1)
					r.filterCache = c
					k := map[string]string{"k1": "v1", "k2": "v2"}
					he1 := &hashEntry{
						keys:  k,
						value: 1,
					}
					he2 := &hashEntry{
						keys:  k,
						value: 3,
					}
					r.registry.registry["k1"] = map[string]hashEntries{}
					r.registry.registry["k1"]["v1"] = hashEntries{he1}
					r.registry.registry["k2"] = map[string]hashEntries{}
					r.registry.registry["k2"]["v2"] = hashEntries{he1}

					Expect(r.Filter(k)).To(HaveLen(1))
					Expect(r.Filter(k)[0].Value).To(Equal(1))

					he1.value = 2 // set new value
					// hack registry so that, if this value is returned, we know the cache was not used
					r.registry.registry["k1"]["v1"] = hashEntries{he2}
					r.registry.registry["k2"]["v2"] = hashEntries{he2}
					Expect(r.Filter(k)[0].Value).To(Equal(2))

					// force clear cache so we get value from registry
					c.RemoveWithHash(toHashString(k))
					Expect(r.Filter(k)[0].Value).To(Equal(3))
				})
				It("Then cache should not have grown", func() {
					Expect(c.recentKeys).To(HaveLen(1))
					Expect(c.cache).To(HaveLen(1))
				})
			})
		})
		Describe("Given a user wants to set a value", func() {
			Context("When the key exist", func() {
				It("Then the value should be replaced", func() {
					r := NewCacheRegistry(1)
					k1 := map[string]string{"k1": "v1", "k2": "v2"}
					he1 := &hashEntry{
						keys:  k1,
						value: 1,
					}
					r.registry.registry["k1"] = map[string]hashEntries{}
					r.registry.registry["k1"]["v1"] = hashEntries{he1}
					r.registry.registry["k2"] = map[string]hashEntries{}
					r.registry.registry["k2"]["v2"] = hashEntries{he1}

					r.Set(k1, 2)
					Expect(r.Get(k1)).To(Equal(2))
				})
			})
			Context("When the key does not exist", func() {
				It("Then a new entry should be made", func() {
					r := NewCacheRegistry(1)
					k := map[string]string{"k1": "v1"}

					Expect(r.registry.registry).To(HaveLen(0))
					r.Set(k, 3)
					Expect(r.Get(k)).To(Equal(3))
				})
			})
			Context("When the key is already in the cache", func() {
				It("Then the new value should be reflected in the cache", func() {
					r := NewCacheRegistry(1)
					k := map[string]string{"k1": "v1", "k2": "v2"}
					he1 := &hashEntry{
						keys:  k,
						value: 1,
					}
					r.registry.registry["k1"] = map[string]hashEntries{}
					r.registry.registry["k1"]["v1"] = hashEntries{he1}
					r.registry.registry["k2"] = map[string]hashEntries{}
					r.registry.registry["k2"]["v2"] = hashEntries{he1}

					Expect(r.Get(k)).To(Equal(1))
					r.Set(k, 2) // this sets the real entry value

					// Replace entry in resgistry so that if the cache is not used, it will get the wrong value
					he2 := &hashEntry{
						keys:  k,
						value: 3,
					}
					r.registry.registry["k1"]["v1"] = hashEntries{he2}
					r.registry.registry["k2"]["v2"] = hashEntries{he2}

					Expect(r.Get(k)).To(Equal(2))
				})
			})
		})
		Describe("Given a user wants to delete a key", func() {
			Context("When the key exist", func() {
				It("Then the value should be deleted from registry", func() {
					r := NewCacheRegistry(1)
					k1 := map[string]string{"k1": "v1", "k2": "v2"}
					r.Set(k1, 1)

					r.Delete(k1)
					Expect(r.registry.registry).To(HaveLen(0))
				})
			})
			Context("When the key does not exist", func() {
				It("Then nothing will happen", func() {
					r := NewCacheRegistry(1)
					k1 := map[string]string{"k1": "v1", "k2": "v2"}
					r.Set(k1, 1)

					k2 := map[string]string{"k1": "v1"}
					r.Delete(k2)
					Expect(r.registry.registry).To(HaveLen(2))
				})
			})
			Context("When the key exist also in both caches", func() {
				It("Then the caches should be updated too", func() {
					r := NewCacheRegistry(5)
					cf := NewSimpleCache(5)
					cg := NewSimpleCache(5)
					r.filterCache = cf
					r.getCache = cg

					k1 := map[string]string{"k1": "v1", "k2": "v2"}
					r.Set(k1, 1)

					k2 := map[string]string{"k1": "v1", "k3": "v3"}
					r.Set(k2, 1)

					r.Get(k1)
					r.Get(k2)
					e := r.Filter(k1)
					Expect(e).To(HaveLen(1))
					Expect(cf.cache).To(HaveLen(1))
					Expect(cg.cache).To(HaveLen(2))

					r.Delete(k1)
					Expect(r.registry.registry).To(HaveLen(2))
					Expect(cf.cache).To(HaveLen(0))
					Expect(cf.recentKeys).To(HaveLen(0))
					Expect(cg.cache).To(HaveLen(1))
					Expect(cg.recentKeys).To(HaveLen(1))
				})
			})
		})
	})
	Describe("Even Better registry", func() {
		Describe("Given wanting to remove hashEntry from a hashEntries", func() {
			Context("When the key exist", func() {
				It("Then the value should be returned and deleted", func() {
					e1 := &hashEntry{
						keys: map[string]string{"a": "b", "c": "d"},
					}
					e2 := &hashEntry{
						keys: map[string]string{"a": "b", "c": "e"},
					}
					e3 := &hashEntry{
						keys: map[string]string{"a": "b"},
					}
					entries := append(hashEntries{}, e1, e2, e3)
					entries, entry := removeFromHashEntries(entries, e3.keys)
					Expect(entry).ToNot(BeNil())
					Expect(entry).To(Equal(e3))
					Expect(entries).To(HaveLen((2)))
					Expect(entries).To(ContainElement(e1))
					Expect(entries).To(ContainElement(e2))
				})
			})
			Context("When the key does not exist", func() {
				It("Then the value should be returned and deleted", func() {
					e1 := &hashEntry{
						keys: map[string]string{"a": "b", "c": "d"},
					}
					e2 := &hashEntry{
						keys: map[string]string{"a": "b", "c": "e"},
					}
					k3 := map[string]string{"a": "b"}
					entries := append(hashEntries{}, e1, e2)
					entries, entry := removeFromHashEntries(entries, k3)
					Expect(entry).To(BeNil())
					Expect(entries).To(HaveLen((2)))
					Expect(entries).To(ContainElement(e1))
					Expect(entries).To(ContainElement(e2))
				})
			})
		})
	})
	Describe("Cache", func() {
		Describe("Given a Key", func() {
			Describe("func sortedKeys", func() {
				Context("When the Key needs to be hashed", func() {
					It("Then the Key's keys needed to be sorted", func() {
						m := map[string]string{"b": "1", "c": "2", "a": "0"}
						sortedKeys := sortedKeys(m)
						for i, k := range sortedKeys {
							Expect(m[k]).To(Equal(strconv.Itoa(i)))
						}
					})
				})
			})
			Describe("func toHashString", func() {
				Context("When the Key needs to be hashed", func() {
					It("Then Keys with the same value should get the same hash", func() {
						k1 := map[string]string{"a": "1", "b": "2"}
						k2 := map[string]string{"b": "2", "a": "1"}
						hk1 := toHashString(k1)
						hk2 := toHashString(k2)
						expectedHash := base64.StdEncoding.EncodeToString([]byte("a")) + ":" + base64.StdEncoding.EncodeToString([]byte("1")) + "," + base64.StdEncoding.EncodeToString([]byte("b")) + ":" + base64.StdEncoding.EncodeToString([]byte("2")) + ","
						Expect(hk1).To(Equal(expectedHash))
						Expect(hk1).To(Equal(hk2))
					})
				})
			})
		})
		Describe("Given a hash", func() {
			Describe("func UpdateWithHash", func() {
				Context("When inserting a new value", func() {
					It("Then both the cache and recentKeys is updated", func() {
						c := NewSimpleCache(3)
						c.cache["k1"] = 1
						c.cache["k2"] = 2
						c.recentKeys = []string{"k1", "k2"}
						c.UpdateWithHash("k3", 3)
						Expect(c.cache).To(HaveLen(3))
						Expect(c.recentKeys[2]).To(Equal("k3"))
					})
				})
			})
			Describe("func RemoveWithHash", func() {
				Context("When removing a hash", func() {
					It("Then both the cache and recentKeys is updated", func() {
						c := NewSimpleCache(3)
						c.cache["k1"] = 1
						c.cache["k2"] = 2
						c.recentKeys = []string{"k1", "k2"}
						c.RemoveWithHash("k2")
						Expect(c.cache).To(HaveLen(1))
						Expect(c.recentKeys[0]).To(Equal("k1"))
					})
				})
			})
			Describe("func addToRecentKeys", func() {
				Context("When adding a hash to recentkKeys that exist", func() {
					It("Then the recentKeys should be reorder to reflect the newly added", func() {
						c := NewSimpleCache(3)
						c.cache["k1"] = 1
						c.cache["k2"] = 2
						c.recentKeys = []string{"k1", "k2"}
						c.addToRecentKeys("k1")
						Expect(c.cache).To(HaveLen(2))
						Expect(c.recentKeys[0]).To(Equal("k2"))
						Expect(c.recentKeys[1]).To(Equal("k1"))
					})
				})
				Context("When adding a hash to recentkKeys that does not exist", func() {
					It("Then the key should be added", func() {
						c := NewSimpleCache(3)
						c.cache["k1"] = 1
						c.cache["k2"] = 2
						c.recentKeys = []string{"k1", "k2"}
						c.addToRecentKeys("k3")
						Expect(c.cache).To(HaveLen(2))
						Expect(c.recentKeys).To(HaveLen(3))
					})
				})
			})
			Describe("func checkCacheSize", func() {
				Context("When inserting and cache has reach size limit", func() {
					It("Then the oldest cache value should be removed and returned", func() {
						c := NewSimpleCache(1)
						c.cache["k2"] = 2
						c.cache["k1"] = 1
						c.recentKeys = []string{"k1", "k2"}
						updated, key, value := c.checkCacheSize()
						Expect(updated).To(BeTrue())
						Expect(key).To(Equal("k1"))
						Expect(value).To(Equal(1))
						Expect(c.cache).To(HaveLen(1))
						Expect(c.recentKeys).To(HaveLen(1))
						Expect(c.cache["k2"]).ToNot(BeNil())
						Expect(c.recentKeys[0]).To(Equal("k2"))
					})
				})
				Context("When inserting and cache has not reach size limit", func() {
					It("Then cache and recentKeys should not be touched", func() {
						c := NewSimpleCache(2)
						c.cache["k1"] = 1
						c.cache["k2"] = 2
						c.recentKeys = []string{"k1", "k2"}
						updated, _, _ := c.checkCacheSize()
						Expect(updated).To(BeFalse())
						Expect(c.cache).To(HaveLen(2))
						Expect(c.recentKeys).To(HaveLen(2))
					})
				})
			})
		})
	})
})
