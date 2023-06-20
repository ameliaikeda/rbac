package values

import (
	"context"

	"github.com/go-logr/logr"
)

// contextKeyType is an unexported type for context.WithValue.
type contextKeyType struct{}

// contextKey is used as a key for pulling values out of a context.
var contextKey contextKeyType

// User is an interface used to provide a subject ID and a slice of roles.
type User interface {
	RBACSubjectID() string
	RBACRoles() []string
}

// FromContext grabs a user from a context, if it has been set as a value.
func FromContext(ctx context.Context) User {
	v := ctx.Value(contextKey)

	if u, ok := v.(User); ok {
		return u
	}

	return nil
}

// Embed embeds a user into the current request context.
func Embed(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, contextKey, user)
}

func SubjectFromContext(ctx context.Context) (string, bool) {
	if u := FromContext(ctx); u != nil {
		return u.RBACSubjectID(), true
	}

	return "", false
}

func RolesFromContext(ctx context.Context) ([]string, bool) {
	if u := FromContext(ctx); u != nil {
		logr.FromContextOrDiscard(ctx).WithName("rbac").Info("looking up roles for user",
			"rbac.subject.id", u.RBACSubjectID())

		roles := u.RBACRoles()

		if len(roles) == 0 {
			logr.FromContextOrDiscard(ctx).Info("rbac: user roles empty",
				"rbac.subject.id", u.RBACSubjectID())
		}

		return roles, true
	}

	return nil, false
}
