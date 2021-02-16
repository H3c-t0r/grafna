package rbac

import (
	"errors"
	"time"
)

var (
	errPolicyNotFound                = errors.New("policy not found")
	errTeamPolicyAlreadyAdded        = errors.New("policy is already added to this team")
	errUserPolicyAlreadyAdded        = errors.New("policy is already added to this user")
	errTeamPolicyNotFound            = errors.New("team policy not found")
	errUserPolicyNotFound            = errors.New("user policy not found")
	errTeamNotFound                  = errors.New("team not found")
	errPermissionNotFound            = errors.New("permission not found")
	errPolicyFailedGenerateUniqueUID = errors.New("failed to generate policy definition UID")
)

// Policy is the model for Policy in RBAC.
type Policy struct {
	Id          int64  `json:"id"`
	OrgId       int64  `json:"orgId"`
	UID         string `xorm:"uid" json:"uid"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

type PolicyDTO struct {
	Id          int64        `json:"id"`
	OrgId       int64        `json:"orgId"`
	UID         string       `xorm:"uid" json:"uid"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`

	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

// Policy is the model for Permission in RBAC.
type Permission struct {
	Id         int64  `json:"id"`
	PolicyId   int64  `json:"-"`
	Permission string `json:"permission"`
	Scope      string `json:"scope"`

	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

type TeamPolicy struct {
	Id       int64
	OrgId    int64
	PolicyId int64
	TeamId   int64

	Updated time.Time
	Created time.Time
}

type UserPolicy struct {
	Id       int64
	OrgId    int64
	PolicyId int64
	UserId   int64

	Updated time.Time
	Created time.Time
}

type GetTeamPoliciesQuery struct {
	OrgId  int64 `json:"-"`
	TeamId int64
}

type GetUserPoliciesQuery struct {
	OrgId  int64 `json:"-"`
	UserId int64
}

type GetUserPermissionsQuery struct {
	OrgId  int64 `json:"-"`
	UserId int64
}

type CreatePermissionCommand struct {
	PolicyId   int64
	Permission string
	Scope      string
}

type UpdatePermissionCommand struct {
	Id         int64
	Permission string
	Scope      string
}

type DeletePermissionCommand struct {
	Id int64
}

type CreatePolicyCommand struct {
	OrgId       int64  `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdatePolicyCommand struct {
	OrgId       int64  `json:"-"`
	UID         string `json:"uid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeletePolicyCommand struct {
	Id    int64
	UID   string `json:"uid"`
	OrgId int64
}

type AddTeamPolicyCommand struct {
	OrgId    int64
	PolicyId int64
	TeamId   int64
}

type RemoveTeamPolicyCommand struct {
	OrgId    int64
	PolicyId int64
	TeamId   int64
}

type AddUserPolicyCommand struct {
	OrgId    int64
	PolicyId int64
	UserId   int64
}

type RemoveUserPolicyCommand struct {
	OrgId    int64
	PolicyId int64
	UserId   int64
}

type EvaluationResult struct {
	HasAccess bool
	Meta      interface{}
}
