package generator

import (
	"context"
	_ "embed"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"

	"github.com/ameliaikeda/rbac/generator/core"
)

func Run(ctx context.Context, opts ...core.OptionFunc) error {
	options := make([]core.OptionFunc, 0, len(opts)+1)
	options = append(options, Defaults)
	options = append(options, opts...)

	gen := New()

	for _, option := range options {
		if err := option(ctx, gen); err != nil {
			return err
		}
	}

	output, err := Generate(ctx, gen)
	if err != nil {
		return err
	}

	if err := writeOutput(output.Roles, gen.RoleMetadata.Path, gen.RoleMetadata.Filename); err != nil {
		return err
	}

	if err := writeOutput(output.Permissions, gen.PermissionMetadata.Path, gen.PermissionMetadata.Filename); err != nil {
		return err
	}

	return nil
}

func writeOutput(output io.Reader, folder, filename string) error {
	b, err := io.ReadAll(output)
	if err != nil {
		return err
	}

	formatted, err := format.Source(b)
	if err != nil {
		return err
	}

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	full := filepath.Join(folder, filename)

	f, err := os.Create(full)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(formatted)
	if err != nil {
		return err
	}

	fmt.Printf("generated: %s\n", full)

	return nil
}

func readFormatted(r io.Reader) ([]byte, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return format.Source(b)
}

func New() *core.Generator {
	return &core.Generator{
		PermissionMetadata: &core.TemplateMetadata{},
		RoleMetadata:       &core.TemplateMetadata{},
		Roles:              make([]core.Role, 0),
		Permissions:        make([]core.Permission, 0),
	}
}

var (
	//go:embed data/role.gotmpl
	roleData string

	//go:embed data/permission.gotmpl
	permissionData string
)

func Defaults(_ context.Context, gen *core.Generator) error {
	gen.ConfigFilename = "rbac.yaml"

	gen.PermissionMetadata.Package = "permission"
	gen.PermissionMetadata.Filename = "permission_gen.go"
	gen.PermissionMetadata.Template = permissionData

	gen.RoleMetadata.Package = "role"
	gen.RoleMetadata.Filename = "role_gen.go"
	gen.RoleMetadata.Template = roleData

	return nil
}
