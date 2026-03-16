package config

// Defaults holds default values for new tickets.
type Defaults struct {
	Priority string   `yaml:"priority"`
	Ready    bool     `yaml:"ready"`
	Tags     []string `yaml:"tags"`
}

// Config holds swimlane configuration.
// ConfigDir is the directory containing the config file (or cwd if no file);
// tickets glob and default_path are resolved relative to ConfigDir.
type Config struct {
	Tickets     string   `yaml:"tickets"`
	DefaultPath string   `yaml:"default_path"`
	Default     Defaults `yaml:"default"`
	ConfigDir   string   `yaml:"-"` // set by loader; not in YAML
}

// ValidPriorities is the set of allowed priority values.
var ValidPriorities = []string{"p0", "p1", "p2", "p3", "p4"}

// ValidStatuses is the set of allowed status values.
var ValidStatuses = []string{"todo", "in-progress", "done"}
