package types

import (
	ktypes "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
)

type ResourceIdentifier struct {
	Kind string `json:"kind" yaml:"kind"`
	Name string `json:"name" yaml:"name"`
}

type Substitutions struct {
	Substitutions []*Substitution `json:"substitutions" yaml:"substitutions"`
}

type ConfigFile struct {
	Id ResourceIdentifier `json:"configFile" yaml:"configFile"`
}

type OutputFile struct {
	Id ResourceIdentifier `json:"outputFile" yaml:"outputFile"`
}

type Substitution struct {
	Id           string      `json:"id" yaml:"id"`
	DefaultValue int         `json:"default" yaml:"default"`
	Replacement  resid.ResId `json:"replacement" yaml:"replacement"`
}

type SourceMeta struct {
	Substitution *Substitution
	Source       *ktypes.SourceSelector
}

type Source struct {
	ResourceIdentifier `json:",inline" yaml:",inline"`
	FieldPath          string      `json:"fieldPath" yaml:"fieldPath"`
	Replacement        interface{} `json:"replacement" yaml:"replacement"`
}

type SubReplacement struct {
	FieldPath   string      `json:"fieldPath" yaml:"fieldPath"`
	Replacement interface{} `json:"replacement" yaml:"replacement"`
}

type YamlValue struct {
	Values interface{} `json:"values" yaml:"values"`
}

type Replacement struct {
	Source  Source   `json:"source" yaml:"source"`
	Targets []Target `json:"targets" yaml:"targets"`
}

type Tagger interface {
	SetTag(newTag string)
	GetTag() []string
	RemoveTagEntry(tag string)
	IsList() bool
}

type Target struct {
	TSelect    ResourceIdentifier `json:"select" yaml:"select"`
	FieldPaths []string           `json:"fieldPaths" yaml:"fieldPaths"`
	Options    interface{}        `json:"options" yaml:"options"`
}

type MetaData struct {
	Name string `json:"name" yaml:"name"`
}

type ApplyReplacementsFnConfig struct {
	ApiVersion   string        `json:"apiVersion" yaml:"apiVersion"`
	Kind         string        `json:"kind" yaml:"kind"`
	Metadata     MetaData      `json:"metadata" yaml:"metadata"`
	Replacements []Replacement `json:"replacements,omitempty" yaml:"replacements,omitempty"`
}

// type SubConfig struct {
// 	// Kind string `json:"kind" yaml:"kind"`
// 	Data ktypes.SourceSelector `json:"data" yaml:"data"`
// }

// type SubConfigData struct {

// 	// Replacement interface{} `json:"replacement" yaml:"replacement"`
// 	FieldPath string               `json:"fieldPath" yaml:"fieldPath"`
// 	Options   *ktypes.FieldOptions `json:"options,omitempty" yaml:"options,omitempty"`
// }
