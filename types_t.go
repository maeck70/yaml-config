package yamlconfig

// Base config structure
type Config_t struct {
	Metadata metadata_t  `yaml:"metadata"`
	Data     interface{} `yaml:"data"`
}

// Metadata for the config file
type metadata_t struct {
	Schema string `yaml:"schema"`
}

type sf_t map[string]SchemaField_t

// Validation Schema
type ConfigValidator_t struct {
	Schema sf_t `yaml:"schema"`
}

// Schema Field
type SchemaField_t struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     any    `yaml:"default"`
	Options     []any  `yaml:"options"`
	OptionType  string `yaml:"optiontype"`
	Min         int64  `yaml:"min"`
	Max         int64  `yaml:"max"`
	Attributes  sf_t   `yaml:"attributes"`
}
