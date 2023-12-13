package config

import (
	"slices"

	"golang.org/x/tools/go/analysis"
)

const (
	DefaultSmallScopeMax       = 7
	DefaultSmallVarnameMax     = -1
	DefaultMediumScopeMax      = 15
	DefaultMediumVarnameMax    = -1
	DefaultLargeScopeMax       = 25
	DefaultLargeVarnameMax     = -1
	DefaultVeryLargeVarnameMax = -1
	DefaultReceiverNameMax     = 2
)

type Config struct {
	Analyzers         Analyzers         `yaml:"analyzers"`
	AnalyzersSettings AnalyzersSettings `yaml:"analyzers-settings"`
	loaded            bool
	err               error
}

type Analyzers struct {
	Disable []string `yaml:"disable"`
}

type AnalyzersSettings struct {
	Contexts     Contexts     `yaml:"contexts"`
	Dontpanic    Dontpanic    `yaml:"dontpanic"`
	Errorstrings Errorstrings `yaml:"errorstrings"`
	Funcfmt      Funcfmt      `yaml:"funcfmt"`
	Getters      Getters      `yaml:"getters"`
	Handlerrors  Handlerrors  `yaml:"handlerrors"`
	Ifacenames   Ifacenames   `yaml:"ifacenames"`
	Mixedcaps    Mixedcaps    `yaml:"mixedcaps"`
	Nilslices    Nilslices    `yaml:"nilslices"`
	Pkgnames     Pkgnames     `yaml:"pkgnames"`
	Recvnames    Recvnames    `yaml:"recvnames"`
	Recvtype     Recvtype     `yaml:"recvtype"`
	Repetition   Repetition   `yaml:"repetition"`
	Typealiases  Typealiases  `yaml:"typealiases"`
	Underscores  Underscores  `yaml:"underscores"`
	Useany       Useany       `yaml:"useany"`
	Useq         Useq         `yaml:"useq"`
	Varnames     Varnames     `yaml:"varnames"`
}

type Contexts struct {
	IncludeGenerated bool `yaml:"include-generated"`
	ExcludeTest      bool `yaml:"exclude-test"`
}

type Dontpanic struct {
	IncludeGenerated bool `yaml:"include-generated"`
	ExcludeTest      bool `yaml:"exclude-test"`
}

type Errorstrings struct {
	IncludeGenerated bool `yaml:"include-generated"`
	ExcludeTest      bool `yaml:"exclude-test"`
}

type Funcfmt struct {
	IncludeGenerated bool `yaml:"include-generated"`
}

type Getters struct {
	Exclude          []string `yaml:"exclude"`
	IncludeGenerated bool     `yaml:"include-generated"`
}

type Handlerrors struct {
	IncludeGenerated bool `yaml:"include-generated"`
	ExcludeTest      bool `yaml:"exclude-test"`
}

type Ifacenames struct {
	All              bool `yaml:"all"`
	IncludeGenerated bool `yaml:"include-generated"`
}

type Mixedcaps struct {
	Exclude          []string `yaml:"exclude"`
	IncludeGenerated bool     `yaml:"include-generated"`
}

type Nilslices struct {
	IncludeGenerated bool `yaml:"include-generated"`
}

type Pkgnames struct {
	IncludeGenerated bool `yaml:"include-generated"`
}

type Recvnames struct {
	IncludeGenerated bool `yaml:"include-generated"`
	Max              int  `yaml:"max"`
}

type Recvtype struct {
	IncludeGenerated bool `yaml:"include-generated"`
}

type Repetition struct {
	Exclude          []string `yaml:"exclude"`
	IncludeGenerated bool     `yaml:"include-generated"`
}

type Typealiases struct {
	Exclude          []string `yaml:"exclude"`
	IncludeGenerated bool     `yaml:"include-generated"`
}

type Underscores struct {
	Exclude          []string `yaml:"exclude"`
	IncludeGenerated bool     `yaml:"include-generated"`
}

type Useany struct {
	IncludeGenerated bool `yaml:"include-generated"`
}

type Useq struct {
	IncludeGenerated bool `yaml:"include-generated"`
}

type Varnames struct {
	Exclude             []string `yaml:"exclude"`
	IncludeGenerated    bool     `yaml:"include-generated"`
	SmallScopeMax       int      `yaml:"small-scope-max"`
	SmallVarnameMax     int      `yaml:"small-varname-max"`
	MediumScopeMax      int      `yaml:"medium-scope-max"`
	MediumVarnameMax    int      `yaml:"medium-varname-max"`
	LargeScopeMax       int      `yaml:"large-scope-max"`
	LargeVarnameMax     int      `yaml:"large-varname-max"`
	VeryLargeVarnameMax int      `yaml:"very-large-varname-max"`
}

func (c *Config) IsDisabled(name string) bool {
	return slices.Contains(c.Analyzers.Disable, name)
}

func Load(pass *analysis.Pass) (*Config, error) {
	c, ok := pass.ResultOf[Loader].(*Config)
	if !ok {
		return nil, nil
	}
	if c.err != nil {
		return nil, c.err
	}
	if !c.loaded {
		return nil, nil
	}
	return c, nil
}
