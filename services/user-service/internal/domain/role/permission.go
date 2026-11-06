package role

import (
	"fmt"
	"slices"
)

type Permission string

const (
	PermissionOrganizationUpdate Permission = "organization.update"
	PermissionMembersInvite      Permission = "members.invite"
	PermissionMembersRemove      Permission = "members.remove"
	PermissionMembersAssignRole  Permission = "members.assign_role"
	PermissionRolesManage        Permission = "roles.manage"
)

var allPermissions = []Permission{
	PermissionOrganizationUpdate,
	PermissionMembersInvite,
	PermissionMembersRemove,
	PermissionMembersAssignRole,
	PermissionRolesManage,
}

func AllPermissions() []Permission {
	return slices.Clone(allPermissions)
}

func ParsePermission(s string) (Permission, error) {
	p := Permission(s)
	if !slices.Contains(allPermissions, p) {
		return "", fmt.Errorf("%w: %q", ErrUnknownPermission, s)
	}
	return p, nil
}

func (p Permission) String() string {
	return string(p)
}
