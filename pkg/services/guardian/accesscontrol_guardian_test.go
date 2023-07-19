package guardian

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	accesscontrolmock "github.com/grafana/grafana/pkg/services/accesscontrol/mock"
	"github.com/grafana/grafana/pkg/services/accesscontrol/ossaccesscontrol"
	"github.com/grafana/grafana/pkg/services/dashboards"
	dashdb "github.com/grafana/grafana/pkg/services/dashboards/database"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/folder/foldertest"
	"github.com/grafana/grafana/pkg/services/licensing/licensingtest"
	"github.com/grafana/grafana/pkg/services/quota/quotatest"
	"github.com/grafana/grafana/pkg/services/supportbundles/supportbundlestest"
	"github.com/grafana/grafana/pkg/services/tag/tagimpl"
	"github.com/grafana/grafana/pkg/services/team/teamimpl"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/services/user/userimpl"
	"github.com/grafana/grafana/pkg/setting"
)

type accessControlGuardianTestCase struct {
	desc           string
	dashUID        string
	dashIsFolder   bool
	permissions    []accesscontrol.Permission
	viewersCanEdit bool
	expected       bool
}

func TestAccessControlDashboardGuardian_CanSave(t *testing.T) {
	tests := []accessControlGuardianTestCase{
		{
			desc:    "should be able to save with dashboard wildcard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "dashboards:*",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to save with folder wildcard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "folders:*",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to save with dashboard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "dashboards:uid:1",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to save with folder general scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "folders:uid:general",
				},
			},
			expected: true,
		},
		/*
			{
				desc:    "should be able to save with folder scope",
				dashUID: "1",
				permissions: []accesscontrol.Permission{
					{
						Action: dashboards.ActionDashboardsWrite,
						Scope:  "folders:uid:1",
					},
				},
				expected: true,
			},
		*/
		{
			desc:    "should not be able to save with incorrect dashboard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "dashboards:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:    "should not be able to save with incorrect folder scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "folders:uid:100",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to save with folder write and dashboard wildcard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "dashboards:*",
				},
			},
			expected: false,
		},
		{
			desc:         "should be able to save with folder write and folder wildcard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "folders:*",
				},
			},
			expected: true,
		},
		{
			desc:         "should not be able to save with folder write and dashboard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "dashboards:uid:1",
				},
			},
			expected: false,
		},
		{
			desc:         "should be able to save with folder write and folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "folders:uid:1",
				},
			},
			expected: true,
		},
		{
			desc:         "should not be able to save with folder write and incorrect dashboard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "dashboards:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to save with folder write and incorrect folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "folders:uid:100",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			guardian := setupAccessControlGuardianTest(t, tt.dashUID, tt.permissions, nil, nil, nil, tt.dashIsFolder)
			can, err := guardian.CanSave()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, can)
		})
	}
}

func TestAccessControlDashboardGuardian_CanEdit(t *testing.T) {
	tests := []accessControlGuardianTestCase{
		{
			desc:    "should be able to edit with dashboard wildcard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "dashboards:*",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to edit with folder wildcard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "folders:*",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to edit with dashboard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "dashboards:uid:1",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to edit with folder scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "folders:uid:general",
				},
			},
			expected: true,
		},
		{
			desc:    "should not be able to edit with incorrect dashboard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "dashboards:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:    "should not be able to edit with incorrect folder scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsWrite,
					Scope:  "folders:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:    "should be able to edit with read action when viewer_can_edit is true",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsRead,
					Scope:  "dashboards:uid:1",
				},
			},
			viewersCanEdit: true,
			expected:       true,
		},
		{
			desc:         "should not be able to edit with folder write and dashboard wildcard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "dashboards:*",
				},
			},
			expected: false,
		},
		{
			desc:         "should be able to edit with folder write and folder wildcard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "folders:*",
				},
			},
			expected: true,
		},
		{
			desc:         "should not be able to edit with folder write and dashboard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "dashboards:uid:1",
				},
			},
			expected: false,
		},
		{
			desc:         "should be able to edit with folder write and folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "folders:uid:1",
				},
			},
			expected: true,
		},
		{
			desc:         "should not be able to edit with folder write and incorrect folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersWrite,
					Scope:  "folders:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:         "should be able to edit with folder read action when viewer_can_edit is true",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersRead,
					Scope:  "folders:uid:1",
				},
			},
			viewersCanEdit: true,
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cfg := setting.NewCfg()
			cfg.ViewersCanEdit = tt.viewersCanEdit
			guardian := setupAccessControlGuardianTest(t, tt.dashUID, tt.permissions, cfg, nil, nil, tt.dashIsFolder)

			can, err := guardian.CanEdit()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, can)
		})
	}
}

