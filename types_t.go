package yamlconfig

// Base config structure
type Config_t struct {
	Metadata Metadata_t  `yaml:"metadata"`
	Schema   interface{} `yaml:"schema"`
}

// Metadata for the config file
type Metadata_t struct {
	Id            string `yaml:"id"`
	SchemaVersion string `yaml:"schema_version"`
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
