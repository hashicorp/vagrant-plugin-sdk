package types

import (
	"fmt"
	"sync"
)

type EndOfRange struct {
	message string
}

func (e *EndOfRange) Error() string {
	return e.message
}

type Range interface {
	Current() int32
	Final() int32
	Full() []int32
	Initial() int32
	Next() (int32, error)
}

func NewRange(start, end int32) (Range, error) {
	if start > end {
		return nil, fmt.Errorf("invalid range, start value is greater than end value")
	}
	return &grange{
		start:   start,
		end:     end,
		current: start,
		m:       sync.Mutex{},
	}, nil
}

func endOfRange(msg string, args ...interface{}) error {
	return &EndOfRange{
		message: fmt.Sprintf(msg, args...),
	}
}

type grange struct {
	start int32
	end   int32

	current int32
	m       sync.Mutex
}

func (g *grange) Current() int32 {
	return g.current
}

func (g *grange) Next() (int32, error) {
	g.m.Lock()
	defer g.m.Unlock()

	if g.current == g.end {
		return g.current, endOfRange("range ends at %d", g.end)
	}

	g.current += 1
	return g.current, nil
}

func (g *grange) Full() []int32 {
	length := (g.end - g.start) + 1
	result := make([]int32, length)
	for i := 0; int32(i) < length; i++ {
		result[i] = g.start + int32(i)
	}
	return result
}

func (g *grange) Initial() int32 {
	return g.start
}

func (g *grange) Final() int32 {
	return g.end
}

var _ Range = (*grange)(nil)
