package rbac

import (
	"fmt"
	"sync"
)

type internalState struct {
	sync.RWMutex
	roles   []Role
	roleMap map[string]Role
}

// state should not be manipulated outside of tests, and is not concurrency-safe to change.
var state = &internalState{
	roles:   make([]Role, 0),
	roleMap: make(map[string]Role),
}

func (s *internalState) allRoles() []Role {
	s.RLock()
	defer s.RUnlock()

	return s.roles
}

func (s *internalState) roleByID(id string) Role {
	s.RLock()
	defer s.RUnlock()

	return s.roleMap[id]
}

func (s *internalState) rolesByID(ids []string) []Role {
	s.RLock()
	defer s.RUnlock()

	roles := make([]Role, 0, len(ids))

	for _, id := range ids {
		if role, ok := s.roleMap[id]; ok {
			roles = append(roles, role)
		}
	}

	return roles
}

func (s *internalState) setRoles(roles []Role) {
	s.Lock()
	defer s.Unlock()

	s.roles = roles

	for _, role := range s.roles {
		// don't allow duplicates because of undefined behavior.
		if _, exists := s.roleMap[role.ID]; exists {
			panic(fmt.Sprintf("rbac: duplicate permission ID when setting up: %s (%s)", role.ID, role.Name))
		}

		s.roleMap[role.ID] = role
	}
}
