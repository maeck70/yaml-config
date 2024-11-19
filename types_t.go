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

type sfattribute_t map[string]SchemaField_t

type sfitem_t []SchemaField_t

// Base Validation Schema
type ConfigValidator_t struct {
	Metadata Metadata_t    `yaml:"metadata"`
	Schema   sfattribute_t `yaml:"schema"`
}

// Schema Field
type SchemaField_t struct {
	Type        string        `yaml:"type"`
	Description string        `yaml:"description"`
	Required    bool          `yaml:"required"`
	Default     any           `yaml:"default"`
	Options     []any         `yaml:"options"`
	OptionType  string        `yaml:"optiontype"`
	Min         int64         `yaml:"min"`
	Max         int64         `yaml:"max"`
	Attributes  sfattribute_t `yaml:"attributes"`
	Items       sfitem_t      `yaml:"items"`
}
