package null

// Value specifies methods that allow introspection into the
// state of a value.
type Value interface {
	// IsValid returns true iff the value is set and non-null
	IsValid() bool
	// IsSet returns true iff the value is set (null inclusive)
	IsSet() bool
}
