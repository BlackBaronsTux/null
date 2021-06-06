package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/volatiletech/null/v8/convert"
	"github.com/volatiletech/randomize"
)

// JSON is a nullable []byte that contains JSON.
//
// You might want to use this in the case where you have say a nullable
// JSON column in postgres for instance, where there is one layer of null for
// the postgres column, and then you also have the opportunity to have null
// as a value contained in the json. When unmarshalling json however you
// cannot set 'null' as a value.
type JSON struct {
	JSON  []byte
	Valid bool
	Set   bool
}

// NewJSON creates a new JSON
func NewJSON(b []byte, valid bool) JSON {
	return JSON{
		JSON:  b,
		Valid: valid,
		Set:   true,
	}
}

// JSONFrom creates a new JSON that will be invalid if nil.
func JSONFrom(b []byte) JSON {
	return NewJSON(b, b != nil)
}

// JSONFromPtr creates a new JSON that will be invalid if nil.
func JSONFromPtr(b *[]byte) JSON {
	if b == nil {
		return NewJSON(nil, false)
	}
	n := NewJSON(*b, true)
	return n
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (j JSON) IsValid() bool {
	return j.Set && j.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (j JSON) IsSet() bool {
	return j.Set
}

// Unmarshal will unmarshal your JSON stored in
// your JSON object and store the result in the
// value pointed to by dest.
func (j JSON) Unmarshal(dest interface{}) error {
	if dest == nil {
		return errors.New("destination is nil, not a valid pointer to an object")
	}

	// Call our implementation of
	// JSON MarshalJSON through json.Marshal
	// to get the value of the JSON object
	res, err := json.Marshal(j)
	if err != nil {
		return err
	}

	return json.Unmarshal(res, dest)
}

// UnmarshalJSON implements json.Unmarshaler.
//
// Example if you have a struct with a null.JSON called v:
//
// 		{}          -> does not call unmarshaljson: !set & !valid
// 		{"v": null} -> calls unmarshaljson, set & !valid
//      {"v": {}}   -> calls unmarshaljson, set & valid (json value is '{}')
//
// That's to say if 'null' is passed in at the json level we do not capture that
// value - instead we set the value-level null flag so that an sql value will
// turn out null.
func (j *JSON) UnmarshalJSON(data []byte) error {
	if data == nil {
		return fmt.Errorf("null: cannot unmarshal nil into Go value of type null.JSON")
	}

	j.Set = true

	if bytes.Equal(data, NullBytes) {
		j.JSON = nil
		j.Valid = false
		return nil
	}

	j.Valid = true
	j.JSON = make([]byte, len(data))
	copy(j.JSON, data)

	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (j *JSON) UnmarshalText(text []byte) error {
	j.Set = true
	if len(text) == 0 {
		j.JSON = nil
		j.Valid = false
	} else {
		j.JSON = append(j.JSON[0:0], text...)
		j.Valid = true
	}

	return nil
}

// Marshal will marshal the passed in object,
// and store it in the JSON member on the JSON object.
func (j *JSON) Marshal(obj interface{}) error {
	res, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Call our implementation of
	// JSON UnmarshalJSON through json.Unmarshal
	// to Set the result to the JSON object
	return json.Unmarshal(res, j)
}

// MarshalJSON implements json.Marshaler.
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j.JSON) == 0 || j.JSON == nil {
		return NullBytes, nil
	}
	return j.JSON, nil
}

// MarshalText implements encoding.TextMarshaler.
func (j JSON) MarshalText() ([]byte, error) {
	if !j.Valid {
		return nil, nil
	}
	return j.JSON, nil
}

// SetValid changes this JSON's value and also sets it to be non-null.
func (j *JSON) SetValid(n []byte) {
	j.JSON = n
	j.Valid = true
	j.Set = true
}

// Ptr returns a pointer to this JSON's value, or a nil pointer if this JSON is null.
func (j JSON) Ptr() *[]byte {
	if !j.Valid {
		return nil
	}
	return &j.JSON
}

// IsZero returns true for null or zero JSON's, for future omitempty support (Go 1.4?)
func (j JSON) IsZero() bool {
	return !j.Valid
}

// Scan implements the Scanner interface.
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		j.JSON, j.Valid, j.Set = nil, false, false
		return nil
	}
	j.Valid, j.Set = true, true
	return convert.ConvertAssign(&j.JSON, value)
}

// Value implements the driver Valuer interface.
func (j JSON) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return j.JSON, nil
}

// Randomize for sqlboiler
func (j *JSON) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	j.JSON = []byte(`"` + randomize.Str(nextInt, 1) + `"`)
	j.Valid = true
}
