package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DefaultTicketsGlob is used when config has no tickets pattern.
const DefaultTicketsGlob = "tickets/**/*.md"

// DefaultPath is used when config has no default_path.
const DefaultPath = "tickets"

// Load reads configuration from the first existing file in lookup order,
// or from configPath if non-empty (overrides lookup).
func Load(configPath string) (*Config, error) {
	if configPath != "" {
		return loadFile(configPath)
	}
	// Lookup order: .swimlane.yaml, swimlane.yaml, ~/.config/swimlane/config.yaml
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}
	candidates := []string{
		filepath.Join(wd, ".swimlane.yaml"),
		filepath.Join(wd, "swimlane.yaml"),
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".config", "swimlane", "config.yaml"))
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return loadFile(p)
		}
	}
	// No config file: return defaults; paths relative to cwd
	return &Config{
		Tickets:     DefaultTicketsGlob,
		DefaultPath: DefaultPath,
		ConfigDir:   wd,
		Default: Defaults{
			Priority: "p2",
			Ready:    true,
			Tags:     nil,
		},
	}, nil
}

func loadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	dir := filepath.Dir(path)
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve config directory: %w", path, err)
	}
	c.ConfigDir = absDir
	if c.Tickets == "" {
		c.Tickets = DefaultTicketsGlob
	}
	if c.DefaultPath == "" {
		c.DefaultPath = DefaultPath
	}
	if err := Validate(&c); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &c, nil
}

// Validate checks required fields and allowed enums.
func Validate(c *Config) error {
	if c.Default.Priority != "" {
		ok := false
		for _, p := range ValidPriorities {
			if c.Default.Priority == p {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("default.priority must be one of p0, p1, p2, p3, p4; got %q", c.Default.Priority)
		}
	}
	return nil
}
