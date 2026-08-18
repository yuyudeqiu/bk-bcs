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

func TestNormalizeSDKDatabaseConfig(t *testing.T) {
	tests := []struct {
		name           string
		dbConfig       godbsdk.Database
		wantSSLEnable  bool
		wantSSLMode    string
		wantTimeout    string
		wantParamsSize int
	}{
		{
			name:           "gaussdb ssl disabled",
			dbConfig:       godbsdk.Database{Typex: godbsdk.Gaussdb},
			wantSSLEnable:  true,
			wantSSLMode:    "disable",
			wantTimeout:    "10",
			wantParamsSize: 1,
		},
		{
			name: "gaussdb ssl enabled",
			dbConfig: godbsdk.Database{Typex: godbsdk.Gaussdb, Ssl: godbsdk.TLS{
				Enable: true,
				Mode:   "verify-full",
			}},
			wantSSLEnable:  true,
			wantSSLMode:    "verify-full",
			wantTimeout:    "10",
			wantParamsSize: 1,
		},
		{
			name: "gaussdb custom timeout",
			dbConfig: godbsdk.Database{
				Typex:  godbsdk.Gaussdb,
				Params: map[string]string{"connect_timeout": "30"},
			},
			wantSSLEnable:  true,
			wantSSLMode:    "disable",
			wantTimeout:    "30",
			wantParamsSize: 1,
		},
		{
			name:           "mysql unchanged",
			dbConfig:       godbsdk.Database{Typex: godbsdk.Mysql},
			wantSSLEnable:  false,
			wantSSLMode:    "",
			wantTimeout:    "",
			wantParamsSize: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			normalizeSDKDatabaseConfig(&test.dbConfig)
			if test.dbConfig.Ssl.Enable != test.wantSSLEnable {
				t.Fatalf("unexpected ssl enable: got %v, want %v", test.dbConfig.Ssl.Enable, test.wantSSLEnable)
			}
			if test.dbConfig.Ssl.Mode != test.wantSSLMode {
				t.Fatalf("unexpected ssl mode: got %q, want %q", test.dbConfig.Ssl.Mode, test.wantSSLMode)
			}
			if test.dbConfig.Params["connect_timeout"] != test.wantTimeout {
				t.Fatalf("unexpected connect timeout: got %q, want %q",
					test.dbConfig.Params["connect_timeout"], test.wantTimeout)
			}
			if len(test.dbConfig.Params) != test.wantParamsSize {
				t.Fatalf("unexpected params size: got %d, want %d", len(test.dbConfig.Params), test.wantParamsSize)
			}
		})
	}
}
