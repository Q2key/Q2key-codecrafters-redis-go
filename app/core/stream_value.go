package core

import (
	"fmt"
	"strings"
	"time"
)

type StreamValue struct {
	Value     map[float64][]int
	HMap      map[int]map[string]string
	LSTime    float64
	LSSeqn    int
	Expired   *time.Time
	ValueType ValueType
}

type StreamDescriptor struct {
	id int
}

func NewStreamValue(lsTime float64) *StreamValue {
	var seq int = 0
	if lsTime == 0 {
		seq = 1
	}

	return &StreamValue{
		Value:     map[float64][]int{},
		HMap:      map[int]map[string]string{},
		LSTime:    lsTime,
		LSSeqn:    seq,
		ValueType: STREAM,
	}
}

func (r *StreamValue) IsExpired() bool {
	if r.Expired == nil {
		return false
	}

	return r.Expired.UnixNano() <= time.Now().UTC().UnixNano()
}

func (r *StreamValue) ToString() string {
	return fmt.Sprintf("%.0f-%d", r.LSTime, r.LSSeqn)
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

func (r *StreamValue) UpdateSeqKey(newTime float64) int {
	v, ok := r.Value[newTime]

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

func (r *StreamValue) WriteSequence(newTime float64, seq int, payload string) {
	v, ok := r.Value[newTime]

	if !ok {
		r.Value[newTime] = []int{seq}
	} else {
		r.Value[newTime] = append(v, seq)
	}

	parts := strings.Split(payload, ":")

	_, ok = r.HMap[seq]
	if !ok {
		r.HMap[seq] = map[string]string{
			parts[0]: parts[1],
		}
	} else {
		r.HMap[seq][parts[0]] = parts[1]
	}

	fmt.Println(r.HMap)
	r.LSSeqn = seq
	r.LSTime = newTime
}

func (r *StreamValue) CanSave(newTime float64, newSequence int) (bool, *string) {
	if newTime == r.LSTime && newSequence > r.LSSeqn {
		return true, nil
	}

	if newTime == 0 && newSequence == 0 {
		mess := "ERR The ID specified in XADD must be greater than 0-0"
		return false, &mess
	}

	if newTime == r.LSTime && r.LSSeqn >= newSequence {
		mess := "ERR The ID specified in XADD is equal or smaller than the target stream top item"
		return false, &mess
	}

	return false, nil
}
