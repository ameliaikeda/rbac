package generator

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"text/template"

	"golang.org/x/tools/go/packages"

	"github.com/ameliaikeda/rbac/generator/core"
)

type Output struct {
	Roles       io.ReadWriter
	Permissions io.ReadWriter
}

func Generate(_ context.Context, generator *core.Generator) (*Output, error) {
	if err := validate(generator); err != nil {
		return nil, err
	}

	roleTemplate, err := template.New("roles").Parse(generator.RoleMetadata.Template)
	if err != nil {
		return nil, err
	}

	permissionTemplate, err := template.New("permissions").Parse(generator.PermissionMetadata.Template)
	if err != nil {

		return nil, err
	}

	output := &Output{
		Roles:       &bytes.Buffer{},
		Permissions: &bytes.Buffer{},
	}

	if err := roleTemplate.Execute(output.Roles, generator); err != nil {
		return nil, err
	}

	if err := permissionTemplate.Execute(output.Permissions, generator); err != nil {
		return nil, err
	}

	return output, nil
}

type ValidationError struct {
	field string
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("validation error: field %s is required.", err.field)
}

func missing(param string) error {
	return &ValidationError{
		field: param,
	}
}

func validate(gen *core.Generator) error {
	switch {
	case len(gen.Roles) == 0:
		return missing("roles")
	case len(gen.Permissions) == 0:
		return missing("permissions")

	case gen.PermissionMetadata.Package == "":
		return missing("config.permissions.package")
	case gen.PermissionMetadata.Template == "":
		return missing("config.permissions.template")
	case gen.PermissionMetadata.Filename == "":
		return missing("config.permissions.filename")
	case gen.PermissionMetadata.ImportPath == "" && gen.RoleMetadata.Package != gen.PermissionMetadata.Package:
		return missing("config.permissions.import-path")

	case gen.RoleMetadata.Package == "":
		return missing("config.roles.package")
	case gen.RoleMetadata.Template == "":
		return missing("config.roles.template")
	case gen.RoleMetadata.Filename == "":
		return missing("config.roles.filename")
	case gen.RoleMetadata.ImportPath == "" && gen.RoleMetadata.Package != gen.PermissionMetadata.Package:
		return missing("config.roles.import-path")
	}

	return nil
}

func ImportPath(path string) (string, error) {
	pkgs, err := packages.Load(&packages.Config{}, path)
	if err != nil {
		return "", err
	}

	if len(pkgs) != 1 {
		return "", errors.New(fmt.Sprintf("rbac: unable to identify package path from path %s", path))
	}

	pkg := pkgs[0]
	if pkg == nil {
		return "", errors.New(fmt.Sprintf("rbac: import package path resolved to nil: %s", path))
	}

	if pkg.Module != nil {
		return pkg.Module.Path, nil
	}

	return pkg.PkgPath, nil
}
