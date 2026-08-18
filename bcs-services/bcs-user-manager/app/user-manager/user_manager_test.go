package usermanager

import "testing"

func TestRequiresLocklessIAMMigration(t *testing.T) {
	tests := []struct {
		dbType   string
		expected bool
	}{
		{dbType: "ob", expected: true},
		{dbType: "oceanbase", expected: true},
		{dbType: "dg", expected: true},
		{dbType: "goldendb", expected: true},
		{dbType: "GoldenDB", expected: true},
		{dbType: "mysql", expected: false},
		{dbType: "postgres", expected: false},
		{dbType: "gaussdb", expected: false},
		{dbType: "opengauss", expected: false},
	}

	for _, test := range tests {
		t.Run(test.dbType, func(t *testing.T) {
			if actual := requiresLocklessIAMMigration(test.dbType); actual != test.expected {
				t.Fatalf("requiresLocklessIAMMigration(%q) = %t, want %t", test.dbType, actual, test.expected)
			}
		})
	}
}

func TestIAMMigrationDatabaseType(t *testing.T) {
	tests := []struct {
		dbType   string
		expected string
	}{
		{dbType: "postgres", expected: "postgres"},
		{dbType: "PostgreSQL", expected: "postgres"},
		{dbType: "gaussdb", expected: "gaussdb"},
		{dbType: "OpenGauss", expected: "opengauss"},
		{dbType: "mysql", expected: ""},
		{dbType: "oceanbase", expected: ""},
		{dbType: "", expected: ""},
	}

	for _, test := range tests {
		t.Run(test.dbType, func(t *testing.T) {
			if actual := iamMigrationDatabaseType(test.dbType); actual != test.expected {
				t.Fatalf("iamMigrationDatabaseType(%q) = %q, want %q", test.dbType, actual, test.expected)
			}
		})
	}
}
