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
	if json == nil || json.val == nil {
		return true
	}
	switch v := json.val.(type) {
	case JSON:
		return v.IsNil()
	case *JSON:
		return v.IsNil()
	default:
		return false
	}
}

func (json *JSON) Get(path string) *JSON {
	if json.IsNil() {
		return &JSON{}
	}

	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}

	if len(path) == 0 {
		return json
	}

	curr := json.val

	components := strings.Split(path, ".")
	for _, i := range components {
		switch t := curr.(type) {
		case *JSON:
			curr = t.Get(i)
		case JSON:
			curr = t.Get(i)
		case any:
			if m, ok := t.(map[string]any); ok {
				if v, ok := m[i]; ok {
					curr = v
				} else {
					return &JSON{}
				}
			} else {
				return &JSON{}
			}
		case map[string]any:
			if v, ok := t[i]; ok {
				curr = v
			} else {
				return &JSON{}
			}
		default:
			return &JSON{}
		}
	}

	switch curr := curr.(type) {
	case *JSON:
		return curr
	case JSON:
		return &curr
	default:
		return &JSON{val: curr, locker: json.locker}
	}
}

func (json *JSON) Set(path string, v any) *JSON {
	if json.locker != nil {
		json.locker.Lock()
		defer json.locker.Unlock()
	}

	if json == v {
		panic("")
	}

	if path == "" {
		json.val = v
		return json
	}

	keys := strings.Split(path, ".")
	var curr map[string]any

	if json.IsNil() {
		curr = make(map[string]any)
		json.val = curr
	} else {
		switch v := json.val.(type) {
		case map[string]any:
			curr = v
		default:
			curr = make(map[string]any)
			json.val = curr
		}
	}

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

	if json.IsNil() {
		return nil
	}
	switch v := json.val.(type) {
	case JSON:
		return v.Value()
	case *JSON:
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

	if json.IsNil() {
		return ""
	}
	switch v := json.Val().(type) {
	case string:
		return v
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
	if json.IsNil() {
		return 0
	}
	return cast.ToFloat64(json.Val())
}

func (json *JSON) Int64() int64 {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}
	if json.IsNil() {
		return 0
	}
	return cast.ToInt64(json.Val())
}
func (json *JSON) Int() int {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}
	if json.IsNil() {
		return 0
	}
	return cast.ToInt(json.Val())
}

func (json *JSON) Bool() bool {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}
	if json.IsNil() {
		return false
	}
	return cast.ToBool(json.Val())
}

func (json *JSON) Slice() []*JSON {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}
	res := make([]*JSON, 0)
	if json.IsNil() {
		return res
	}
	slice := cast.ToSlice(json.Val())
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

func (json JSON) MarshalJSON() ([]byte, error) {
	if json.locker != nil {
		json.locker.RLock()
		defer json.locker.RUnlock()
	}

	res, err := jsonlib.Marshal(json.Val())
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
