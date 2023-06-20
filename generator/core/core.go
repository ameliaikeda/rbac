package core

import (
	"context"

	"github.com/ameliaikeda/rbac"
)

type OptionFunc func(context.Context, *Generator) error

// Generator takes a set of options and creates the roles and permissions for RBAC usage.
// They all have defaults, but everything can be overridden with an OptionFunc.
type Generator struct {
	PermissionMetadata *TemplateMetadata
	RoleMetadata       *TemplateMetadata
	Roles              []Role
	Permissions        []Permission
	ConfigFilename     string
	BaseDirectory      string
}

// TemplateMetadata is used to generate a file from a template.
type TemplateMetadata struct {
	// Package name to use.
	// Default: role or permission, singular.
	Package string `json:"package" yaml:"package"`

	// Filename is the filename to generate.
	// Default: (package)_gen.go
	Filename string `json:"filename" yaml:"filename"`

	// Path is the folder in which a file should be generated.
	// Default: ./rbac/(package)
	Path string `json:"path" yaml:"path"`

	// ImportPath is the path used to import this package from another, e.g. github.com/example.
	ImportPath string `json:"-" yaml:"-"`

	// Template is a custom go text/template to use for generation.
	Template string `json:"template" yaml:"template"`

	// Tags indicated which struct tags, if any, to add to generated structs.
	// Currently handled: json, yaml, db.
	// To add anything extra, override Template.
	Tags []string `json:"tags" yaml:"tags"`
}

type Role struct {
	rbac.Role
	Permissions []Permission

	// we also add additional fields specific to generation, for ease of use.
	GoTags string
	GoName string
}

type Permission struct {
	rbac.Permission

	// we also add additional fields specific to generation, for ease of use.
	GoTags string
	GoName string
}
