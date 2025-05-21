package models

type LanguageSettings struct {
	Code     string `json:"code" bson:"code" mapstructure:"code"`
	Language string `json:"language" bson:"language" mapstructure:"language"`
}

// Settings represents the configuration settings structure
type Settings struct {
	Id       string           `bson:"id,omitempty" json:"id" mapstructure:"id"`
	Language LanguageSettings `bson:"language" json:"language" mapstructure:"language"`
}
