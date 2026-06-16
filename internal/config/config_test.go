package config

import "testing"

// fakeEnv returns a getenv function backed by a map.
func fakeEnv(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoadDefaults(t *testing.T) {
	c, err := load(fakeEnv(nil))
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}
	if c.Env != EnvDevelopment {
		t.Errorf("Env = %q, want %q", c.Env, EnvDevelopment)
	}
	if c.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, want :8080", c.HTTPAddr)
	}
	if c.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want info", c.LogLevel)
	}
	if c.LogFormat != "json" {
		t.Errorf("LogFormat = %q, want json", c.LogFormat)
	}
	if c.DatabaseURL != "" {
		t.Errorf("DatabaseURL = %q, want empty", c.DatabaseURL)
	}
}

func TestLoadOverrides(t *testing.T) {
	c, err := load(fakeEnv(map[string]string{
		"FARAID_ENV":          EnvProduction,
		"FARAID_HTTP_ADDR":    "localhost:9000",
		"FARAID_LOG_LEVEL":    "debug",
		"FARAID_LOG_FORMAT":   "text",
		"FARAID_DATABASE_URL": "postgres://localhost/faraid",
	}))
	if err != nil {
		t.Fatalf("load overrides: %v", err)
	}
	if c.Env != EnvProduction || c.HTTPAddr != "localhost:9000" ||
		c.LogLevel != "debug" || c.LogFormat != "text" ||
		c.DatabaseURL != "postgres://localhost/faraid" {
		t.Errorf("overrides not applied: %+v", c)
	}
}

func TestValidationErrors(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
	}{
		{"bad env", map[string]string{"FARAID_ENV": "staging"}},
		{"bad level", map[string]string{"FARAID_LOG_LEVEL": "verbose"}},
		{"bad format", map[string]string{"FARAID_LOG_FORMAT": "xml"}},
		{"bad addr", map[string]string{"FARAID_HTTP_ADDR": "8080"}},
		{"prod without db", map[string]string{"FARAID_ENV": EnvProduction}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := load(fakeEnv(c.env)); err == nil {
				t.Errorf("expected error for %s", c.name)
			}
		})
	}
}

func TestProductionWithDatabaseURL(t *testing.T) {
	_, err := load(fakeEnv(map[string]string{
		"FARAID_ENV":          EnvProduction,
		"FARAID_DATABASE_URL": "postgres://localhost/faraid",
	}))
	if err != nil {
		t.Errorf("production with database url should be valid: %v", err)
	}
}

func TestLoad(t *testing.T) {
	t.Setenv("FARAID_ENV", EnvTest)
	t.Setenv("FARAID_LOG_LEVEL", "info")
	t.Setenv("FARAID_LOG_FORMAT", "json")
	t.Setenv("FARAID_HTTP_ADDR", ":8080")
	t.Setenv("FARAID_DATABASE_URL", "")
	c, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.Env != EnvTest {
		t.Errorf("Env = %q, want %q", c.Env, EnvTest)
	}
}

func TestOrDefault(t *testing.T) {
	if got := orDefault("", "def"); got != "def" {
		t.Errorf("orDefault empty = %q, want def", got)
	}
	if got := orDefault("set", "def"); got != "set" {
		t.Errorf("orDefault set = %q, want set", got)
	}
}
