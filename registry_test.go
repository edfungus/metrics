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
	Describe("Even Better registry", func() {
		Describe("Given a user wants to get a value", func() {
			Context("When the key exist", func() {
				var r *CachedRegistry
				var c *SimpleCache
				It("Then the value should be returned", func() {
					r = NewCacheRegistry()
					c = NewSimpleCache(1)
					r.getCache = c
					k := map[string]string{"k1": "v1", "k2": "v2"}
					he := &hashEntry{
						keys:  k,
						value: 1,
					}
					r.registry.registry["k1"] = map[string][]*hashEntry{}
					r.registry.registry["k1"]["v1"] = []*hashEntry{he}
					r.registry.registry["k2"] = map[string][]*hashEntry{}
					r.registry.registry["k2"]["v2"] = []*hashEntry{he}

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
				It("Then the value should be returned", func() {
					r := NewCacheRegistry()
					k := map[string]string{"k1": "v1", "k2": "v2"}
					he := &hashEntry{
						keys:  k,
						value: 1,
					}
					r.registry.registry["k1"] = map[string][]*hashEntry{}
					r.registry.registry["k1"]["v1"] = []*hashEntry{he}
					r.registry.registry["k2"] = map[string][]*hashEntry{}
					r.registry.registry["k2"]["v2"] = []*hashEntry{he}

					newK := map[string]string{"k1": "v1"}
					Expect(r.Get(newK)).To(BeNil())
				})
			})
			Context("When the key has already been retrieved", func() {
				var r *CachedRegistry
				var c *SimpleCache
				It("Then the value should be returned from cache", func() {
					r = NewCacheRegistry()
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
					r.registry.registry["k1"] = map[string][]*hashEntry{}
					r.registry.registry["k1"]["v1"] = []*hashEntry{he1}
					r.registry.registry["k2"] = map[string][]*hashEntry{}
					r.registry.registry["k2"]["v2"] = []*hashEntry{he1}

					Expect(r.Get(k)).To(Equal(1))

					he1.value = 2 // set new value
					// hack registry so that, if this value is returned, we know the cache was not used
					r.registry.registry["k1"]["v1"] = []*hashEntry{he2}
					r.registry.registry["k2"]["v2"] = []*hashEntry{he2}
					Expect(r.Get(k)).To(Equal(2))

					// force clear cache so we get value from registry
					c.removeWithHash(toHashString(k))
					Expect(r.Get(k)).To(Equal(3))
				})
				It("Then cache should not have grown", func() {
					Expect(c.recentKeys).To(HaveLen(1))
					Expect(c.cache).To(HaveLen(1))
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
			Describe("func removeWithHash", func() {
				Context("When removing a hash", func() {
					It("Then both the cache and recentKeys is updated", func() {
						c := NewSimpleCache(3)
						c.cache["k1"] = 1
						c.cache["k2"] = 2
						c.recentKeys = []string{"k1", "k2"}
						c.removeWithHash("k2")
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
