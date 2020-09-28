package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/razor-1/null/v9/convert"
)

// Int16 is an nullable int16.
type Int16 struct {
	Int16 int16
	Valid bool
	set   bool
}

// NewInt16 creates a new Int16
func NewInt16(i int16, valid, set bool) Int16 {
	return Int16{
		Int16: i,
		Valid: valid,
		set:   set,
	}
}

// Int16From creates a new Int16 that will always be valid.
func Int16From(i int16) Int16 {
	return NewInt16(i, true, true)
}

// Int16FromPtr creates a new Int16 that be null if i is nil.
func Int16FromPtr(i *int16) Int16 {
	if i == nil {
		return NewInt16(0, false, true)
	}
	return NewInt16(*i, true, true)
}

func (i Int16) IsSet() bool {
	return i.set
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int16) UnmarshalJSON(data []byte) error {
	i.set = true
	if bytes.Equal(data, NullBytes) {
		i.Int16, i.Valid = 0, false
		return nil
	}

	var x int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x > math.MaxInt16 {
		return fmt.Errorf("json: %d overflows max int16 value", x)
	}

	i.Int16, i.Valid = int16(x), true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int16) UnmarshalText(text []byte) error {
	i.set = true
	if text == nil || len(text) == 0 {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 16)
	i.Valid = err == nil
	if i.Valid {
		i.Int16 = int16(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (i Int16) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int16), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int16) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int16), 10)), nil
}

// SetValid changes this Int16's value and also sets it to be non-null.
func (i *Int16) SetValid(n int16) {
	i.Int16 = n
	i.Valid = true
	i.set = true
}

// Ptr returns a pointer to this Int16's value, or a nil pointer if this Int16 is null.
func (i Int16) Ptr() *int16 {
	if !i.Valid {
		return nil
	}
	return &i.Int16
}

// IsZero returns true for invalid Int16's, for future omitempty support (Go 1.4?)
func (i Int16) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (i *Int16) Scan(value interface{}) error {
	if value == nil {
		i.Int16, i.Valid, i.set = 0, false, false
		return nil
	}
	i.Valid, i.set = true, true
	return convert.ConvertAssign(&i.Int16, value)
}

// Value implements the driver Valuer interface.
func (i Int16) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Int16), nil
}

// Randomize for sqlboiler
func (i *Int16) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		i.Int16 = 0
		i.Valid = false
	} else {
		i.Int16 = int16(nextInt() % math.MaxInt16)
		i.Valid = true
	}
}
