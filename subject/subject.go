// Package subject declares methods and constants for checking a permission's subjects.
//
// Most commonly used are Wildcard and Self; the former means the permission applies for any subject.
// The latter, Self, means the resolved subject ID (usually a user's database ID, or similar) is checked.
package subject

import (
	"fmt"
)

const (
	// Wildcard should be used when a role should have access to any subject via a permission.
	Wildcard = "*"

	// Self should be used when a permission is only granted to a resource that matches the user's ID.
	Self = "rbac.self"
)

// Matches checks that two values for subjects match.
// If `expected` is Wildcard, this function always returns true.
//
// The constant `Self` is always replaced with the expected SubjectID of the user in calling functions.
// Calling Matches(Self, any) is not valid on its own.
//
// Match values are always converted to a string, in priority order:
// - Strings are passed as-is.
// - Integers and floats are converted to strings.
// - Anything with an `RBACSubjectID() string` method uses the result of that.
// - If the given type has a fmt.Stringer method, we use the result.
// - If nothing else is valid, we raise a panic.
func Matches(expected, actual any) bool {
	if expected == Wildcard {
		return true
	}

	return convert(expected) == convert(actual)
}

// convert takes a subject and converts it to a string via most means available.
func convert(sub any) string {
	type rbacSubject interface {
		RBACSubjectID() string
	}

	if s, ok := sub.(rbacSubject); ok {
		return s.RBACSubjectID()
	}

	switch s := sub.(type) {
	case string:
		return s

	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprintf("%v", s)

	case rbacSubject:
		return s.RBACSubjectID()

	case fmt.Stringer:
		return s.String()
	}

	// if we're given an incompatible type (struct, slice/array, complex, uintptr, chan), bail out.
	panic(fmt.Sprintf("rbac: subject doesn't implement RBACSubjectID, or can't coerce to string: %T given", sub))
}
