package cardinal

import (
	"sync"
	"time"

	"github.com/dataence/bloom"
	"github.com/dataence/bloom/scalable"
)

type (
	Cardinal struct {
		sync.Mutex

		buf []*Filter

		duration       time.Duration
		chunk_duration time.Duration
		last_t         time.Time
		i              int

		filter *Filter
	}

	Filter struct {
		bloom.Bloom
		uniques uint64
		total   uint64
	}
)

const (
	CHUNKS     = 10
	CHUNK_SIZE = 4096
)

func New(duration time.Duration) *Cardinal {

	buf := make([]*Filter, CHUNKS)

	// initialize with modest size to ensure
	for i, _ := range buf {
		buf[i] = &Filter{scalable.New(CHUNK_SIZE), 0, 0}
	}

	return &Cardinal{
		buf:            buf,
		filter:         buf[0],
		duration:       duration,
		chunk_duration: duration / CHUNKS,
	}

}

func (c *Cardinal) Add(token []byte) {

	c.Lock()

	t := time.Now().Truncate(c.chunk_duration)

	if !t.Equal(c.last_t) {
		c.last_t = t
		c.i++
		next_i := c.i % len(c.buf)
		// always create a new filter with the size of the previous
		// as the estimated number of items to minimize collisions
		c.buf[next_i] = &Filter{scalable.New(min(CHUNK_SIZE, c.filter.Count())), 0, 0}
		c.filter = c.buf[next_i]
	}

	// check all filters before incrementing
	if !c.check(token) {
		c.filter.Add(token)
		c.filter.uniques++
	}

	c.filter.total++

	c.Unlock()

}

func (c *Cardinal) Check(token []byte) (r bool) {
	c.Lock()
	r = c.check(token)
	c.Unlock()
	return
}

func (c *Cardinal) Cardinality() (r float64) {
	c.Lock()
	r = c.cardinality()
	c.Unlock()
	return
}

func (c *Cardinal) Count() (r uint64) {
	c.Lock()
	r = c.count()
	c.Unlock()
	return
}

func (c *Cardinal) Uniques() (r uint64) {
	c.Lock()
	r = c.uniques()
	c.Unlock()
	return
}

func (c *Cardinal) Duration() time.Duration {
	return c.duration
}

func (c *Cardinal) Reset() {
	c.Lock()
	for _, filter := range c.buf {
		filter.reset()
	}
	c.Unlock()
}

func (c *Cardinal) check(token []byte) bool {

	for _, filter := range c.buf {
		if filter.Check(token) {
			return true
		}
	}

	return false

}

func (c *Cardinal) cardinality() float64 {

	var uniques, total uint64

	for _, filter := range c.buf {
		uniques, total = uniques+filter.uniques, total+filter.total
	}

	return float64(uniques) / float64(total)

}

func (c *Cardinal) count() (total uint64) {

	for _, filter := range c.buf {
		total += filter.total
	}

	return

}

func (c *Cardinal) uniques() (uniques uint64) {

	for _, filter := range c.buf {
		uniques += filter.uniques
	}

	return

}

func (f *Filter) reset() {
	f.Bloom.Reset()
	f.uniques = 0
	f.total = 0
}

func min(a, b uint) uint {
	if a < b {
		return b
	}
	return a
}
