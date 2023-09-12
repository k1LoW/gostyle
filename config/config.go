package config

import (
	"fmt"
	"slices"

	"golang.org/x/tools/go/analysis"
)

type Config struct {
	Analyzers        Analyzers        `yaml:"analyzers"`
	AnalyzerSettings AnalyzerSettings `yaml:"analyzer-settings"`
	loaded           bool
	err              error
}

type Analyzers struct {
	Disable []string `yaml:"disable"`
}

type AnalyzerSettings struct {
	Ifacenames  Ifacenames  `yaml:"ifacenames"`
	Mixedcaps   Mixedcaps   `yaml:"mixedcaps"`
	Pkgnames    Pkgnames    `yaml:"pkgnames"`
	Recvnames   Recvnames   `yaml:"recvnames"`
	Underscores Underscores `yaml:"underscores"`
}

type Ifacenames struct {
	All              bool `yaml:"all"`
	IncludeGenerated bool `yaml:"include-generated"`
}

type Mixedcaps struct {
	Exclude          []string `yaml:"exclude"`
	IncludeGenerated bool     `yaml:"include-generated"`
}

type Pkgnames struct {
	IncludeGenerated bool `yaml:"include-generated"`
}

type Recvnames struct {
	IncludeGenerated bool `yaml:"include-generated"`
}

type Underscores struct {
	Exclude          []string `yaml:"exclude"`
	IncludeGenerated bool     `yaml:"include-generated"`
}

func (c *Config) IsDisabled(name string) bool {
	return slices.Contains(c.Analyzers.Disable, name)
}

func Load(pass *analysis.Pass) (*Config, error) {
	c, ok := pass.ResultOf[Loader].(*Config)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from config.Loader: %T", pass.ResultOf[Loader])
	}
	if c.err != nil {
		return nil, c.err
	}
	if !c.loaded {
		return nil, nil
	}
	return c, nil
}
