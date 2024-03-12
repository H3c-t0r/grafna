package sqlstore

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/org"
	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/grafana/grafana/pkg/setting"
)

type sqlStoreTest struct {
	name       string
	dbType     string
	dbHost     string
	dbURL      string
	dbUser     string
	dbPwd      string
	expConnStr string
	features   featuremgmt.FeatureToggles
	err        error
}

var sqlStoreTestCases = []sqlStoreTest{
	{
		name:       "MySQL IPv4",
		dbType:     "mysql",
		dbHost:     "1.2.3.4:5678",
		expConnStr: ":@tcp(1.2.3.4:5678)/test_db?collation=utf8mb4_unicode_ci&allowNativePasswords=true&clientFoundRows=true",
	},
	{
		name:       "Postgres IPv4",
		dbType:     "postgres",
		dbHost:     "1.2.3.4:5678",
		expConnStr: "user='' host=1.2.3.4 port=5678 dbname=test_db sslmode='' sslcert='' sslkey='' sslrootcert=''",
	},
	{
		name:       "Postgres IPv4 (Default Port)",
		dbType:     "postgres",
		dbHost:     "1.2.3.4",
		expConnStr: "user='' host=1.2.3.4 port=5432 dbname=test_db sslmode='' sslcert='' sslkey='' sslrootcert=''",
	},
	{
		name:       "Postgres username and password",
		dbType:     "postgres",
		dbHost:     "1.2.3.4",
		dbUser:     "grafana",
		dbPwd:      "password",
		expConnStr: "user=grafana host=1.2.3.4 port=5432 dbname=test_db sslmode='' sslcert='' sslkey='' sslrootcert='' password=password",
	},
	{
		name:       "Postgres username no password",
		dbType:     "postgres",
		dbHost:     "1.2.3.4",
		dbUser:     "grafana",
		dbPwd:      "",
		expConnStr: "user=grafana host=1.2.3.4 port=5432 dbname=test_db sslmode='' sslcert='' sslkey='' sslrootcert=''",
	},
	{
		name:       "MySQL IPv4 (Default Port)",
		dbType:     "mysql",
		dbHost:     "1.2.3.4",
		expConnStr: ":@tcp(1.2.3.4)/test_db?collation=utf8mb4_unicode_ci&allowNativePasswords=true&clientFoundRows=true",
	},
	{
		name:       "MySQL IPv6",
		dbType:     "mysql",
		dbHost:     "[fe80::24e8:31b2:91df:b177]:1234",
		expConnStr: ":@tcp([fe80::24e8:31b2:91df:b177]:1234)/test_db?collation=utf8mb4_unicode_ci&allowNativePasswords=true&clientFoundRows=true",
	},
	{
		name:       "Postgres IPv6",
		dbType:     "postgres",
		dbHost:     "[fe80::24e8:31b2:91df:b177]:1234",
		expConnStr: "user='' host=fe80::24e8:31b2:91df:b177 port=1234 dbname=test_db sslmode='' sslcert='' sslkey='' sslrootcert=''",
	},
	{
		name:       "MySQL IPv6 (Default Port)",
		dbType:     "mysql",
		dbHost:     "[::1]",
		expConnStr: ":@tcp([::1])/test_db?collation=utf8mb4_unicode_ci&allowNativePasswords=true&clientFoundRows=true",
	},
	{
		name:       "Postgres IPv6 (Default Port)",
		dbType:     "postgres",
		dbHost:     "[::1]",
		expConnStr: "user='' host=::1 port=5432 dbname=test_db sslmode='' sslcert='' sslkey='' sslrootcert=''",
	},
	{
		name:  "Invalid database URL",
		dbURL: "://invalid.com/",
		err:   &url.Error{Op: "parse", URL: "://invalid.com/", Err: errors.New("missing protocol scheme")},
	},
	{
		name:       "MySQL with ANSI_QUOTES mode",
		dbType:     "mysql",
		dbHost:     "[::1]",
		features:   featuremgmt.WithFeatures(featuremgmt.FlagMysqlAnsiQuotes),
		expConnStr: ":@tcp([::1])/test_db?collation=utf8mb4_unicode_ci&allowNativePasswords=true&clientFoundRows=true&sql_mode='ANSI_QUOTES'",
	},
}

func TestIntegrationSQLConnectionString(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	for _, testCase := range sqlStoreTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			sqlstore := &SQLStore{}
			sqlstore.Cfg = makeSQLStoreTestConfig(t, testCase)
			connStr, err := sqlstore.buildConnectionString()
			require.Equal(t, testCase.err, err)

			assert.Equal(t, testCase.expConnStr, connStr)
		})
	}
}

