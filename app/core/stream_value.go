package core

import (
	"fmt"
	"strings"
	"time"
)

type SequenceKeysByTsMark map[float64][]int

type (
	PairsBySequenceIdx map[string][]string
	StreamValue        struct {
		Value      SequenceKeysByTsMark
		Paris      PairsBySequenceIdx
		lastTsMark float64
		LastSidx   int
		Expired    *time.Time
		ValueType  ValueType
	}
)

type StreamDescriptor struct {
	id int
}

func NewStreamValue(tsMark float64) *StreamValue {
	var seq int = 0
	if tsMark == 0 {
		seq = 1
	}

	return &StreamValue{
		Value:      SequenceKeysByTsMark{},
		Paris:      PairsBySequenceIdx{},
		lastTsMark: tsMark,
		LastSidx:   seq,
		ValueType:  STREAM,
	}
}

func (r *StreamValue) IsExpired() bool {
	if r.Expired == nil {
		return false
	}

	return r.Expired.UnixNano() <= time.Now().UTC().UnixNano()
}

func (r *StreamValue) ToString() string {
	return fmt.Sprintf("%.0f-%d", r.lastTsMark, r.LastSidx)
}

func (r *StreamValue) GetType() string {
	return "stream"
}

func (r *StreamValue) SetExpired(expired time.Time) {
	r.Expired = &expired
}

func (r *StreamValue) SetValue(value interface{}) {
}

func (r *StreamValue) KeyExists(newTime float64) bool {
	_, ok := r.Value[newTime]
	return ok
}

func (r *StreamValue) UpdateSeqKey(tsMark float64) int {
	v, ok := r.Value[tsMark]

	if !ok {
		return 1
	}

	l := len(v)
	s := v[l-1]

	if l > 1 {
		return s + 1
	}

	return 0
}

func (r *StreamValue) WriteSequence(tsMark float64, sidx int, payload string) {
	v, ok := r.Value[tsMark]

	if !ok {
		r.Value[tsMark] = []int{sidx}
	} else {
		r.Value[tsMark] = append(v, sidx)
	}

	parts := strings.Split(payload, ":")

	r.Paris[formKey(tsMark, sidx)] = []string{parts[0], parts[1]}

	r.LastSidx = sidx
	r.lastTsMark = tsMark
}

func (r *StreamValue) CanSave(tsMark float64, sequence int) (bool, *string) {
	if tsMark == r.lastTsMark && sequence > r.LastSidx {
		return true, nil
	}

	if tsMark == 0 && sequence == 0 {
		mess := "ERR The ID specified in XADD must be greater than 0-0"
		return false, &mess
	}

	if tsMark == r.lastTsMark && r.LastSidx >= sequence {
		mess := "ERR The ID specified in XADD is equal or smaller than the target stream top item"
		return false, &mess
	}

	return false, nil
}

func formKey(ts float64, idx int) string {
	return fmt.Sprintf("%.0f-%d", ts, idx)
}
