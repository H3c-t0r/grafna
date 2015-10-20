package api

import (
	l "log"
	"github.com/toolkits/file"
	"flag"
	"sync"
	"encoding/json"

	"github.com/Unknwon/macaron"
	"github.com/Cepave/grafana/pkg/api/dtos"
	"github.com/Cepave/grafana/pkg/middleware"
	m "github.com/Cepave/grafana/pkg/models"
	"github.com/macaron-contrib/binding"
)

var OpenFalconConfigFile = flag.String("configGlobal", "cfg.json", "configuration file")
type DatabaseConfig struct {
	Addr       string	`json:"addr"`
	Account    string	`json:"account"`
	Password   string	`json:"password"`
}

type DatasourceConfig struct {
	Type			string	`json:"type"`
	UrlDashboard	string	`json:"urlDashboard"`
	UrlQuery		string	`json:"urlQuery"`
}

type GlobalConfig struct {
	Database      *DatabaseConfig    `json:"database"`
	Datasource    *DatasourceConfig  `json:"datasource"`
}

var (
	configOpenFalcon     *GlobalConfig
	lock = new(sync.RWMutex)
)

/**
 * @function name:	func parseConfig(cfg string)
 * @description:	This function parses config file cfg.json.
 * @related issues:	OWL-115, OWL-085
 * @param:			cfg string
 * @return:			void
 * @author:			Don Hsieh
 * @since:			09/14/2015
 * @last modified: 	10/07/2015
 * @called by:		func main()
 */
func parseConfig(cfg string) {
	if !file.IsExist(cfg) {
		l.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}
	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		l.Fatalln("read config file:", cfg, "fail:", err)
	}

	var configGlobal GlobalConfig
	err = json.Unmarshal([]byte(configContent), &configGlobal)
	if err != nil {
		l.Fatalln("parse config file:", cfg, "fail:", err)
		return
	}
	lock.Lock()
	defer lock.Unlock()
	configOpenFalcon = &configGlobal
}

