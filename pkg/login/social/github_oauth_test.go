package social

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/setting"
)

const testGHUserTeamsJSON = `[
  {
    "id": 1,
    "node_id": "MDQ6VGVhbTE=",
    "url": "https://api.github.com/teams/1",
    "html_url": "https://github.com/orgs/github/teams/justice-league",
    "name": "Justice League",
    "slug": "justice-league",
    "description": "A great team.",
    "privacy": "closed",
    "permission": "admin",
    "members_url": "https://api.github.com/teams/1/members{/member}",
    "repositories_url": "https://api.github.com/teams/1/repos",
    "parent": null,
    "members_count": 3,
    "repos_count": 10,
    "created_at": "2017-07-14T16:53:42Z",
    "updated_at": "2017-08-17T12:37:15Z",
    "organization": {
      "login": "github",
      "id": 1,
      "node_id": "MDEyOk9yZ2FuaXphdGlvbjE=",
      "url": "https://api.github.com/orgs/github",
      "repos_url": "https://api.github.com/orgs/github/repos",
      "events_url": "https://api.github.com/orgs/github/events",
      "hooks_url": "https://api.github.com/orgs/github/hooks",
      "issues_url": "https://api.github.com/orgs/github/issues",
      "members_url": "https://api.github.com/orgs/github/members{/member}",
      "public_members_url": "https://api.github.com/orgs/github/public_members{/member}",
      "avatar_url": "https://github.com/images/error/octocat_happy.gif",
      "description": "A great organization",
      "name": "github",
      "company": "GitHub",
      "blog": "https://github.com/blog",
      "location": "San Francisco",
      "email": "octocat@github.com",
      "is_verified": true,
      "has_organization_projects": true,
      "has_repository_projects": true,
      "public_repos": 2,
      "public_gists": 1,
      "followers": 20,
      "following": 0,
      "html_url": "https://github.com/octocat",
      "created_at": "2008-01-14T04:33:35Z",
      "updated_at": "2017-08-17T12:37:15Z",
      "type": "Organization"
    }
  }
]`

const testGHUserJSON = `{
  "login": "octocat",
  "id": 1,
  "node_id": "MDQ6VXNlcjE=",
  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
  "gravatar_id": "",
  "url": "https://api.github.com/users/octocat",
  "html_url": "https://github.com/octocat",
  "followers_url": "https://api.github.com/users/octocat/followers",
  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
  "organizations_url": "https://api.github.com/users/octocat/orgs",
  "repos_url": "https://api.github.com/users/octocat/repos",
  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
  "received_events_url": "https://api.github.com/users/octocat/received_events",
  "type": "User",
  "site_admin": false,
  "name": "monalisa octocat",
  "company": "GitHub",
  "blog": "https://github.com/blog",
  "location": "San Francisco",
  "email": "octocat@github.com",
  "hireable": false,
  "bio": "There once was...",
  "twitter_username": "monatheoctocat",
  "public_repos": 2,
  "public_gists": 1,
  "followers": 20,
  "following": 0,
  "created_at": "2008-01-14T04:33:35Z",
  "updated_at": "2008-01-14T04:33:35Z",
  "private_gists": 81,
  "total_private_repos": 100,
  "owned_private_repos": 100,
  "disk_usage": 10000,
  "collaborators": 8,
  "two_factor_authentication": true,
  "plan": {
    "name": "Medium",
    "space": 400,
    "private_repos": 20,
    "collaborators": 0
  }
}`

