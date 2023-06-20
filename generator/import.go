package generator

import (
	"context"
	"path/filepath"

	"github.com/ameliaikeda/rbac/generator/core"
)

func BasePath(path string) core.OptionFunc {
	return func(_ context.Context, gen *core.Generator) error {
		gen.BaseDirectory = path
		return nil
	}
}

func ResolvePaths(_ context.Context, gen *core.Generator) error {
	roleAbsolute, err := filepath.Abs(filepath.Join(gen.BaseDirectory, gen.RoleMetadata.Path, gen.RoleMetadata.Package))
	if err != nil {
		return err
	}

	rolePkg, err := ImportPath(roleAbsolute)
	if err != nil {
		return err
	}

	permAbsolute, err := filepath.Abs(
		filepath.Join(gen.BaseDirectory, gen.PermissionMetadata.Path, gen.PermissionMetadata.Package),
	)
	if err != nil {
		return err
	}

	permPkg, err := ImportPath(permAbsolute)
	if err != nil {
		return err
	}

	gen.RoleMetadata.Path = roleAbsolute
	gen.RoleMetadata.ImportPath = rolePkg

	gen.PermissionMetadata.ImportPath = permPkg
	gen.PermissionMetadata.Path = permAbsolute

	return nil
}
