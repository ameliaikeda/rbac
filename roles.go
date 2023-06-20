package rbac

import (
	"context"
)

// Role is a base unit that can be assigned a group or a user, and contains a set of permissions.
//
// - A user MAY have many roleLookup
// - A user MAY have many groups
// - A group MAY have many roleLookup
//
// A user's roleLookup are usually determined as a flattened set of all group roleLookup, plus all direct user roleLookup.
// The method of obtaining a user's roleLookup is left up to the implementer or middleware.
//
// Roles are considered equal if the ID value matches, but a role with an empty ID will never match another role.
type Role struct {
	ID          string
	Name        string
	Description string
	Permissions []Permission

	// CustomMappings holds any extra data at generation time used to map this role to another system.
	CustomMappings map[string]string
}

// Can checks if a role has a specific permission. If a subject is passed, they are verified via logical AND.
func (r Role) Can(ctx context.Context, perm Permission, subjects ...any) bool {
	if !r.Has(perm) {
		return false
	}

	return perm.ValidSubjects(ctx, subjects...)
}

// Has checks if a role has a specific permission, regardless of subjects.
func (r Role) Has(perm Permission) bool {
	for _, p := range r.Permissions {
		if perm.Equals(p) {
			return true
		}
	}

	return false
}