func TestAccessControlDashboardGuardian_CanView(t *testing.T) {
	tests := []accessControlGuardianTestCase{
		{
			desc:    "should be able to view with dashboard wildcard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsRead,
					Scope:  "dashboards:*",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to view with folder wildcard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsRead,
					Scope:  "folders:*",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to view with dashboard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsRead,
					Scope:  "dashboards:uid:1",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to view with folder general scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsRead,
					Scope:  "folders:uid:general",
				},
			},
			expected: true,
		},
		/*
			{
				desc:    "should be able to view with folder scope",
				dashUID: "1",
				permissions: []accesscontrol.Permission{
					{
						Action: dashboards.ActionDashboardsRead,
						Scope:  "folders:uid:1",
					},
				},
				expected: true,
			},
		*/
		{
			desc:    "should not be able to view with incorrect dashboard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsRead,
					Scope:  "dashboards:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:    "should not be able to view with incorrect folder scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsRead,
					Scope:  "folders:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to view with folders read and dashboard wildcard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersRead,
					Scope:  "dashboards:*",
				},
			},
			expected: false,
		},
		{
			desc:         "should be able to view with folders read and folder wildcard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersRead,
					Scope:  "folders:*",
				},
			},
			expected: true,
		},
		{
			desc:         "should not be able to view with folders read and dashboard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersRead,
					Scope:  "dashboards:uid:1",
				},
			},
			expected: false,
		},
		{
			desc:         "should be able to view with folders read and folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersRead,
					Scope:  "folders:uid:1",
				},
			},
			expected: true,
		},
		{
			desc:         "should not be able to view with folders read incorrect dashboard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersRead,
					Scope:  "dashboards:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to view with folders read and incorrect folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersRead,
					Scope:  "folders:uid:10",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			guardian := setupAccessControlGuardianTest(t, tt.dashUID, tt.permissions, nil, nil, nil, tt.dashIsFolder)

			can, err := guardian.CanView()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, can)
		})
	}
}
func TestAccessControlDashboardGuardian_CanAdmin(t *testing.T) {
	tests := []accessControlGuardianTestCase{
		{
			desc:    "should be able to admin with dashboard wildcard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsPermissionsRead,
					Scope:  "dashboards:*",
				},
				{
					Action: dashboards.ActionDashboardsPermissionsWrite,
					Scope:  "dashboards:*",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to admin with folder wildcard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsPermissionsRead,
					Scope:  "folders:*",
				},
				{
					Action: dashboards.ActionDashboardsPermissionsWrite,
					Scope:  "folders:*",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to admin with dashboard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsPermissionsRead,
					Scope:  "dashboards:uid:1",
				},
				{
					Action: dashboards.ActionDashboardsPermissionsWrite,
					Scope:  "dashboards:uid:1",
				},
			},
			expected: true,
		},
		/*
				{
					desc:    "should be able to admin with folder general scope",
					dashUID: "1",
					permissions: []accesscontrol.Permission{
						{
							Action: dashboards.ActionDashboardsPermissionsRead,
							Scope:  "folders:uid:general",
						},
						{
							Action: dashboards.ActionDashboardsPermissionsWrite,
							Scope:  "folders:uid:general",
						},
					},
					expected: true,
				},
			{
				desc:    "should be able to admin with folder scope",
				dashUID: "1",
				permissions: []accesscontrol.Permission{
					{
						Action: dashboards.ActionDashboardsPermissionsRead,
						Scope:  "folders:uid:1",
					},
					{
						Action: dashboards.ActionDashboardsPermissionsWrite,
						Scope:  "folders:uid:1",
					},
				},
				expected: true,
			},
		*/
		{
			desc:    "should not be able to admin with incorrect dashboard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsPermissionsRead,
					Scope:  "dashboards:uid:10",
				},
				{
					Action: dashboards.ActionDashboardsPermissionsWrite,
					Scope:  "dashboards:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:    "should not be able to admin with incorrect folder scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsPermissionsRead,
					Scope:  "folders:uid:10",
				},
				{
					Action: dashboards.ActionDashboardsPermissionsWrite,
					Scope:  "folders:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to admin with folder read and write and dashboard wildcard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersPermissionsRead,
					Scope:  "dashboards:*",
				},
				{
					Action: dashboards.ActionFoldersPermissionsWrite,
					Scope:  "dashboards:*",
				},
			},
			expected: false,
		},
		{
			desc:         "should be able to admin with folder read and write and wildcard scope",
			dashIsFolder: true,
			dashUID:      "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersPermissionsRead,
					Scope:  "folders:*",
				},
				{
					Action: dashboards.ActionFoldersPermissionsWrite,
					Scope:  "folders:*",
				},
			},
			expected: true,
		},
		{
			desc:         "should not be able to admin with folder read and wildcard scope",
			dashIsFolder: true,
			dashUID:      "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersPermissionsRead,
					Scope:  "folders:*",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to admin with folder write and wildcard scope",
			dashIsFolder: true,
			dashUID:      "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersPermissionsWrite,
					Scope:  "folders:*",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to admin with folder read and write and dashboard scope",
			dashIsFolder: true,
			dashUID:      "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersPermissionsRead,
					Scope:  "dashboards:uid:1",
				},
				{
					Action: dashboards.ActionFoldersPermissionsWrite,
					Scope:  "dashboards:uid:1",
				},
			},
			expected: false,
		},
		/*
			{
				desc:         "should be able to admin with folder read and write and folder scope",
				dashUID:      "1",
				dashIsFolder: true,
				permissions: []accesscontrol.Permission{
					{
						Action: dashboards.ActionFoldersPermissionsRead,
						Scope:  "folders:uid:1",
					},
					{
						Action: dashboards.ActionFoldersPermissionsWrite,
						Scope:  "folders:uid:1",
					},
				},
				expected: true,
			},
		*/
		{
			desc:         "should not be able to admin with folder read and folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersPermissionsRead,
					Scope:  "folders:uid:1",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to admin with folder write and folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersPermissionsWrite,
					Scope:  "folders:uid:1",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to admin with folder read and write and incorrect dashboard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersPermissionsRead,
					Scope:  "dashboards:uid:10",
				},
				{
					Action: dashboards.ActionFoldersPermissionsWrite,
					Scope:  "dashboards:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to admin with folder read and write and incorrect folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersPermissionsRead,
					Scope:  "folders:uid:10",
				},
				{
					Action: dashboards.ActionFoldersPermissionsWrite,
					Scope:  "folders:uid:10",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			guardian := setupAccessControlGuardianTest(t, tt.dashUID, tt.permissions, nil, nil, nil, tt.dashIsFolder)

			can, err := guardian.CanAdmin()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, can)
		})
	}
}

