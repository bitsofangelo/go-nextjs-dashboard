package optional

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Optional carries three states: missing, explicitly null, or a value.
type Optional[T any] struct {
	Val       T
	IsPresent bool
	IsNull    bool
}

// UnmarshalJSON marks IsPresent=true, then
//   - if data == "null" → IsNull=true
//   - else → unmarshal into Val
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.IsPresent = true
	if bytes.Equal(data, []byte("null")) {
		o.IsNull = true
		// zero value
		var zero T
		o.Val = zero
		return nil
	}
	o.IsNull = false
	return json.Unmarshal(data, &o.Val)
}

// MarshalJSON emits either "null" (missing or null) or the JSON of Val.
func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.IsPresent || o.IsNull {
		return []byte("null"), nil
	}
	return json.Marshal(o.Val)
}

// Value implements driver.Valuer for database writes.
//   - nil → NULL
//   - otherwise → underlying Val (must be a type driver accepts)
func (o Optional[T]) Value() (driver.Value, error) {
	if o.IsNull {
		return nil, nil
	}

	v := any(o.Val)

	if u, ok := v.(uuid.UUID); ok {
		return u.String(), nil
	}

	if dv, ok := v.(driver.Value); ok {
		return dv, nil
	}

	switch tv := v.(type) {
	case int64, float64, bool, string, []byte:
		return tv, nil
	default:
		return nil, fmt.Errorf("optional: unsupported type %T for driver.Value", o.Val)
	}
}

// Scan implements sql.Scanner for database reads.
//   - src == nil → IsNull=true
//   - else → try to scan into T (either via T’s own Scanner or via a direct cast)
func (o *Optional[T]) Scan(src any) error {
	o.IsPresent = true
	if src == nil {
		o.IsNull = true
		fmt.Println("here")
		var zero T
		o.Val = zero
		return nil
	}
	o.IsNull = false

	var tmp T

	// If T itself implements sql.Scanner, use that:
	if scanner, ok := any(&tmp).(sql.Scanner); ok {
		err := scanner.Scan(src)
		if err != nil {
			return fmt.Errorf("optional scan: %w", err)
		}
		o.Val = tmp
		return nil
	}

	// Otherwise try a direct conversion:
	if v, ok := src.(T); ok {
		o.Val = v
		return nil
	}

	// Handle boolean
	if _, ok := any(tmp).(bool); ok {
		switch v := src.(type) {
		case int64:
			o.Val = any(v != 0).(T)
			return nil
		case []byte:
			o.Val = any(string(v) == "1").(T)
			return nil
		case string:
			o.Val = any(v == "1").(T)
			return nil
		}
	}

	// Special‐case: many drivers return []byte for text columns
	if b, ok := src.([]byte); ok {
		// only works if T is e.g. string or []byte
		switch any(o.Val).(type) {
		case string:
			o.Val = any(string(b)).(T)
			return nil
		case []byte:
			o.Val = any(b).(T)
			return nil
		}
	}

	return fmt.Errorf(
		"cannot scan %T into Optional[%s]",
		src, reflect.TypeOf(o.Val).Name(),
	)
}

func (o Optional[T]) IsZero() bool {
	return !o.IsPresent
}

func (o *Optional[T]) Ptr() *T {
	if o.IsPresent && !o.IsNull {
		return &o.Val
	}
	return nil
}

func Of[T any](v T) Optional[T] {
	return Optional[T]{Val: v, IsPresent: true, IsNull: false}
}

func FromPtr[T any](v *T) Optional[T] {
	var (
		n      T
		isNull bool
	)

	if v != nil {
		n = *v
	} else {
		isNull = true
	}

	return Optional[T]{
		Val:       n,
		IsPresent: true,
		IsNull:    isNull,
	}
}

func StringToUUID(o Optional[string]) (Optional[uuid.UUID], error) {
	if !o.IsPresent {
		return Optional[uuid.UUID]{}, nil
	}
	if o.IsNull {
		return Optional[uuid.UUID]{IsPresent: true, IsNull: true}, nil
	}
	u, err := uuid.Parse(o.Val)
	if err != nil {
		return Optional[uuid.UUID]{}, fmt.Errorf("cannot parse %s as UUID: %w", o.Val, err)
	}
	return Optional[uuid.UUID]{Val: u, IsPresent: true}, nil
}

func StringToTime(o Optional[string], layout string) (Optional[time.Time], error) {
	if !o.IsPresent {
		return Optional[time.Time]{}, nil
	}
	if o.IsNull {
		return Optional[time.Time]{IsPresent: true, IsNull: true}, nil
	}
	t, err := time.Parse(layout, o.Val)
	if err != nil {
		return Optional[time.Time]{}, fmt.Errorf("cannot parse %v as time: %w", o.Val, err)
	}
	return Optional[time.Time]{Val: t, IsPresent: true}, nil
}
