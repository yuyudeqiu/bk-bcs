package framework

import "testing"

func TestIsAllowedMySQLTestDBName(t *testing.T) {
	t.Setenv("BCS_TEST_MYSQL_DBNAME_PREFIX", "")
	if !isAllowedMySQLTestDBName("bcs_user_test") {
		t.Fatal("default sqlstore test database should be allowed")
	}
	if !isAllowedMySQLTestDBName("bcs_user_iam_migration_test") {
		t.Fatal("default iam migration test database should be allowed")
	}
	if isAllowedMySQLTestDBName("bcs_user") {
		t.Fatal("non-test database should not be allowed")
	}

	t.Setenv("BCS_TEST_MYSQL_DBNAME_PREFIX", "custom_test_, another_test_, ")
	if !isAllowedMySQLTestDBName("custom_test_oceanbase") {
		t.Fatal("configured test database prefix should be allowed")
	}
	if !isAllowedMySQLTestDBName("another_test_mysql") {
		t.Fatal("each configured test database prefix should be allowed")
	}
	if isAllowedMySQLTestDBName("unrelated_database") {
		t.Fatal("empty configured prefix should not allow arbitrary database names")
	}
}
