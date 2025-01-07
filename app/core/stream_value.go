package core

import "time"

type StreamValue struct {
	Value     map[string]interface{}
	Expired   *time.Time
	ValueType ValueType
}

func (r *StreamValue) IsExpired() bool {
	if r.Expired == nil {
		return false
	}

	return r.Expired.UnixNano() <= time.Now().UTC().UnixNano()
}

func (r *StreamValue) ToString() string {
	return "Not impl"
}

func (r *StreamValue) GetType() string {
	return "stream"
}

func (r *StreamValue) SetExpired(expired time.Time) {
	r.Expired = &expired
}

func (r *StreamValue) SetValue(value interface{}) {
	val, ok := value.(map[string]interface{})
	if ok {
		r.Value = val
	}
}
