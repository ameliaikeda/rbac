package rbac

import (
	"context"

	"github.com/go-logr/logr"

	"github.com/ameliaikeda/rbac/values"
)

func log(ctx context.Context, msg string, kv ...any) {
	kv = append(kv, "library", "rbac")
	logr.FromContextOrDiscard(ctx).WithName("rbac").Info(msg, kv...)
}

// Can uses the current context values to determine if an action can be taken.
//
// Usage: if rbac.Can(ctx, permissions.SpecificationCreate) {}
func Can(ctx context.Context, perm Permission, subjects ...any) bool {
	result := false

	for _, role := range Roles(ctx) {
		if role.Can(ctx, perm, subjects...) {
			result = true
		}
	}

	log(ctx, "checking permissions",
		"rbac.permission.id", perm.ID,
		"rbac.result", result)

	return result
}

// Roles returns a list of roleLookup in the current context.
// The roleLookup must have been set up globally.
func Roles(ctx context.Context) []Role {
	if user := User(ctx); user != nil {
		roles := state.rolesByID(user.RBACRoles())

		if len(roles) == 0 {
			log(ctx, "no roles present on subject", "rbac.subject.id", user.RBACSubjectID())
		}

		return roles
	}

	return nil
}

func User(ctx context.Context) values.User {
	log(ctx, "looking up user in context")
	user := values.FromContext(ctx)
	if user == nil {
		log(ctx, "user not found in context")
	}

	return user
}
