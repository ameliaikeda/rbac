package rbac

import (
	"context"

	"github.com/ameliaikeda/rbac/subject"
	"github.com/ameliaikeda/rbac/values"
)

// Permission holds information on something a role is allowed to do, as well as any subjects within it.
type Permission struct {
	ID          string
	Name        string
	Description string

	// Subjects are things this permission can be applied against, such as a database ID, or a special marker.
	Subjects []string
}

// Equals checks two permissions are the same.
// They can have differing subjects, but if the IDs match, they are treated the same for checks that a user has a role.
func (p Permission) Equals(cmp Permission) bool {
	return p.ID == cmp.ID
}

// ValidSubjects checks that all given subjects are valid.
// If you need a logical OR, see AnyValidSubject.
func (p Permission) ValidSubjects(ctx context.Context, subjects ...any) bool {
	for _, sub := range subjects {
		if !p.ValidSubject(ctx, sub) {
			return false
		}
	}

	return true
}

// AnyValidSubject is the inverse of ValidSubject, using logical OR.
func (p Permission) AnyValidSubject(ctx context.Context, subjects ...any) bool {
	for _, sub := range subjects {
		if p.ValidSubject(ctx, sub) {
			return true
		}
	}

	return false
}

// ValidSubject loops through all subjects on a Permission, returning true if any match.
// Special behavior is assigned to the constants subject.Wildcard and subject.Self.
//
// - subject.Wildcard will allow any subject as if it matched.
// - subject.Self will use any available auth in the context to validate against a subject (user) ID.
func (p Permission) ValidSubject(ctx context.Context, check any) bool {
	for _, expected := range p.Subjects {
		expected := expected

		// if `subject.Self` is listed on the permission, replace it with an auth user
		if expected == subject.Self {
			if str, ok := values.SubjectFromContext(ctx); ok {
				expected = str
			} else {
				continue // we do not have an auth user set up; self can never be checked and so is skipped.
			}
		}

		if subject.Matches(expected, check) {
			return true
		}
	}

	return false
}

// WithSubjects adds subjects to the current permission.
// Usage is e.g. permission.Create.WithSubjects([]string{"foo"})
func (p Permission) WithSubjects(subjects []string) Permission {
	p.Subjects = subjects

	return p
}
