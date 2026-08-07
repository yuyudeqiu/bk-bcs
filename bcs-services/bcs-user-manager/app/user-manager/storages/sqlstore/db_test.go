package sqlstore

import (
	"testing"

	godbsdk "code.cwoa.net/carlchen2/cw-godb-sdk/core/config"
)

func TestSDKDatabaseType(t *testing.T) {
	tests := []struct {
		name     string
		dbType   string
		expected godbsdk.DatabaseType
	}{
		{name: "mysql", dbType: "mysql", expected: godbsdk.Mysql},
		{name: "oceanbase short name", dbType: "ob", expected: godbsdk.Mysql},
		{name: "oceanbase", dbType: "oceanbase", expected: godbsdk.Mysql},
		{name: "goldendb short name", dbType: "dg", expected: godbsdk.Mysql},
		{name: "goldendb", dbType: "goldendb", expected: godbsdk.Mysql},
		{name: "goldendb uppercase", dbType: "GoldenDB", expected: godbsdk.Mysql},
		{name: "postgres", dbType: "postgres", expected: godbsdk.Postgres},
		{name: "dameng", dbType: "dameng", expected: godbsdk.Dameng},
		{name: "gaussdb", dbType: "gaussdb", expected: godbsdk.Gaussdb},
		{name: "opengauss", dbType: "opengauss", expected: godbsdk.Gaussdb},
		{name: "opengauss uppercase", dbType: "OpenGauss", expected: godbsdk.Gaussdb},
		{name: "default", dbType: "unknown", expected: godbsdk.Mysql},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := sdkDatabaseType(test.dbType); actual != test.expected {
				t.Fatalf("sdkDatabaseType(%q) = %q, want %q", test.dbType, actual, test.expected)
			}
		})
	}
}
