package cardinal

import (
	"github.com/jasonmoo/bloom"
	"github.com/jasonmoo/bloom/scalable"
	"sync"
	"time"
)

type (
	Cardinal struct {
		sync.Mutex
		buf            []*Filter
		chunk_duration time.Duration

		last_i int
	}

	Filter struct {
		bloom.Bloom
		uniques uint
	}
)

func New(duration time.Duration, n int) *Cardinal {

	const (
		chunks = 10
		min_n  = 1000
	)

	buf, chunk_n := make([]*Filter, chunks), uint(max(min_n, n/chunks))

	for i, _ := range buf {
		buf[i] = &Filter{scalable.New(chunk_n), 0}
	}

	return &Cardinal{
		buf:            buf,
		chunk_duration: duration / chunks,
	}

}

func (c *Cardinal) Add(token []byte) {

	c.Lock()

	i := int(time.Now().Truncate(c.chunk_duration).UnixNano() % int64(len(c.buf)))

	filter := c.buf[i]

	if c.last_i != i {
		c.last_i = i
		filter.Reset()
		filter.uniques = 0
	}

	if !filter.Check(token) {
		filter.uniques++
	}

	filter.Add(token)

	c.Unlock()

}

func (c *Cardinal) Check(token []byte) bool {

	c.Lock()
	defer c.Unlock()

	for _, filter := range c.buf {
		if filter.Check(token) {
			return true
		}
	}
	return false

}

func (c *Cardinal) Cardinality() float64 {

	c.Lock()
	defer c.Unlock()

	var uniques, total uint

	for _, filter := range c.buf {
		uniques, total = uniques+filter.uniques, total+filter.Count()
	}

	return float64(uniques) / float64(total)

}

func (c *Cardinal) Count() (total uint) {

	c.Lock()

	for _, filter := range c.buf {
		total += filter.Count()
	}

	c.Unlock()

	return

}

func (c *Cardinal) Reset() {

	c.Lock()

	for _, filter := range c.buf {
		filter.Reset()
		filter.uniques = 0
	}

	c.Unlock()

}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
