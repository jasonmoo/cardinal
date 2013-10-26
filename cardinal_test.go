package cardinal

import (
	"math"
	"strconv"
	"testing"
	"time"
)

var (
	a    = []byte("a")
	b    = []byte("b")
	c    = New(time.Second, 10000)
	runs = 10
)

func TestReset(t *testing.T) {

	c.Add(a)

	c.Reset()

	if c.Count() != 0 || !math.IsNaN(c.Cardinality()) {
		t.Errorf("Reset failed got: %d, %f (count, card)", c.Count(), c.Cardinality())
	}

}

func TestCardinality(t *testing.T) {

	c.Reset()

	// each item unique == 1.0 every time
	for i := 0; i < runs; i++ {
		c.Add([]byte(strconv.Itoa(i)))
		if card := c.Cardinality(); card != 1 {
			t.Errorf("Expected cardinality: 1, got: %f", card)
		}
	}

	c.Reset()

	// 2 items, 1 uniq = .5
	c.Add(a)
	c.Add(a)

	if card := c.Cardinality(); card != 0.5 {
		t.Errorf("Expected cardinality: 0.5, got: %f", card)
	}

}

func BenchmarkAdd(b *testing.B) {

	c := New(time.Second, 1000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Add(a)
	}

}

func BenchmarkCheck(b *testing.B) {

	a, c := []byte{'a'}, New(time.Second, 1000)

	for i := 0; i < b.N; i++ {
		c.Add(a)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Check(a)
	}

}

func BenchmarkCardinality(b *testing.B) {

	a, c := []byte{'a'}, New(time.Second, 1000)

	for i := 0; i < 100000; i++ {
		c.Add(a)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Cardinality()
	}

}
