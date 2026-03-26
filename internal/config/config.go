package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/adrg/xdg"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

const envPrefix = "SERVERBEACON_"

var (
	once sync.Once
	_cfg *Config
	K    = koanf.New(".")
)

// Main config object
// NOTE: When adding fields or structs, don't forget
// to update the createDefaultConfigWithEnvVars() mapping.
type Config struct {
	Debug       bool      `koanf:"debug"`
	APiSettings APIConfig `koanf:"api"`
	DB          DBConfig  `koanf:"db"`
}

// Settings for REST API server
type APIConfig struct {
	// Default: http
	Proto string `koanf:"proto"`
	// Default: 0.0.0.0
	Host string `koanf:"host"`
	// Default: 18080
	Port int64 `koanf:"port"`
}

// Settings for app database
type DBConfig struct {
	Path string `koanf:"path" path:"expand"`
}

// Determines the Koanf parser to use based on filetype
func parserForFile(path string) (koanf.Parser, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		return yaml.Parser(), nil
	case ".json":
		return json.Parser(), nil
	case ".toml":
		return toml.Parser(), nil
	case ".env":
		return dotenv.Parser(), nil
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}
}

// createDefaultConfigWithEnvVars creates a default config map with values from env vars if available
// func createDefaultConfigWithEnvVars() map[string]interface{} {
// 	config := map[string]interface{}{
// 		"debug": false,
// 		"api": map[string]interface{}{
// 			"proto": getEnvOrDefault(envPrefix+"PROTO", "http"),
// 			"host":  getEnvOrDefault(envPrefix+"HOST", "0.0.0.0"),
// 			"port":  getEnvOrDefault(envPrefix+"PORT", "18080"),
// 		},
// 		"db": map[string]interface{}{
// 			"path": "serverbeacon.sqlite3",
// 		},
// 	}

// 	return config
// }

// Return the default config file path (~/.local/share/serverbeacon/config.yml)
func GetDefaultConfigPath() string {
	return filepath.Join(xdg.DataHome, "serverbeacon", "config.yml")
}

// Return environment prefix
func GetEnvPrefix() string {
	return envPrefix
}

// FindConfigFile checks for a .local variant of the config file first,
// falling back to the original if .local doesn't exist.
// If configFile is empty, returns the default XDG config path.
func FindConfigFile(configFile string) string {
	// Search for config file if none is provided
	if configFile != "" {
		ext := filepath.Ext(configFile)
		base := strings.TrimSuffix(configFile, ext)
		localFile := base + ".local" + ext

		if _, err := os.Stat(localFile); err == nil {
			return localFile
		}
		if _, err := os.Stat(configFile); err == nil {
			return configFile
		}
		return ""
	}

	// Return config (or local config) if file is found
	candidates := []string{"config.local.yml", "config.yml"}
	for _, f := range candidates {
		if _, err := os.Stat(f); err == nil {
			return f
		}
	}

	return ""
}

// LoadConfig loads configuration from a file, environment variables, and/org CLI args.
// Returns the parsed config struct
func LoadConfig(flagSet *pflag.FlagSet, configFile string) (*Config, error) {
	// Initialize config object
	K = koanf.New(".")

	// Set config defaults
	defaults := map[string]interface{}{
		"debug": false,
		"api": map[string]interface{}{
			"proto": "http",
			"host":  "0.0.0.0",
			"port":  18080,
		},
	}

	// Set defaults in Koanf config object
	if err := K.Load(confmap.Provider(defaults, "."), nil); err != nil {
		return nil, fmt.Errorf("error loading defaults: %w", err)
	}

	// Attempt to load from config file
	configFile = FindConfigFile(configFile)
	if configFile != "" {
		if _, err := os.Stat(configFile); err == nil {
			// Detect parser from file extension
			parser, err := parserForFile(configFile)
			if err != nil {
				return nil, fmt.Errorf("unsupported config file format: %w", err)
			}

			// Load config from file
			if err := K.Load(file.Provider(configFile), parser); err != nil {
				return nil, fmt.Errorf("error loading config file: %w", err)
			}
		}
	}

	// Load env vars
	if err := K.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil); err != nil {
		return nil, fmt.Errorf("error loading env vars: %w", err)
	}

	// Load CLI args
	if flagSet != nil {
		if err := K.Load(posflag.Provider(flagSet, ".", K), nil); err != nil {
			return nil, fmt.Errorf("error loading flags: %w", err)
		}
	}

	// Unmarshal config into Config struct
	var cfg Config
	if err := K.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

