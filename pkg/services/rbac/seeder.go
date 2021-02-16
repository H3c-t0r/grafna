package rbac

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type seeder struct {
	Service *RBACService
}

var builtInPolicies = []PolicyDTO{
	{
		Name:        "grafana:builtin:users:read:self",
		Description: "v1",
		Permissions: []Permission{
			{
				Permission: "users:read",
				Scope:      "users:self",
			},
			{
				Permission: "users.tokens:list",
				Scope:      "users:self",
			},
			{
				Permission: "users.teams:read",
				Scope:      "users:self",
			},
		},
	},
}

func (s seeder) Seed(ctx context.Context, orgID int64) error {
	// FIXME: As this will run on startup, we want to optimize running this
	existingPolicies, err := s.Service.GetPolicies(ctx, orgID)
	if err != nil {
		return err
	}
	policySet := map[string]*Policy{}
	for _, policy := range existingPolicies {
		if policy == nil {
			continue
		}
		policySet[policy.Name] = policy
	}

	for _, policy := range builtInPolicies {
		policy.OrgId = orgID

		current, exists := policySet[policy.Name]
		if exists {
			if policy.Description == current.Description {
				continue
			}
		}

		p, err := s.createOrUpdatePolicy(ctx, policy, current)
		if err != nil {
			s.Service.log.Error("failed to create/update policy", "name", policy.Name, "err", err)
			continue
		}
		if p == 0 {
			// remote version was equal or newer than current version
			continue
		}

		existingPermissions, err := s.Service.GetPolicyPermissions(ctx, p)
		if err != nil {
			s.Service.log.Info("failed to get current permissions for policy", "name", policy.Name, "err", err)
		}

		err = s.idempotentUpdatePermissions(ctx, p, policy.Permissions, existingPermissions)
		if err != nil {
			s.Service.log.Error("failed to update policy permissions", "name", policy.Name, "err", err)
		}
	}

	return nil
}

func (s seeder) createOrUpdatePolicy(ctx context.Context, policy PolicyDTO, old *Policy) (int64, error) {
	if old == nil {
		p, err := s.Service.CreatePolicy(ctx, CreatePolicyCommand{
			OrgId:       policy.OrgId,
			Name:        policy.Name,
			Description: policy.Description,
		})
		if err != nil {
			return 0, err
		}
		return p.Id, nil
	}

	// FIXME: We probably want to be able to have a description as well
	currentVersion, err := strconv.Atoi(policy.Description[1:])
	if err != nil {
		return 0, fmt.Errorf(
			"failed to read version for policy %s (\"%s\"): %w",
			policy.Name,
			policy.Description,
			err,
		)
	}

	var oldVersion int
	if strings.HasPrefix(old.Description, "v") {
		oldVersion, err = strconv.Atoi(old.Description[1:])
		if err != nil {
			return 0, fmt.Errorf(
				"failed to read previous version for policy %s (\"%s\"): %w",
				policy.Name,
				old.Description,
				err,
			)
		}
	}

	if oldVersion >= currentVersion {
		return 0, nil
	}

	_, err = s.Service.UpdatePolicy(ctx, UpdatePolicyCommand{
		UID:         old.UID,
		Name:        policy.Name,
		Description: policy.Description,
	})
	if err != nil {
		return 0, err
	}
	return old.Id, nil
}

func (s seeder) idempotentUpdatePermissions(ctx context.Context, policyID int64, new []Permission, old []Permission) error {
	added, removed := diffPermissionList(new, old)
	fmt.Println("diff", new, old)

	for _, p := range added {
		_, err := s.Service.CreatePermission(ctx, &CreatePermissionCommand{
			PolicyId:   policyID,
			Permission: p.Permission,
			Scope:      p.Scope,
		})
		if err != nil {
			return fmt.Errorf("could not create permission %s (%s): %w", p.Permission, p.Scope, err)
		}
	}

	for _, p := range removed {
		err := s.Service.DeletePermission(ctx, &DeletePermissionCommand{
			Id: p.Id,
		})
		if err != nil {
			return fmt.Errorf("could not delete permission %s (%s): %w", p.Permission, p.Scope, err)
		}
	}

	return nil
}

func diffPermissionList(new, old []Permission) (added, removed []Permission) {
	sortPermissionList(new)
	sortPermissionList(old)

	added = []Permission{}
	removed = []Permission{}

	newPos := 0
	oldPos := 0

	for {
		if newPos >= len(new) && oldPos >= len(old) {
			break
		}
		if newPos >= len(new) {
			removed = append(removed, old[oldPos:]...)
			break
		}
		if oldPos >= len(old) {
			added = append(added, new[newPos:]...)
			break
		}

		n, o := new[newPos], old[oldPos]

		switch {
		case n.Permission > o.Permission || n.Permission == o.Permission && n.Scope > o.Scope:
			oldPos++
			removed = append(removed, o)
		case n.Permission < o.Permission || n.Permission == o.Permission && n.Scope < o.Scope:
			newPos++
			added = append(added, o)
		default:
			newPos++
			oldPos++
		}
	}

	return added, removed
}

func sortPermissionList(l []Permission) {
	sort.Slice(l, func(i, j int) bool {
		if l[i].Permission == l[j].Permission {
			return l[i].Scope < l[j].Scope
		}

		return l[i].Permission < l[j].Permission
	})
}