func TestAccessControlDashboardGuardian_CanDelete(t *testing.T) {
	tests := []accessControlGuardianTestCase{
		{
			desc:    "should be able to delete with dashboard wildcard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsDelete,
					Scope:  "dashboards:*",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to delete with folder wildcard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsDelete,
					Scope:  "folders:*",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to delete with dashboard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsDelete,
					Scope:  "dashboards:uid:1",
				},
			},
			expected: true,
		},
		{
			desc:    "should be able to delete with folder general scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsDelete,
					Scope:  "folders:uid:general",
				},
			},
			expected: true,
		},
		/*
			{
				desc:    "should be able to delete with folder scope",
				dashUID: "1",
				permissions: []accesscontrol.Permission{
					{
						Action: dashboards.ActionDashboardsDelete,
						Scope:  "folders:uid:1",
					},
				},
				expected: true,
			},
		*/
		{
			desc:    "should not be able to delete with incorrect dashboard scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsDelete,
					Scope:  "dashboards:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:    "should not be able to delete with incorrect folder scope",
			dashUID: "1",
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionDashboardsDelete,
					Scope:  "folders:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to delete with folder delete and dashboard wildcard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersDelete,
					Scope:  "dashboards:*",
				},
			},
			expected: false,
		},
		{
			desc:         "should be able to delete with folder deletea and folder wildcard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersDelete,
					Scope:  "folders:*",
				},
			},
			expected: true,
		},
		{
			desc:         "should not be able to delete with folder delete and dashboard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersDelete,
					Scope:  "dashboards:uid:1",
				},
			},
			expected: false,
		},
		{
			desc:         "should be able to delete with folder delete and folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersDelete,
					Scope:  "folders:uid:1",
				},
			},
			expected: true,
		},
		{
			desc:         "should not be able to delete with folder delete and incorrect dashboard scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersDelete,
					Scope:  "dashboards:uid:10",
				},
			},
			expected: false,
		},
		{
			desc:         "should not be able to delete with folder delete and incorrect folder scope",
			dashUID:      "1",
			dashIsFolder: true,
			permissions: []accesscontrol.Permission{
				{
					Action: dashboards.ActionFoldersDelete,
					Scope:  "folders:uid:10",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			guardian := setupAccessControlGuardianTest(t, tt.dashUID, tt.permissions, nil, nil, nil, tt.dashIsFolder)

			can, err := guardian.CanDelete()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, can)
		})
	}
}

