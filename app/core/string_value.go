package core

import "time"

type StringValue struct {
	Value     string
	Expired   *time.Time
	ValueType ValueType
}

func (r *StringValue) IsExpired() bool {
	if r.Expired == nil {
		return false
	}

	return r.Expired.UnixNano() <= time.Now().UTC().UnixNano()
}

func (r *StringValue) ToString() string {
	return r.Value
}

func (r *StringValue) GetType() string {
	return "string"
}

func (r *StringValue) SetExpired(expired time.Time) {
	r.Expired = &expired
}

func (r *StringValue) GetValue() string {
	return r.Value
}

func (r *StringValue) SetValue(value interface{}) {
	val, ok := value.(string)
	if ok {
		r.Value = val
	}
}
