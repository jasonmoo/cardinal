package cardinal

import (
	"github.com/jasonmoo/bloom"
	"github.com/jasonmoo/bloom/scalable"
	"time"
)

type (
	Cardinal struct {
		buf            []*Filter
		chunk_duration time.Duration

		last_i int
	}

	Filter struct {
		bloom.Bloom
		uniques uint
	}
)

func New(duration time.Duration) *Cardinal {

	const (
		chunks          = 10
		chunk_size uint = 4096
	)

	buf := make([]*Filter, chunks)

	for i, _ := range buf {
		buf[i] = &Filter{scalable.New(chunk_size), 0}
	}

	return &Cardinal{
		buf:            buf,
		chunk_duration: duration / chunks,
	}

}

func (c *Cardinal) Add(token []byte) {

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

}

func (c *Cardinal) Check(token []byte) bool {

	for _, filter := range c.buf {
		if filter.Check(token) {
			return true
		}
	}

	return false

}

func (c *Cardinal) Cardinality() float64 {

	var uniques, total uint

	for _, filter := range c.buf {
		uniques, total = uniques+filter.uniques, total+filter.Count()
	}

	return float64(uniques) / float64(total)

}

func (c *Cardinal) Uniques() (total uint) {

	for _, filter := range c.buf {
		total += filter.uniques
	}

	return

}

func (c *Cardinal) Count() (total uint) {

	for _, filter := range c.buf {
		total += filter.Count()
	}

	return

}

func (c *Cardinal) Reset() {

	for _, filter := range c.buf {
		filter.Reset()
		filter.uniques = 0
	}

}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
