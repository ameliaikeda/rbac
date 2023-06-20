package yaml

import (
	"context"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"

	"github.com/ameliaikeda/rbac"
	"github.com/ameliaikeda/rbac/generator/core"
	"github.com/ameliaikeda/rbac/generator/mapping"
	"github.com/ameliaikeda/rbac/subject"
)

// Config is the primary struct used to decode the configuration JSON file.
type Config struct {
	Metadata    Metadata              `yaml:"config"`
	Roles       map[string]Role       `yaml:"roles"`
	Permissions map[string]Permission `yaml:"permissions"`
}

// Metadata is used to alter role generation.
type Metadata struct {
	Permissions *core.TemplateMetadata `yaml:"permissions"`
	Roles       *core.TemplateMetadata `yaml:"roles"`
}

// Role holds all configuration info for a role.
type Role struct {
	ID          string   `yaml:"id"`
	Key         string   `yaml:"-"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Permissions []string `yaml:"permissions"`

	// extra params

	// GoName overrides the automatically generated Go Name for the roles.
	// The default is <path>/roles.PascalCase(ID). If ID and Key differ, we defer to Key.
	// This allows using e.g. admin (roles.Admin) as the key, but keeping the ID as my-orgname-admin.
	GoName string `yaml:"go-name"`

	// ActiveDirectory is a mapping to record against the role for mapping active directory groups in single sign-on.
	ActiveDirectory string `yaml:"ad-mapping"`
}

// Permission holds the configuration info for all permissions.
// NB: This may be removed later and worked out automatically from Role.
type Permission struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	GoName      string `yaml:"go-name"`
}

func OptionsFromYAML(filename string) core.OptionFunc {
	if filename == "" {
		filename = "rbac.yaml"
	}

	return func(ctx context.Context, gen *core.Generator) (err error) {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}

		defer func() {
			// don't overwrite a YAML error if we can't close the file.
			// in most systems this is because the file was closed twice, or it no longer exists.
			fileErr := f.Close()
			if err == nil {
				err = fileErr
			}
		}()

		var config Config

		decoder := yaml.NewDecoder(f)
		if err := decoder.Decode(&config); err != nil {
			return err
		}

		if config.Metadata.Roles != nil {
			mergeMetadata(gen.RoleMetadata, config.Metadata.Roles)
		}

		if config.Metadata.Permissions != nil {
			mergeMetadata(gen.PermissionMetadata, config.Metadata.Permissions)
		}

		permsMap := make(map[string]core.Permission)
		perms := make([]core.Permission, 0, len(config.Permissions))

		for id, perm := range config.Permissions {
			p := marshalPermission(id, perm)

			permsMap[id] = p
			perms = append(perms, p)
		}

		gen.Permissions = perms

		roles := make([]core.Role, 0, len(config.Roles))
		for id, role := range config.Roles {
			roles = append(roles, marshalRole(id, role, permsMap))
		}

		gen.Roles = roles

		return nil
	}
}

func mergeMetadata(meta *core.TemplateMetadata, add *core.TemplateMetadata) {
	if add.Path != "" {
		meta.Path = add.Path
	}

	if add.Package != "" {
		meta.Package = add.Package
	}

	if add.Filename != "" {
		meta.Package = add.Package
	}

	if add.Template != "" {
		meta.Template = add.Template
	}

	if len(add.Tags) > 0 {
		meta.Tags = add.Tags
	}
}

func marshalRole(key string, role Role, permissions map[string]core.Permission) core.Role {
	// if key is blank, use the ID.
	if role.Key == "" {
		role.Key = key
	}

	if role.ID == "" {
		role.ID = key
	}

	if role.GoName == "" {
		role.GoName = strcase.ToCamel(key)
	}

	rbacRole := rbac.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CustomMappings: map[string]string{
			mapping.ActiveDirectoryGroupName: role.ActiveDirectory,
		},
	}

	coreRole := core.Role{
		Role:   rbacRole,
		GoName: role.GoName,
	}

	return marshalRolePermissions(role, coreRole, permissions)
}

func marshalRolePermissions(template Role, role core.Role, permissions map[string]core.Permission) core.Role {
	perms := make([]core.Permission, 0, len(template.Permissions))
	for _, name := range template.Permissions {
		name, subjects := parseRolePermission(name)

		if p, ok := permissions[name]; ok {
			p.Subjects = subjects

			perms = append(perms, p)
		}
	}

	role.Permissions = perms

	return role
}

func parseRolePermission(name string) (string, []string) {
	if strings.Contains(name, ":") {
		// we are dealing with a subject permission, e.g. - edit:self, create:*
		contents := strings.SplitN(name, ":", 2)

		subjects := make([]string, 0)

		for _, sub := range strings.Split(contents[1], ",") {
			switch sub {
			case subject.Wildcard, "any":
				subjects = append(subjects, subject.Wildcard)
			case subject.Self, "self":
				subjects = append(subjects, subject.Self)
			default:
				subjects = append(subjects, sub)
			}
		}

		return contents[0], subjects
	}

	return name, nil
}

func marshalPermission(key string, permission Permission) core.Permission {
	if permission.ID == "" {
		permission.ID = key
	}

	if permission.GoName == "" {
		permission.GoName = strcase.ToCamel(key)
	}

	return core.Permission{
		Permission: rbac.Permission{
			ID:          permission.ID,
			Name:        permission.Name,
			Description: permission.Description,
		},
		GoName: permission.GoName,
	}
}