func TestSocialGitHub_UserInfo(t *testing.T) {
	var boolPointer *bool
	tests := []struct {
		name                     string
		userRawJSON              string
		userTeamsRawJSON         string
		settingAutoAssignOrgRole string
		settingAllowGrafanaAdmin bool
		settingSkipOrgRoleSync   bool
		roleAttributePath        string
		autoAssignOrgRole        string
		want                     *BasicUserInfo
		wantErr                  bool
	}{
		{
			name:              "Basic User info",
			userRawJSON:       testGHUserJSON,
			userTeamsRawJSON:  testGHUserTeamsJSON,
			autoAssignOrgRole: "",
			roleAttributePath: "",
			want: &BasicUserInfo{
				Id:     "1",
				Name:   "monalisa octocat",
				Email:  "octocat@github.com",
				Login:  "octocat",
				Role:   "Viewer",
				Groups: []string{"https://github.com/orgs/github/teams/justice-league", "@github/justice-league"},
			},
		},
		{
			name:              "Admin mapping takes precedence over auto assign org role",
			roleAttributePath: "[login==octocat] && 'Admin' || 'Viewer'",
			userRawJSON:       testGHUserJSON,
			autoAssignOrgRole: "Editor",
			userTeamsRawJSON:  testGHUserTeamsJSON,
			want: &BasicUserInfo{
				Id:     "1",
				Name:   "monalisa octocat",
				Email:  "octocat@github.com",
				Login:  "octocat",
				Role:   "Admin",
				Groups: []string{"https://github.com/orgs/github/teams/justice-league", "@github/justice-league"},
			},
		},
		{
			name:              "Editor mapping via groups",
			roleAttributePath: "contains(groups[*], '@github/justice-league') && 'Editor' || 'Viewer'",
			userRawJSON:       testGHUserJSON,
			autoAssignOrgRole: "Editor",
			userTeamsRawJSON:  testGHUserTeamsJSON,
			want: &BasicUserInfo{
				Id:     "1",
				Name:   "monalisa octocat",
				Email:  "octocat@github.com",
				Login:  "octocat",
				Role:   "Editor",
				Groups: []string{"https://github.com/orgs/github/teams/justice-league", "@github/justice-league"},
			},
		},
		{
			name:                   "Should be empty role if setting skipOrgRoleSync is set to true",
			roleAttributePath:      "contains(groups[*], '@github/justice-league') && 'Editor' || 'Viewer'",
			settingSkipOrgRoleSync: true,
			userRawJSON:            testGHUserJSON,
			userTeamsRawJSON:       testGHUserTeamsJSON,
			want: &BasicUserInfo{
				Id:     "1",
				Name:   "monalisa octocat",
				Email:  "octocat@github.com",
				Login:  "octocat",
				Role:   "",
				Groups: []string{"https://github.com/orgs/github/teams/justice-league", "@github/justice-league"},
			},
		},
		{
			name:                     "Should return nil pointer if allowGrafanaAdmin and skipOrgRoleSync setting is set to true",
			roleAttributePath:        "contains(groups[*], '@github/justice-league') && 'Editor' || 'Viewer'",
			settingSkipOrgRoleSync:   true,
			settingAllowGrafanaAdmin: true,
			userRawJSON:              testGHUserJSON,
			userTeamsRawJSON:         testGHUserTeamsJSON,
			want: &BasicUserInfo{
				Id:             "1",
				Name:           "monalisa octocat",
				Email:          "octocat@github.com",
				Login:          "octocat",
				Role:           "",
				Groups:         []string{"https://github.com/orgs/github/teams/justice-league", "@github/justice-league"},
				IsGrafanaAdmin: boolPointer,
			},
		},
		{
			name:              "fallback to default org role",
			roleAttributePath: "",
			userRawJSON:       testGHUserJSON,
			autoAssignOrgRole: "Editor",
			userTeamsRawJSON:  testGHUserTeamsJSON,
			want: &BasicUserInfo{
				Id:     "1",
				Name:   "monalisa octocat",
				Email:  "octocat@github.com",
				Login:  "octocat",
				Role:   "Editor",
				Groups: []string{"https://github.com/orgs/github/teams/justice-league", "@github/justice-league"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				// return JSON if matches user endpoint
				if strings.HasSuffix(request.URL.String(), "/user") {
					writer.Header().Set("Content-Type", "application/json")
					_, err := writer.Write([]byte(tt.userRawJSON))
					require.NoError(t, err)
				} else if strings.HasSuffix(request.URL.String(), "/user/teams?per_page=100") {
					writer.Header().Set("Content-Type", "application/json")
					_, err := writer.Write([]byte(tt.userTeamsRawJSON))
					require.NoError(t, err)
				} else {
					writer.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			s, err := NewGitHubProvider(map[string]any{
				"allowed_organizations": "",
				"api_url":               server.URL + "/user",
				"team_ids":              "",
				"role_attribute_path":   tt.roleAttributePath,
			}, &setting.Cfg{
				AutoAssignOrgRole:     tt.autoAssignOrgRole,
				GitHubSkipOrgRoleSync: tt.settingSkipOrgRoleSync,
			}, featuremgmt.WithFeatures())
			require.NoError(t, err)

			token := &oauth2.Token{
				AccessToken: "fake_token",
			}

			got, err := s.UserInfo(context.Background(), server.Client(), token)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSocialGitHub_InitializeExtraFields(t *testing.T) {
	type settingFields struct {
		teamIds              []int
		allowedOrganizations []string
	}
	testCases := []struct {
		name     string
		settings map[string]any
		want     settingFields
	}{
		{
			name: "teamIds is set",
			settings: map[string]any{
				"team_ids": "1234,5678",
			},
			want: settingFields{
				teamIds:              []int{1234, 5678},
				allowedOrganizations: []string{},
			},
		},
		{
			name: "allowedOrganizations is set",
			settings: map[string]any{
				"allowed_organizations": "uuid-1234,uuid-5678",
			},
			want: settingFields{
				teamIds:              []int{},
				allowedOrganizations: []string{"uuid-1234", "uuid-5678"},
			},
		},
		{
			name: "teamIds and allowedOrganizations are empty",
			settings: map[string]any{
				"team_ids":              "",
				"allowed_organizations": "",
			},
			want: settingFields{
				teamIds:              []int{},
				allowedOrganizations: []string{},
			},
		},
		{
			name: "should not error when teamIds are not integers",
			settings: map[string]any{
				"team_ids": "abc1234,5678",
			},
			want: settingFields{
				teamIds:              []int{},
				allowedOrganizations: []string{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := NewGitHubProvider(tc.settings, &setting.Cfg{}, featuremgmt.WithFeatures())
			require.NoError(t, err)

			require.Equal(t, tc.want.teamIds, s.teamIds)
			require.Equal(t, tc.want.allowedOrganizations, s.allowedOrganizations)
		})
	}
}