type accessControlGuardianCanCreateTestCase struct {
	desc        string
	isFolder    bool
	folderID    int64
	permissions []accesscontrol.Permission
	expected    bool
}

func TestAccessControlDashboardGuardian_CanCreate(t *testing.T) {
	tests := []accessControlGuardianCanCreateTestCase{
		{
			desc:     "should be able to create dashboard in general folder",
			isFolder: false,
			folderID: 0,
			permissions: []accesscontrol.Permission{
				{Action: dashboards.ActionDashboardsCreate, Scope: "folders:uid:general"},
			},
			expected: true,
		},
		{
			desc:     "should be able to create dashboard in any folder",
			isFolder: false,
			folderID: 0,
			permissions: []accesscontrol.Permission{
				{Action: dashboards.ActionDashboardsCreate, Scope: "folders:*"},
			},
			expected: true,
		},
		{
			desc:        "should not be able to create dashboard without permissions",
			isFolder:    false,
			folderID:    0,
			permissions: []accesscontrol.Permission{},
			expected:    false,
		},
		{
			desc:     "should be able to create folder with correct permissions",
			isFolder: true,
			folderID: 0,
			permissions: []accesscontrol.Permission{
				{Action: dashboards.ActionFoldersCreate},
			},
			expected: true,
		},
		{
			desc:        "should not be able to create folders without permissions",
			isFolder:    true,
			folderID:    0,
			permissions: []accesscontrol.Permission{},
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			guardian := setupAccessControlGuardianTest(t, "0", tt.permissions, nil, nil, nil, tt.isFolder)

			can, err := guardian.CanCreate(tt.folderID, tt.isFolder)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, can)
		})
	}
}

type accessControlGuardianGetHiddenACLTestCase struct {
	desc        string
	permissions []accesscontrol.ResourcePermission
	hiddenUsers map[string]struct{}
	isFolder    bool
}

func TestAccessControlDashboardGuardian_GetHiddenACL(t *testing.T) {
	tests := []accessControlGuardianGetHiddenACLTestCase{
		{
			desc: "should only return permissions containing hidden users",
			permissions: []accesscontrol.ResourcePermission{
				{RoleName: "managed:users:1:permissions", UserId: 1, UserLogin: "user1", IsManaged: true},
				{RoleName: "managed:teams:1:permissions", TeamId: 1, Team: "team1", IsManaged: true},
				{RoleName: "managed:users:2:permissions", UserId: 2, UserLogin: "user2", IsManaged: true},
				{RoleName: "managed:users:3:permissions", UserId: 3, UserLogin: "user3", IsManaged: true},
				{RoleName: "managed:users:4:permissions", UserId: 4, UserLogin: "user4", IsManaged: true},
			},
			hiddenUsers: map[string]struct{}{"user2": {}, "user3": {}},
		},
		{
			desc: "should only return permissions containing hidden users",
			permissions: []accesscontrol.ResourcePermission{
				{RoleName: "managed:users:1:permissions", UserId: 1, UserLogin: "user1", IsManaged: true},
				{RoleName: "managed:teams:1:permissions", TeamId: 1, Team: "team1", IsManaged: true},
				{RoleName: "managed:users:2:permissions", UserId: 2, UserLogin: "user2", IsManaged: true},
				{RoleName: "managed:users:3:permissions", UserId: 3, UserLogin: "user3", IsManaged: true},
				{RoleName: "managed:users:4:permissions", UserId: 4, UserLogin: "user4", IsManaged: true},
			},
			hiddenUsers: map[string]struct{}{"user2": {}, "user3": {}},
			isFolder:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mocked := accesscontrolmock.NewMockedPermissionsService()
			mocked.On("MapActions", mock.Anything).Return("View")
			mocked.On("GetPermissions", mock.Anything, mock.Anything, mock.Anything).Return(tt.permissions, nil)
			guardian := setupAccessControlGuardianTest(t, "1", nil, nil, mocked, mocked, tt.isFolder)

			cfg := setting.NewCfg()
			cfg.HiddenUsers = tt.hiddenUsers
			permissions, err := guardian.GetHiddenACL(cfg)
			require.NoError(t, err)
			var hiddenUserNames []string
			for name := range tt.hiddenUsers {
				hiddenUserNames = append(hiddenUserNames, name)
			}
			assert.Len(t, permissions, len(hiddenUserNames))
			for _, p := range permissions {
				assert.Contains(t, hiddenUserNames, fmt.Sprintf("user%d", p.UserID))
			}
		})
	}
}

