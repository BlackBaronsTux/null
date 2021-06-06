package null

import (
	"bytes"
	"encoding/json"
	"testing"
)

var (
	jsonJSON = []byte(`"hello"`)
)

func TestJSONFrom(t *testing.T) {
	t.Parallel()

	i := JSONFrom([]byte(`"hello"`))
	assertJSON(t, i, "JSONFrom()")

	zero := JSONFrom(nil)
	if zero.Valid {
		t.Error("JSONFrom(nil)", "is valid, but should be invalid")
	}

	zero = JSONFrom([]byte{})
	if !zero.Valid {
		t.Error("JSONFrom([]byte{})", "is invalid, but should be valid")
	}
}

func TestJSONFromPtr(t *testing.T) {
	t.Parallel()

	n := []byte(`"hello"`)
	iptr := &n
	i := JSONFromPtr(iptr)
	assertJSON(t, i, "JSONFromPtr()")

	null := JSONFromPtr(nil)
	assertNullJSON(t, null, "JSONFromPtr(nil)")
}

type Test struct {
	Name string
	Age  int
}

func TestMarshal(t *testing.T) {
	t.Parallel()

	var i JSON

	test := &Test{Name: "hello", Age: 15}

	err := i.Marshal(test)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(i.JSON, []byte(`{"Name":"hello","Age":15}`)) {
		t.Errorf("Mismatch between received and expected, got: %s", string(i.JSON))
	}
	if i.Valid == false {
		t.Error("Expected valid true, got Valid false")
	}

	err = i.Marshal(nil)
	if err != nil {
		t.Error(err)
	}

	if i.Valid == true {
		t.Error("Expected Valid false, got Valid true")
	}
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	var i JSON

	test := &Test{}

	err := i.Unmarshal(test)
	if err != nil {
		t.Error(err)
	}

	x := &Test{Name: "hello", Age: 15}
	err = i.Marshal(x)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(i.JSON, []byte(`{"Name":"hello","Age":15}`)) {
		t.Errorf("Mismatch between received and expected, got: %s", string(i.JSON))
	}

	err = i.Unmarshal(test)
	if err != nil {
		t.Error(err)
	}

	if test.Age != 15 {
		t.Errorf("Expected 15, got %d", test.Age)
	}
	if test.Name != "hello" {
		t.Errorf("Expected name, got %s", test.Name)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	var i JSON
	err := json.Unmarshal(jsonJSON, &i)
	maybePanic(err)
	assertJSON(t, i, "[]byte json")

	var ni JSON
	err = ni.UnmarshalJSON([]byte{})
	if err != nil {
		t.Error(err)
	}
	if ni.Valid == false {
		t.Errorf("expected Valid to be true, got false")
	}
	if !bytes.Equal(ni.JSON, nil) {
		t.Errorf("Expected JSON to be nil, but was not: %#v %#v", ni.JSON, []byte(nil))
	}

	var null JSON
	err = null.UnmarshalJSON(nil)
	if err == nil {
		t.Error("passing a nil should fail")
	}
}

func TestUnmarshalJSONInStruct(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Val JSON `json:"val,omitempty"`
	}

	// In this case UnmarshalJSON is never called and it should not be
	// considered set nor valid.
	t1 := testStruct{}
	err := json.Unmarshal([]byte(`{}`), &t1)
	if err != nil {
		t.Error(err)
	}
	if t1.Val.Set {
		t.Error("should not be set, no value was given")
	}
	if t1.Val.Valid {
		t.Error("should not be valid, no value was given")
	}

	// In this case UnmarshalJSON is called with [110 117 108 108]
	// in this case the value contained in the JSON should not exist
	// and it should be set and !valid.
	//
	// This is so {"val": null} unmarshalling can turn into an sql null value.
	t2 := testStruct{}
	err = json.Unmarshal([]byte(`{"val": null}`), &t2)
	if err != nil {
		t.Error(err)
	}
	if !t2.Val.Set {
		t.Error("should be set")
	}
	if t2.Val.Valid {
		t.Error("should not be valid")
	}
}

func TestTextUnmarshalJSON(t *testing.T) {
	t.Parallel()

	var i JSON
	err := i.UnmarshalText([]byte(`"hello"`))
	maybePanic(err)
	assertJSON(t, i, "UnmarshalText() []byte")

	var blank JSON
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullJSON(t, blank, "UnmarshalText() empty []byte")
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	i := JSONFrom([]byte(`"hello"`))
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, `"hello"`, "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewJSON(nil, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalJSONText(t *testing.T) {
	t.Parallel()

	i := JSONFrom([]byte(`"hello"`))
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, `"hello"`, "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewJSON(nil, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestJSONPointer(t *testing.T) {
	t.Parallel()

	i := JSONFrom([]byte(`"hello"`))
	ptr := i.Ptr()
	if !bytes.Equal(*ptr, []byte(`"hello"`)) {
		t.Errorf("bad %s []byte: %#v ≠ %s\n", "pointer", ptr, `"hello"`)
	}

	null := NewJSON(nil, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s []byte: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestJSONIsZero(t *testing.T) {
	t.Parallel()

	i := JSONFrom([]byte(`"hello"`))
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewJSON(nil, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewJSON(nil, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestJSONSetValid(t *testing.T) {
	t.Parallel()

	change := NewJSON(nil, false)
	assertNullJSON(t, change, "SetValid()")
	change.SetValid([]byte(`"hello"`))
	assertJSON(t, change, "SetValid()")
}

func TestJSONScan(t *testing.T) {
	t.Parallel()

	var i JSON
	err := i.Scan(`"hello"`)
	maybePanic(err)
	assertJSON(t, i, "scanned []byte")

	var null JSON
	err = null.Scan(nil)
	maybePanic(err)
	assertNullJSON(t, null, "scanned null")
}

func assertJSON(t *testing.T, i JSON, from string) {
	t.Helper()
	if !bytes.Equal(i.JSON, []byte(`"hello"`)) {
		t.Errorf("bad %s []byte: %#v ≠ %#v\n", from, string(i.JSON), string([]byte(`"hello"`)))
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullJSON(t *testing.T, i JSON, from string) {
	t.Helper()
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
