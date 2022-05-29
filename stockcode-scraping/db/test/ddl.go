package test

type AttributeDefinitions struct {
	AttributeName []string
	AttributeType string
}

type Properties struct {
	TableName            string `yaml:"TableName"`
	AttributeDefinitions `yaml:",inline"`
}

type StCode struct {
	Type       interface{} `yaml:"Type"`
	Properties `yaml:",inline"`
}

type Resources struct {
	StCode `yaml:",inline"`
}

type Tbl struct {
	AWSTemplateFormatVersion string `yaml:"AWSTemplateFormatVersion"`
	Resources                `yaml:",inline"`
}
