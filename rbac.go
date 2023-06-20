// Package rbac is a generic library based on code generation to create a role-based access control system.
//
// It is designed in a way that makes it simple, and includes middleware to gate requests based on context.
package rbac

// SetDefaultRoles is something that should ideally be called from an init function.
// While it is concurrency-safe for read and write access, it's not advisable to change state between requests.
func SetDefaultRoles(roles []Role) {
	if state == nil {
		panic("rbac: can't set internal permissions; state is nil")
	}

	state.setRoles(roles)
}
