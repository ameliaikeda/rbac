# rbac - A role-based access control library for Go.

To install and run `rbac` into `$GOBIN/rbac`, run `go install`:

```
go install github.com/ameliaikeda/rbac/cmd/rbac@v0.1.0
```

**Note**: This library is still in active development and can change at any time.

# YAML Format

```yaml
---
# see TemplateMetadata on the wiki
config:
  roles:
    package: "role"
    filename: "role_gen.go"
    path: "rbac" # this means $(pwd)/rbac/role/role_gen.go is created (default)
    
  permissions:
    # same as above

permissions:
  create_item:
    name: "Create Item"
  edit_item:
    name: "Edit Item"

roles:
  admin:
    id: "00000000-0000-0000-0000-000000000000"
    name: "Admin"
    description: "Admin Users"
    key: "admin" # this overrides the key for this object above
    permissions:
      - create_item
      - edit_item
   
   user:
    id: "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF"
    name: "User"
    description: "Standard Users"
    key: "user"
    permissions:
      - create_item
      - edit_item:self # :self and :any are special cases listed further in.
```

See the wiki (TODO) for more info on the full YAML format, including `go-name` directives.

# Usage

```
rbac [-config rbac.yaml]
```

This will create `rbac/role/role_gen.go` and `rbac/permission/permission_gen.go` by default, configurable via the YAML file above.

It can also be put into a `go:generate` comment:

```
//go:generate go run github.com/ameliaikeda/rbac/cmd/rbac@v0.1.0 -config config.yaml
```

# Developer interface

Basically everything can be set up via two interfaces:

1. A middleware that embeds a user object and its roles in the context of a request via `values.Embed`
2. A `rbac.Can(ctx, permissions.Edit, subject) bool` method that verifies if a user can do a gated action.

You can optionally (instead of using `subject` as a string for comparison, like an ID etc), implement this interface on `subject`:

```go
type rbacSubject interface {
  RBACSubjectID() string
}
```

If a struct does not implement `String` or `RBACSubjectID()`, and cannot be coerced to a string, it will panic.
Likewise, `slice`, `chan`, and other non-basic Go types will fail if tyhey do not implement the above interfaces.
All primitive types but `uintptr` and `complex*` are coercable and will work.

In practice, this means you can simply implement `RBACSubjectID` on your User models.
