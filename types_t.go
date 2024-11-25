package yamlconfig

// Metadata for the config and schema file
type Metadata_t struct {
	Id            string `yaml:"id"`
	SchemaVersion string `yaml:"schema_version"`
}

// Base Config structure
type Config_t struct {
	Metadata Metadata_t  `yaml:"metadata"`
	Data     interface{} `yaml:"data"`
}

// Base Validation Schema
type ConfigValidator_t struct {
	Metadata Metadata_t `yaml:"metadata"`
	Schema   Schema_t   `yaml:"schema"`
}

// Schema is the map of fields on that level
type Schema_t map[string]SchemaField_t

type GroupField_t struct {
	Type        string                 `yaml:"type"`
	Description string                 `yaml:"description"`
	Attributes  map[string]interface{} `yaml:"attributes"`
}

// SchemaField contains the attributes for the fields
type SchemaField_t struct {
	Type        string                   `yaml:"type"`
	Description string                   `yaml:"description"`
	Required    bool                     `yaml:"required"`
	Default     any                      `yaml:"default"`
	Options     []any                    `yaml:"options"`
	OptionType  string                   `yaml:"optiontype"`
	Min         int                      `yaml:"min"`
	Max         int                      `yaml:"max"`
	Attributes  map[string]SchemaField_t `yaml:"attributes"`
	List        map[string]SchemaField_t `yaml:"list"`
	Valid       []string                 `yaml:"valid"`
	Group       GroupField_t             `yaml:"group"`
}
