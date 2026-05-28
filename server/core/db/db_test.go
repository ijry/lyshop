package db

import "testing"

func TestResolveDatabaseTarget(t *testing.T) {
	tests := []struct {
		name       string
		rawDSN     string
		wantDriver string
		wantDSN    string
	}{
		{
			name:       "empty dsn defaults to sqlite",
			rawDSN:     "",
			wantDriver: "sqlite",
			wantDSN:    defaultSQLiteDSN,
		},
		{
			name:       "mysql dsn stays mysql",
			rawDSN:     "root:password@tcp(mysql:3306)/lyshop?charset=utf8mb4&parseTime=True&loc=Local",
			wantDriver: "mysql",
			wantDSN:    "root:password@tcp(mysql:3306)/lyshop?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name:       "sqlite file path",
			rawDSN:     "data/lyshop.db",
			wantDriver: "sqlite",
			wantDSN:    "data/lyshop.db",
		},
		{
			name:       "sqlite scheme",
			rawDSN:     "sqlite://data/lyshop.sqlite",
			wantDriver: "sqlite",
			wantDSN:    "data/lyshop.sqlite",
		},
		{
			name:       "sqlite memory mode",
			rawDSN:     "file:test?mode=memory&cache=shared",
			wantDriver: "sqlite",
			wantDSN:    "file:test?mode=memory&cache=shared",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDriver, gotDSN := resolveDatabaseTarget(tt.rawDSN)
			if gotDriver != tt.wantDriver {
				t.Fatalf("driver mismatch: got %q, want %q", gotDriver, tt.wantDriver)
			}
			if gotDSN != tt.wantDSN {
				t.Fatalf("dsn mismatch: got %q, want %q", gotDSN, tt.wantDSN)
			}
		})
	}
}
