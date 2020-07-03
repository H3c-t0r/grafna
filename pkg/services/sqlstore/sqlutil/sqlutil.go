package sqlutil

type TestDB struct {
	DriverName string
	ConnStr    string
}

// To run all tests in a local test database, set ConnStr to "grafana_test.db"
//var TestDB_Sqlite3 = TestDB{DriverName: "sqlite3", ConnStr: ":memory:"}

var TestDB_Sqlite3 = TestDB{DriverName: "sqlite3", ConnStr: "grafana_test.db"}
var TestDB_Mysql = TestDB{DriverName: "mysql", ConnStr: "grafana:password@tcp(localhost:3306)/grafana_tests?collation=utf8mb4_unicode_ci"}
var TestDB_Postgres = TestDB{DriverName: "postgres", ConnStr: "user=grafanatest password=grafanatest host=localhost port=5432 dbname=grafanatest sslmode=disable"}
var TestDB_Mssql = TestDB{DriverName: "mssql", ConnStr: "server=localhost;port=1433;database=grafanatest;user id=grafana;password=Password!"}
