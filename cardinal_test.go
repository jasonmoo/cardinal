package cardinal

import (
	"math"
	"strconv"
	"testing"
	"time"
)

var (
	a    = "a"
	b    = "b"
	c    = New(time.Second)
	runs = 10
)

func TestAddAndCount(t *testing.T) {

	if ct := c.Count(); ct != 0 {
		t.Errorf("Expected count: 0, got: %d", ct)
	}

	if ct := c.Uniques(); ct != 0 {
		t.Errorf("Expected uniques: 0, got: %d", ct)
	}

	c.Add(a)

	if ct := c.Count(); ct != 1 {
		t.Errorf("Expected count: 1, got: %d", ct)
	}

	if ct := c.Uniques(); ct != 1 {
		t.Errorf("Expected uniques: 1, got: %d", ct)
	}

	c.Add(a)

	if ct := c.Count(); ct != 2 {
		t.Errorf("Expected count: 2, got: %d", ct)
	}

	if ct := c.Uniques(); ct != 1 {
		t.Errorf("Expected uniques: 1, got: %d", ct)
	}

}

func TestCheck(t *testing.T) {

	// missing value
	if c.Check(b) {
		t.Errorf("Expected false for missing value Check()")
	}

	// present value
	if !c.Check(a) {
		t.Errorf("Expected true for present value Check()")
	}

}

func TestReset(t *testing.T) {

	c.Reset()

	switch {
	case c.Count() != 0, c.Uniques() != 0, !math.IsNaN(c.Cardinality()):
		t.Errorf("Reset failed got: %d, %d, %f (count, uniques, card)", c.Count(), c.Uniques(), c.Cardinality())
	}

}

func TestCardinality(t *testing.T) {

	c.Reset()

	// each item unique == 1.0 every time
	for i := 0; i < runs; i++ {
		c.Add(strconv.Itoa(i))
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

	c := New(time.Second)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Add(a)
	}

}

func BenchmarkCheck(b *testing.B) {

	a, c := "a", New(time.Second)
	c.Add(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Check(a)
	}

}

func BenchmarkCardinality(b *testing.B) {

	a, c := "a", New(time.Second)

	for i := 0; i < 100000; i++ {
		c.Add(a)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Cardinality()
	}

}
