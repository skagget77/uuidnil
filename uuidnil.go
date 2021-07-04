// Package uuidnil allow github.com/google/uuid.UUIDs to be unmarshaled from
// empty and malformed JSON strings.
package uuidnil

import (
	"bytes"
	"encoding/json"
	"log"
	"reflect"
)

// Option for the Wrap function.
type Option int

func (o Option) trace() bool {
	return o&TraceLog != 0
}

func (o Option) debug() bool {
	return o&DebugLog != 0
}

func (o Option) empty() bool {
	return o&AllowEmpty != 0
}

func (o Option) invalid() bool {
	return o&AllowInvalid != 0
}

const (
	// TraceLog signals the wrapper to output trace information. Implies
	// DebugLog.
	TraceLog Option = 1

	// DebugLog signals the wrapper to output debug information. Less verbose
	// than the trace log.
	DebugLog Option = 2

	// AllowEmpty signals the wrapper to allow empty strings to be unmarshaled
	// to uuid.Nil.
	AllowEmpty Option = 4

	// AllowInvalid signals the wrapper to allow invalid UUIDs to be unmarshaled
	// to uuid.Nil. Implies AllowEmpty.
	AllowInvalid Option = 8
)

// Wrap wraps v in a wrapper that allow for more lenient unmarshaling of
// uuid.UUID values from JSON data. The options specifies how leanient the
// unmarshaler should be.
func Wrap(v interface{}, options ...Option) *wrapper {
	var opts Option
	for _, option := range options {
		opts |= option
	}

	return &wrapper{val: v, opts: opts}
}

type wrapper struct {
	val  interface{}
	opts Option
}

// UnmarshalJSON unmarshals the JSON data to a proxy value which is then copied
// to the wrapped value.
func (w *wrapper) UnmarshalJSON(data []byte) error {
	if w.opts.trace() {
		log.Println("[TRACE] UnmarshalJSON")
	}

	value := reflect.ValueOf(w.val)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return &json.InvalidUnmarshalError{Type: value.Type()}
	}

	var target reflect.Value
	var assign assignFunc
	if w.opts.empty() || w.opts.invalid() {
		target = value
		value, assign = proxyValue(value, w.opts)

		if err := assign(reflect.Indirect(value), reflect.Indirect(target)); err != nil {
			return err
		}
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(value.Interface()); err != nil {
		return err
	}

	if w.opts.empty() || w.opts.invalid() {
		if err := assign(reflect.Indirect(target), reflect.Indirect(value)); err != nil {
			return err
		}
	}

	return nil
}
