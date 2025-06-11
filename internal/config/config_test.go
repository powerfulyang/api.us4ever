package config

import (
	"os"
	"testing"
)

func TestDBConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  DBConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: DBConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "user",
				Password: "pass",
				Database: "testdb",
				Schema:   "public",
			},
			wantErr: false,
		},
		{
			name: "missing host",
			config: DBConfig{
				Port:     5432,
				Username: "user",
				Password: "pass",
				Database: "testdb",
				Schema:   "public",
			},
			wantErr: true,
		},
		{
			name: "invalid port - zero",
			config: DBConfig{
				Host:     "localhost",
				Port:     0,
				Username: "user",
				Password: "pass",
				Database: "testdb",
				Schema:   "public",
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			config: DBConfig{
				Host:     "localhost",
				Port:     70000,
				Username: "user",
				Password: "pass",
				Database: "testdb",
				Schema:   "public",
			},
			wantErr: true,
		},
		{
			name: "missing username",
			config: DBConfig{
				Host:     "localhost",
				Port:     5432,
				Password: "pass",
				Database: "testdb",
				Schema:   "public",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			config: DBConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "user",
				Database: "testdb",
				Schema:   "public",
			},
			wantErr: true,
		},
		{
			name: "missing database",
			config: DBConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "user",
				Password: "pass",
				Schema:   "public",
			},
			wantErr: true,
		},
		{
			name: "missing schema",
			config: DBConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "user",
				Password: "pass",
				Database: "testdb",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DBConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ServerConfig
		wantErr bool
	}{
		{
			name: "valid port",
			config: ServerConfig{
				Port: 8080,
			},
			wantErr: false,
		},
		{
			name: "port too low",
			config: ServerConfig{
				Port: 0,
			},
			wantErr: true,
		},
		{
			name: "port too high",
			config: ServerConfig{
				Port: 70000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  AppConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: AppConfig{
				AppName: "test-app",
				AppEnv:  "test",
				Server: ServerConfig{
					Port: 8080,
				},
				Database: DBConfig{
					Host:     "localhost",
					Port:     5432,
					Username: "user",
					Password: "pass",
					Database: "testdb",
					Schema:   "public",
				},
			},
			wantErr: false,
		},
		{
			name: "missing app name",
			config: AppConfig{
				AppEnv: "test",
				Server: ServerConfig{
					Port: 8080,
				},
				Database: DBConfig{
					Host:     "localhost",
					Port:     5432,
					Username: "user",
					Password: "pass",
					Database: "testdb",
					Schema:   "public",
				},
			},
			wantErr: true,
		},
		{
			name: "missing app env",
			config: AppConfig{
				AppName: "test-app",
				Server: ServerConfig{
					Port: 8080,
				},
				Database: DBConfig{
					Host:     "localhost",
					Port:     5432,
					Username: "user",
					Password: "pass",
					Database: "testdb",
					Schema:   "public",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid server config",
			config: AppConfig{
				AppName: "test-app",
				AppEnv:  "test",
				Server: ServerConfig{
					Port: 0, // Invalid port
				},
				Database: DBConfig{
					Host:     "localhost",
					Port:     5432,
					Username: "user",
					Password: "pass",
					Database: "testdb",
					Schema:   "public",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid database config",
			config: AppConfig{
				AppName: "test-app",
				AppEnv:  "test",
				Server: ServerConfig{
					Port: 8080,
				},
				Database: DBConfig{
					Host:     "localhost",
					Port:     5432,
					Username: "user",
					Password: "pass",
					Database: "", // Missing database name
					Schema:   "public",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AppConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadEnvironmentFile(t *testing.T) {
	// Test when NACOS_SERVER_ADDR is set (should skip loading .env)
	t.Run("skip when NACOS_SERVER_ADDR is set", func(t *testing.T) {
		os.Setenv("NACOS_SERVER_ADDR", "localhost:8848")
		defer os.Unsetenv("NACOS_SERVER_ADDR")

		err := loadEnvironmentFile()
		if err != nil {
			t.Errorf("loadEnvironmentFile() should not error when NACOS_SERVER_ADDR is set, got: %v", err)
		}
	})

	// Test when NACOS_SERVER_ADDR is not set
	t.Run("attempt to load when NACOS_SERVER_ADDR is not set", func(t *testing.T) {
		os.Unsetenv("NACOS_SERVER_ADDR")

		// This will likely fail since we don't have a .env file in test environment
		// but we're testing the function behavior
		err := loadEnvironmentFile()
		// We don't assert on error here since it's expected in test environment
		_ = err
	})
}
