// +build all benchmark
package registry

import (
	"math/rand"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	getCacheSize   = 1000
	getEntryCount  = 5000 // How many entries inserted to registry for get benches
	getRepeatCount = 1000 // How many times to repeat get for each benchmark
)

var _ = Describe("Registry", func() {
	Describe("Get", func() {
		Context("Simple Registry", func() {
			var r Registry
			var k []Key
			It("Setup registry", func() {
				r = NewSimpleRegistry()
				k = insertRandomEntries(r, getEntryCount)
			})
			Measure("Getting the same entry", func(b Benchmarker) {
				benchGetSameKey(r, k, b)
			}, 10)
			Measure("Getting a random entry", func(b Benchmarker) {
				benchGetRandomKey(r, k, b)
			}, 10)
		})
		Context("Better Registry", func() {
			var r Registry
			var k []Key
			It("Setup registry", func() {
				r = NewBetterRegistry()
				k = insertRandomEntries(r, getEntryCount)
			})
			Measure("Getting the same entry", func(b Benchmarker) {
				benchGetSameKey(r, k, b)
			}, 10)
			Measure("Getting a random entry", func(b Benchmarker) {
				benchGetRandomKey(r, k, b)
			}, 10)
		})
		Context("Cached Registry", func() {
			var r Registry
			var k []Key
			It("Setup registry", func() {
				r = NewCacheRegistry(getCacheSize)
				k = insertRandomEntries(r, getEntryCount)
			})
			Measure("Getting the same entry", func(b Benchmarker) {
				benchGetSameKey(r, k, b)
			}, 10)
			Measure("Getting a random entry", func(b Benchmarker) {
				benchGetRandomKey(r, k, b)
			}, 10)
		})
	})
})

func insertRandomEntries(r Registry, count int) []Key {
	k := []Key{}
	for i := 0; i < count; i++ {
		k = append(k, Key{"a" + strconv.Itoa(i): "b", "c": "d" + strconv.Itoa(i)})
		r.Set(k[i], i)
	}
	return k
}

func benchGetSameKey(r Registry, k []Key, b Benchmarker) time.Duration {
	return b.Time("runtime", func() {
		n := len(k) - 1
		for j := 0; j < getRepeatCount; j++ {
			i := r.Get(k[n])
			Expect(i).To(Equal(n))
		}
	})
}

func benchGetRandomKey(r Registry, k []Key, b Benchmarker) time.Duration {
	return b.Time("runtime", func() {
		maxEntryIndex := len(k) - 1
		for j := 0; j < getRepeatCount; j++ {
			n := rand.Intn(maxEntryIndex)
			i := r.Get(k[n])
			Expect(i).To(Equal(n))
		}
	})
}