// Register adds http routes
func Register(r *macaron.Macaron) {
	flag.Parse()
	parseConfig(*OpenFalconConfigFile)

	reqSignedIn := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true})
	reqGrafanaAdmin := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true, ReqGrafanaAdmin: true})
	reqEditorRole := middleware.RoleAuth(m.ROLE_EDITOR, m.ROLE_ADMIN)
	regOrgAdmin := middleware.RoleAuth(m.ROLE_ADMIN)
	bind := binding.Bind

	// not logged in views
	r.Get("/", reqSignedIn, Index)
	r.Get("/logout", Logout)
	r.Post("/login", bind(dtos.LoginCommand{}), wrap(LoginPost))
	r.Get("/login/:name", OAuthLogin)
	r.Get("/login", LoginView)

	// authed views
	r.Get("/profile/", reqSignedIn, Index)
	r.Get("/org/", reqSignedIn, Index)
	r.Get("/org/new", reqSignedIn, Index)
	r.Get("/datasources/", reqSignedIn, Index)
	r.Get("/datasources/edit/*", reqSignedIn, Index)
	r.Get("/org/users/", reqSignedIn, Index)
	r.Get("/org/apikeys/", reqSignedIn, Index)
	r.Get("/dashboard/import/", reqSignedIn, Index)
	r.Get("/admin/settings", reqGrafanaAdmin, Index)
	r.Get("/admin/users", reqGrafanaAdmin, Index)
	r.Get("/admin/users/create", reqGrafanaAdmin, Index)
	r.Get("/admin/users/edit/:id", reqGrafanaAdmin, Index)
	r.Get("/dashboard/*", reqSignedIn, Index)

	// sign up
	r.Get("/signup", Index)
	r.Post("/api/user/signup", bind(m.CreateUserCommand{}), wrap(SignUp))

	// reset password
	r.Get("/user/password/send-reset-email", Index)
	r.Get("/user/password/reset", Index)

	r.Post("/api/user/password/send-reset-email", bind(dtos.SendResetPasswordEmailForm{}), wrap(SendResetPasswordEmail))
	r.Post("/api/user/password/reset", bind(dtos.ResetUserPasswordForm{}), wrap(ResetPassword))

	// dashboard snapshots
	r.Post("/api/snapshots/", bind(m.CreateDashboardSnapshotCommand{}), CreateDashboardSnapshot)
	r.Get("/dashboard/snapshot/*", Index)

	r.Get("/api/snapshots/:key", GetDashboardSnapshot)
	r.Get("/api/snapshots-delete/:key", DeleteDashboardSnapshot)

	// api renew session based on remember cookie
	r.Get("/api/login/ping", LoginApiPing)

	// authed api
	r.Group("/api", func() {

		// user (signed in)
		r.Group("/user", func() {
			r.Get("/", wrap(GetSignedInUser))
			r.Put("/", bind(m.UpdateUserCommand{}), wrap(UpdateSignedInUser))
			r.Post("/using/:id", wrap(UserSetUsingOrg))
			r.Get("/orgs", wrap(GetSignedInUserOrgList))
			r.Post("/stars/dashboard/:id", wrap(StarDashboard))
			r.Delete("/stars/dashboard/:id", wrap(UnstarDashboard))
			r.Put("/password", bind(m.ChangeUserPasswordCommand{}), wrap(ChangeUserPassword))
		})

		// users (admin permission required)
		r.Group("/users", func() {
			r.Get("/", wrap(SearchUsers))
			r.Get("/:id", wrap(GetUserById))
			r.Get("/:id/orgs", wrap(GetUserOrgList))
			r.Put("/:id", bind(m.UpdateUserCommand{}), wrap(UpdateUser))
		}, reqGrafanaAdmin)

		// current org
		r.Group("/org", func() {
			r.Get("/", wrap(GetOrgCurrent))
			r.Put("/", bind(m.UpdateOrgCommand{}), wrap(UpdateOrgCurrent))
			r.Post("/users", bind(m.AddOrgUserCommand{}), wrap(AddOrgUserToCurrentOrg))
			r.Get("/users", wrap(GetOrgUsersForCurrentOrg))
			r.Patch("/users/:userId", bind(m.UpdateOrgUserCommand{}), wrap(UpdateOrgUserForCurrentOrg))
			r.Delete("/users/:userId", wrap(RemoveOrgUserForCurrentOrg))
		}, regOrgAdmin)

		// create new org
		r.Post("/orgs", bind(m.CreateOrgCommand{}), wrap(CreateOrg))

		// search all orgs
		r.Get("/orgs", reqGrafanaAdmin, wrap(SearchOrgs))

		// orgs (admin routes)
		r.Group("/orgs/:orgId", func() {
			r.Put("/", bind(m.UpdateOrgCommand{}), wrap(UpdateOrg))
			r.Get("/users", wrap(GetOrgUsers))
			r.Post("/users", bind(m.AddOrgUserCommand{}), wrap(AddOrgUser))
			r.Patch("/users/:userId", bind(m.UpdateOrgUserCommand{}), wrap(UpdateOrgUser))
			r.Delete("/users/:userId", wrap(RemoveOrgUser))
		}, reqGrafanaAdmin)

		// auth api keys
		r.Group("/auth/keys", func() {
			r.Get("/", wrap(GetApiKeys))
			r.Post("/", bind(m.AddApiKeyCommand{}), wrap(AddApiKey))
			r.Delete("/:id", wrap(DeleteApiKey))
		}, regOrgAdmin)

		// Data sources
		r.Group("/datasources", func() {
			r.Get("/", GetDataSources)
			r.Post("/", bind(m.AddDataSourceCommand{}), AddDataSource)
			r.Put("/:id", bind(m.UpdateDataSourceCommand{}), UpdateDataSource)
			r.Delete("/:id", DeleteDataSource)
			r.Get("/:id", GetDataSourceById)
			r.Get("/plugins", GetDataSourcePlugins)
		}, regOrgAdmin)

		r.Get("/frontend/settings/", GetFrontendSettings)
		r.Any("/datasources/proxy/:id/*", reqSignedIn, ProxyDataSourceRequest)
		r.Any("/datasources/proxy/:id", reqSignedIn, ProxyDataSourceRequest)

		// Dashboard
		r.Group("/dashboards", func() {
			r.Combo("/db/:slug").Get(GetDashboard).Delete(DeleteDashboard)
			r.Post("/db", reqEditorRole, bind(m.SaveDashboardCommand{}), PostDashboard)
			r.Get("/file/:file", GetDashboardFromJsonFile)
			r.Get("/home", GetHomeDashboard)
			r.Get("/tags", GetDashboardTags)
		})

		// Search
		r.Get("/search/", Search)

		// metrics
		r.Get("/metrics/test", GetTestMetrics)
	}, reqSignedIn)

	// admin api
	r.Group("/api/admin", func() {
		r.Get("/settings", AdminGetSettings)
		r.Post("/users", bind(dtos.AdminCreateUserForm{}), AdminCreateUser)
		r.Put("/users/:id/password", bind(dtos.AdminUpdateUserPasswordForm{}), AdminUpdateUserPassword)
		r.Put("/users/:id/permissions", bind(dtos.AdminUpdateUserPermissionsForm{}), AdminUpdateUserPermissions)
		r.Delete("/users/:id", AdminDeleteUser)
	}, reqGrafanaAdmin)

	// rendering
	r.Get("/render/*", reqSignedIn, RenderToPng)

	r.NotFound(NotFoundHandler)
}