// ensureConfigFile creates the config file if it doesn't exist
// func ensureConfigFile(configFile string) error {
// 	// Create config directory if it doesn't exist
// 	configDir := filepath.Dir(configFile)
// 	if err := os.MkdirAll(configDir, 0755); err != nil {
// 		return fmt.Errorf("failed to create config directory: %w", err)
// 	}

// 	// Generate default config with env vars and prompted token
// 	configData := createDefaultConfigWithEnvVars()

// 	// Marshal to YAML
// 	data, err := yamlv3.Marshal(configData)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal default config: %w", err)
// 	}

// 	// Write to file
// 	if err := os.WriteFile(configFile, data, 0644); err != nil {
// 		return fmt.Errorf("failed to write config file: %w", err)
// 	}

// 	fmt.Printf("\nConfig file created: %s\n", configFile)
// 	return nil
// }

// expandPaths walks the config struct and expands ~ in any field tagged with path:"expand"
// func (c *Config) expandPaths() {
// 	expandStructPaths(reflect.ValueOf(c).Elem())
// }

// expandStructPaths recursively walks a struct and expands paths in tagged fields
// func expandStructPaths(v reflect.Value) {
// 	if v.Kind() != reflect.Struct {
// 		return
// 	}

// 	t := v.Type()
// 	for i := 0; i < v.NumField(); i++ {
// 		field := v.Field(i)
// 		fieldType := t.Field(i)

// 		// Check if field has path:"expand" tag
// 		if tag := fieldType.Tag.Get("path"); tag == "expand" {
// 			if field.Kind() == reflect.String && field.CanSet() {
// 				field.SetString(expandPath(field.String()))
// 			}
// 		}

// 		// Recursively handle nested structs
// 		if field.Kind() == reflect.Struct {
// 			expandStructPaths(field)
// 		}
// 	}
// }

// expandPath returns the expanded path, handling ~ for home directory and converting to absolute path
// func expandPath(path string) string {
// 	// Handle ~ expansion
// 	if strings.HasPrefix(path, "~/") {
// 		home, err := os.UserHomeDir()
// 		if err == nil {
// 			path = filepath.Join(home, path[2:])
// 		}
// 	}

// 	// Convert to absolute path
// 	absPath, err := filepath.Abs(path)
// 	if err == nil {
// 		return absPath
// 	}

// 	return path // Return original if expansion fails
// }

// getEnvOrDefault gets an environment variable or returns the default value
// func getEnvOrDefault(key, defaultValue string) string {
// 	if value := os.Getenv(key); value != "" {
// 		return value
// 	}

// 	return defaultValue
// }

// Init loads and caches the config (call from entrypoints like cmd/root.go)
func Init(flagSet *pflag.FlagSet, configFile string) error {
	var initErr error
	once.Do(func() {
		K = koanf.New(".")

		// Check locations for existin config file
		configFile = FindConfigFile(configFile)

		c, err := LoadConfig(flagSet, configFile)
		if err != nil {
			initErr = fmt.Errorf("config.Init failed: %w", err)
			return
		}
		_cfg = c
	})

	// Debug config
	// fmt.Printf("Config: %+v\n", c)

	return initErr
}

// GetConfig returns the initialized Config
func GetConfig() *Config {
	if _cfg == nil {
		panic("config.Init() must be called before GetConfig()")
	}
	return _cfg
}

// GetKoanf returns the raw koanf instance
func GetKoanf() *koanf.Koanf {
	return K
}