func TestIntegrationIsUniqueConstraintViolation(t *testing.T) {
	store := InitTestDB(t)

	testCases := []struct {
		desc string
		f    func(*testing.T, *DBSession) error
	}{
		{
			desc: "successfully detect primary key violations",
			f: func(t *testing.T, sess *DBSession) error {
				// Attempt to insert org with provided ID (primary key) twice
				now := time.Now()
				org := org.Org{Name: "test org primary key violation", Created: now, Updated: now, ID: 42}
				err := sess.InsertId(&org, store.Dialect)
				require.NoError(t, err)

				// Provide a different name to avoid unique constraint violation
				org.Name = "test org 2"
				return sess.InsertId(&org, store.Dialect)
			},
		},
		{
			desc: "successfully detect unique constrain violations",
			f: func(t *testing.T, sess *DBSession) error {
				// Attempt to insert org with reserved name
				now := time.Now()
				org := org.Org{Name: "test org unique constrain violation", Created: now, Updated: now, ID: 43}
				err := sess.InsertId(&org, store.Dialect)
				require.NoError(t, err)

				// Provide a different ID to avoid primary key violation
				org.ID = 44
				return sess.InsertId(&org, store.Dialect)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := store.WithDbSession(context.Background(), func(sess *DBSession) error {
				return tc.f(t, sess)
			})
			require.Error(t, err)
			assert.True(t, store.Dialect.IsUniqueConstraintViolation(err))
		})
	}
}

func makeSQLStoreTestConfig(t *testing.T, tc sqlStoreTest) *setting.Cfg {
	t.Helper()

	if tc.features == nil {
		tc.features = featuremgmt.WithFeatures()
	}
	// nolint:staticcheck
	cfg := setting.NewCfgWithFeatures(tc.features.IsEnabledGlobally)

	sec, err := cfg.Raw.NewSection("database")
	require.NoError(t, err)
	_, err = sec.NewKey("type", tc.dbType)
	require.NoError(t, err)
	_, err = sec.NewKey("host", tc.dbHost)
	require.NoError(t, err)
	_, err = sec.NewKey("url", tc.dbURL)
	require.NoError(t, err)
	_, err = sec.NewKey("user", tc.dbUser)
	require.NoError(t, err)
	_, err = sec.NewKey("name", "test_db")
	require.NoError(t, err)
	_, err = sec.NewKey("password", tc.dbPwd)
	require.NoError(t, err)

	return cfg
}

func TestBuildConnectionStringPostgres(t *testing.T) {
	testCases := []struct {
		name            string
		SslMode         string
		SSLSNI          string
		CaCertPath      string
		ClientKeyPath   string
		ClientCertPath  string
		expectedConnStr string
	}{
		{
			name:            "Postgres with sslmode disable",
			SslMode:         "disable",
			expectedConnStr: "user=grafana host=127.0.0.1 port=5432 dbname=grafana_test sslmode=disable sslcert='' sslkey='' sslrootcert='' password=password",
		},
		{
			name:            "Postgres with sslmode verify-ca",
			SslMode:         "verify-ca",
			CaCertPath:      "/path/to/ca_cert",
			ClientKeyPath:   "/path/to/client_key",
			ClientCertPath:  "/path/to/client_cert",
			expectedConnStr: "user=grafana host=127.0.0.1 port=5432 dbname=grafana_test sslmode=verify-ca sslcert=/path/to/client_cert sslkey=/path/to/client_key sslrootcert=/path/to/ca_cert password=password",
		},
		{
			name:            "Postgres with sslmode verify-ca without SNI",
			SslMode:         "verify-ca",
			CaCertPath:      "/path/to/ca_cert",
			ClientKeyPath:   "/path/to/client_key",
			ClientCertPath:  "/path/to/client_cert",
			SSLSNI:          "0",
			expectedConnStr: "user=grafana host=127.0.0.1 port=5432 dbname=grafana_test sslmode=verify-ca sslcert=/path/to/client_cert sslkey=/path/to/client_key sslrootcert=/path/to/ca_cert sslsni=0 password=password",
		},
		{
			name:            "Postgres with sslmode verify-ca with SNI",
			SslMode:         "verify-ca",
			CaCertPath:      "/path/to/ca_cert",
			ClientKeyPath:   "/path/to/client_key",
			ClientCertPath:  "/path/to/client_cert",
			SSLSNI:          "1",
			expectedConnStr: "user=grafana host=127.0.0.1 port=5432 dbname=grafana_test sslmode=verify-ca sslcert=/path/to/client_cert sslkey=/path/to/client_key sslrootcert=/path/to/ca_cert sslsni=1 password=password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlstore := &SQLStore{}
			sqlstore.Cfg = setting.NewCfg()
			sec, err := sqlstore.Cfg.Raw.NewSection("database")
			require.NoError(t, err)
			_, err = sec.NewKey("type", migrator.Postgres)
			require.NoError(t, err)
			_, err = sec.NewKey("host", "127.0.0.1")
			require.NoError(t, err)
			_, err = sec.NewKey("port", "5432")
			require.NoError(t, err)
			_, err = sec.NewKey("user", "grafana")
			require.NoError(t, err)
			_, err = sec.NewKey("name", "grafana_test")
			require.NoError(t, err)
			_, err = sec.NewKey("password", "password")
			require.NoError(t, err)
			_, err = sec.NewKey("ssl_mode", tc.SslMode)
			require.NoError(t, err)
			_, err = sec.NewKey("ca_cert_path", tc.CaCertPath)
			require.NoError(t, err)
			_, err = sec.NewKey("client_key_path", tc.ClientKeyPath)
			require.NoError(t, err)
			_, err = sec.NewKey("client_cert_path", tc.ClientCertPath)
			require.NoError(t, err)
			_, err = sec.NewKey("ssl_sni", tc.SSLSNI)
			require.NoError(t, err)
			connectionString, err := sqlstore.buildConnectionString()
			require.NoError(t, err)
			assert.Equal(t, tc.expectedConnStr, connectionString)
		})
	}
}
