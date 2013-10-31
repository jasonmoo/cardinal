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

		buf []*Filter

		duration       time.Duration
		chunk_duration time.Duration
		last_t         time.Time
		i              uint

		filter *Filter
	}

	Filter struct {
		bloom.Bloom
		uniques uint
		total   uint
	}
)

func New(duration time.Duration) *Cardinal {

	const (
		chunks          = 10
		chunk_size uint = 4096
	)

	buf := make([]*Filter, chunks)

	for i, _ := range buf {
		buf[i] = &Filter{scalable.New(chunk_size), 0, 0}
	}

	return &Cardinal{
		buf:            buf,
		duration:       duration,
		chunk_duration: duration / chunks,
	}

}

func (c *Cardinal) Add(token []byte) {

	c.Lock()

	t := time.Now().Truncate(c.chunk_duration)

	if !t.Equal(c.last_t) {
		c.last_t = t
		c.i++
		c.filter = c.buf[c.i%uint(len(c.buf))]
		c.filter.reset()
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

func (c *Cardinal) Count() (r uint) {
	c.Lock()
	r = c.count()
	c.Unlock()
	return
}

func (c *Cardinal) Uniques() (r uint) {
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

	var uniques, total uint

	for _, filter := range c.buf {
		uniques, total = uniques+filter.uniques, total+filter.total
	}

	return float64(uniques) / float64(total)

}

func (c *Cardinal) count() (total uint) {

	for _, filter := range c.buf {
		total += filter.total
	}

	return

}

func (c *Cardinal) uniques() (uniques uint) {

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
