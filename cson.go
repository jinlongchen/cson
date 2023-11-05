package cson

import (
	jsonlib "encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/spf13/cast"
)

type JSON struct {
	locker *sync.RWMutex
	val    any
}

func NewJSON(val any) *JSON {
	return &JSON{val: val, locker: nil}
}

func (json *JSON) Safe() *JSON {
	if json.locker != nil {
		return json
	}
	return &JSON{
		locker: &sync.RWMutex{},
		val:    json.val,
	}
}

func (json *JSON) IsNil() bool {
	return json == nil || json.val == nil
}

func (json *JSON) Get(path string) *JSON {
	if json.val == nil {
		return &JSON{}
	}

	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}

	curr := json.val
	components := strings.Split(path, ".")
	for _, i := range components {
		if m, ok := curr.(map[string]any); ok {
			if v, ok := m[i]; ok {
				curr = v
			} else {
				return &JSON{}
			}
		} else {
			return &JSON{}
		}
	}

	return &JSON{val: curr, locker: json.locker}
}

func (json *JSON) Set(path string, v any) *JSON {
	if json.locker != nil {
		json.locker.Lock()
		defer json.locker.Unlock()
	}

	if path == "" {
		json.val = v
		return json
	}

	if json.val == nil {
		json.val = make(map[string]any)
	}

	switch json.val.(type) {
	case map[string]any:
	default:
		json.val = make(map[string]any)
	}

	keys := strings.Split(path, ".")
	curr := json.val.(map[string]any)
	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		next := curr[key]
		switch next := next.(type) {
		case map[string]any:
			curr = next
		default:
			newNext := make(map[string]any)
			curr[key] = newNext
			curr = newNext
		}
	}
	curr[keys[len(keys)-1]] = v
	return json
}

func (json *JSON) Value() any {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}

	if json.val == nil {
		return nil
	}
	switch v := json.val.(type) {
	case JSON:
		return v.Value()
	default:
		return v
	}
}

func (json *JSON) Val() any {
	return json.Value()
}

func (json *JSON) String() string {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}

	if json.val == nil {
		return ""
	}
	switch v := json.val.(type) {
	case string:
		return v
	case JSON:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
func (json *JSON) Str() string {
	return json.String()
}

func (json *JSON) Float64() float64 {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}
	if json.val == nil {
		return 0
	}
	switch v := json.val.(type) {
	case JSON:
		return v.Float64()
	default:
		return cast.ToFloat64(v)
	}
}

func (json *JSON) Int64() int64 {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}
	if json.val == nil {
		return 0
	}
	switch v := json.val.(type) {
	case JSON:
		return v.Int64()
	default:
		return cast.ToInt64(v)
	}
}

func (json *JSON) Bool() bool {
	if json.val == nil {
		return false
	}
	switch v := json.val.(type) {
	case JSON:
		return v.Bool()
	default:
		return cast.ToBool(v)
	}
}

func (json *JSON) Slice() []*JSON {
	res := make([]*JSON, 0)
	if json.val == nil {
		return res
	}
	var val any
	switch v := json.val.(type) {
	case JSON:
		val = v.Val()
	default:
		val = v
	}
	slice := cast.ToSlice(val)
	for _, item := range slice {
		res = append(res, NewJSON(item))
	}
	return res
}
func (json *JSON) Eq(a any) bool {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}
	return reflect.DeepEqual(json.val, a)
}

func (json *JSON) MarshalSON() ([]byte, error) {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}

	res, err := jsonlib.Marshal(json.val)
	return res, err
}

func (json *JSON) UnmarshalJSON(data []byte) error {
	v := new(any)
	err := jsonlib.Unmarshal(data, v)
	if err != nil {
		return err
	}
	json.val = v
	return nil
}
