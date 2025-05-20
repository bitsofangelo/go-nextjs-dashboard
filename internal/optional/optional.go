package optional

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

// Optional carries three states: missing, explicitly null, or a value.
type Optional[T any] struct {
	V         T
	IsPresent bool
	IsNull    bool
}

// UnmarshalJSON marks IsPresent=true, then
//   - if data == "null" → IsNull=true
//   - else → unmarshal into V
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.IsPresent = true
	if bytes.Equal(data, []byte("null")) {
		o.IsNull = true
		// zero value
		var zero T
		o.V = zero
		return nil
	}
	o.IsNull = false
	return json.Unmarshal(data, &o.V)
}

// MarshalJSON emits either "null" (missing or null) or the JSON of V.
func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.IsPresent || o.IsNull {
		return []byte("null"), nil
	}
	return json.Marshal(o.V)
}

// Value implements driver.Valuer for database writes.
//   - nil → NULL
//   - otherwise → underlying V (must be a type driver accepts)
func (o Optional[T]) Value() (driver.Value, error) {
	if o.IsNull {
		return nil, nil
	}

	v := any(o.V)

	if u, ok := v.(uuid.UUID); ok {
		return u.String(), nil
	}

	if dv, ok := v.(driver.Value); ok {
		return dv, nil
	}

	switch tv := any(o.V).(type) {
	case int64, float64, bool, string, []byte:
		return tv, nil
	default:
		return nil, fmt.Errorf("optional: unsupported type %T for driver.Value", o.V)
	}
}

// Scan implements sql.Scanner for database reads.
//   - src == nil → IsNull=true
//   - else → try to scan into T (either via T’s own Scanner or via a direct cast)
func (o *Optional[T]) Scan(src any) error {
	o.IsPresent = true
	if src == nil {
		o.IsNull = true
		var zero T
		o.V = zero
		return nil
	}
	o.IsNull = false

	// If T itself implements sql.Scanner, use that:
	var tmp T
	if scanner, ok := any(&tmp).(sql.Scanner); ok {
		err := scanner.Scan(src)
		if err != nil {
			return fmt.Errorf("optional scan: %w", err)
		}
		o.V = tmp
		return nil
	}

	// Otherwise try a direct conversion:
	if v, ok := src.(T); ok {
		o.V = v
		return nil
	}

	// Special‐case: many drivers return []byte for text columns
	if b, ok := src.([]byte); ok {
		// only works if T is e.g. string or []byte
		switch any(o.V).(type) {
		case string:
			o.V = any(string(b)).(T)
			return nil
		case []byte:
			o.V = any(b).(T)
			return nil
		}
	}

	return fmt.Errorf(
		"cannot scan %T into Optional[%s]",
		src, reflect.TypeOf(o.V).Name(),
	)
}

func (o *Optional[T]) Ptr() *T {
	if o.IsPresent && !o.IsNull {
		return &o.V
	}
	return nil
}

func Of[T any](v T) Optional[T] {
	return Optional[T]{V: v, IsPresent: true, IsNull: false}
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
		V:         n,
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

	u, err := uuid.Parse(o.V)
	if err != nil {
		return Optional[uuid.UUID]{}, fmt.Errorf("cannot parse %s as UUID: %w", o.V, err)
	}

	return Optional[uuid.UUID]{V: u, IsPresent: true}, nil
}
