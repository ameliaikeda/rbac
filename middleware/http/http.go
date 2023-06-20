// Package http sets up various middleware methods for HTTP requests, and includes options to set defaults.
//
// Methods of fetching subjects and roles include:
// - JWT token scopes (via the go-jwt library)
// - any
// - any method that returns a slice of roles; RoleProvider and UserProvider fit this.
package http