func setupAccessControlGuardianTest(t *testing.T, uid string,
	permissions []accesscontrol.Permission,
	cfg *setting.Cfg,
	dashboardPermissions accesscontrol.DashboardPermissionsService, folderPermissions accesscontrol.FolderPermissionsService, isFolder bool) DashboardGuardian {
	t.Helper()
	store := db.InitTestDB(t)

	toSave := dashboards.NewDashboard(uid)
	toSave.SetUID(uid)

	// seed dashboard
	quotaService := quotatest.New(false, nil)
	dashStore, err := dashdb.ProvideDashboardStore(store, store.Cfg, featuremgmt.WithFeatures(), tagimpl.ProvideService(store, store.Cfg), quotaService)
	require.NoError(t, err)
	dash, err := dashStore.SaveDashboard(context.Background(), dashboards.SaveDashboardCommand{
		Dashboard: toSave.Data,
		UserID:    1,
		OrgID:     1,
		IsFolder:  isFolder,
	})
	require.NoError(t, err)
	fakeDashboardService := dashboards.NewFakeDashboardService(t)
	qResult := &dashboards.Dashboard{}
	fakeDashboardService.On("GetDashboard", mock.Anything, mock.AnythingOfType("*dashboards.GetDashboardQuery")).Run(func(args mock.Arguments) {
		q := args.Get(1).(*dashboards.GetDashboardQuery)
		qResult = &dashboards.Dashboard{
			ID:       q.ID,
			UID:      q.UID,
			OrgID:    q.OrgID,
			IsFolder: isFolder,
		}
	}).Maybe().Return(qResult, nil)

	ac := accesscontrolmock.New().WithPermissions(permissions)
	// TODO replace with actual folder store implementation after resolving import cycles
	//folderSvc := folderimpl.ProvideService(ac, bus.ProvideBus(tracing.InitializeTracerForTest()), cfg, dashStore, folderimpl.ProvideDashboardFolderStore(store), store, featuremgmt.WithFeatures())
	folderSvc := foldertest.NewFakeService()
	folderStore := foldertest.NewFakeFolderStore(t)
	ac.RegisterScopeAttributeResolver(dashboards.NewDashboardUIDScopeResolver(folderStore, fakeDashboardService, folderSvc))
	ac.RegisterScopeAttributeResolver(dashboards.NewFolderUIDScopeResolver(folderSvc))
	ac.RegisterScopeAttributeResolver(dashboards.NewFolderIDScopeResolver(folderStore, folderSvc))

	license := licensingtest.NewFakeLicensing()
	license.On("FeatureEnabled", "accesscontrol.enforcement").Return(true).Maybe()
	teamSvc := teamimpl.ProvideService(store, store.Cfg)
	userSvc, err := userimpl.ProvideService(store, nil, store.Cfg, nil, nil, quotatest.New(false, nil), supportbundlestest.NewFakeBundleService())
	require.NoError(t, err)

	if folderPermissions == nil {
		folderPermissions, err = ossaccesscontrol.ProvideFolderPermissions(
			cfg, routing.NewRouteRegister(), store, ac, license, &dashboards.FakeDashboardStore{}, folderSvc, ac, teamSvc, userSvc)
		require.NoError(t, err)
	}
	if dashboardPermissions == nil {
		dashboardPermissions, err = ossaccesscontrol.ProvideDashboardPermissions(
			cfg, routing.NewRouteRegister(), store, ac, license, &dashboards.FakeDashboardStore{}, folderSvc, ac, teamSvc, userSvc)
		require.NoError(t, err)
	}

	g, err := NewAccessControlDashboardGuardianByDashboard(context.Background(), cfg, dash, &user.SignedInUser{OrgID: 1}, store, ac, folderPermissions, dashboardPermissions, fakeDashboardService)
	require.NoError(t, err)
	return g
}
