package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/importer"
	"go/types"
	"io"
	"os"
	"slices"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/unitchecker"
)

const (
	DefaultSmallScopeMax       = 7
	DefaultSmallVarnameMax     = -1
	DefaultMediumScopeMax      = 15
	DefaultMediumVarnameMax    = -1
	DefaultLargeScopeMax       = 25
	DefaultLargeVarnameMax     = -1
	DefaultVeryLargeVarnameMax = -1
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
	Getters     Getters     `yaml:"getters"`
	Ifacenames  Ifacenames  `yaml:"ifacenames"`
	Mixedcaps   Mixedcaps   `yaml:"mixedcaps"`
	Nilslices   Nilslices   `yaml:"nilslices"`
	Pkgnames    Pkgnames    `yaml:"pkgnames"`
	Recvnames   Recvnames   `yaml:"recvnames"`
	Recvtype    Recvtype    `yaml:"recvtype"`
	Repetition  Repetition  `yaml:"repetition"`
	Typealiases Typealiases `yaml:"typealiases"`
	Underscores Underscores `yaml:"underscores"`
	Useq        Useq        `yaml:"useq"`
	Varnames    Varnames    `yaml:"varnames"`
}

type Getters struct {
	Exclude          []string `yaml:"exclude"`
	IncludeGenerated bool     `yaml:"include-generated"`
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

func NewTypesConfig(pass *analysis.Pass) (types.Config, error) { //nostyle:repetition
	args := flag.Args()
	if len(args) == 0 {
		return types.Config{Importer: importer.Default()}, nil
	}
	filename := args[0]
	cfg, err := readConfig(filename)
	if err != nil {
		return types.Config{}, err
	}
	return types.Config{Importer: importer.ForCompiler(pass.Fset, cfg.Compiler, func(path string) (io.ReadCloser, error) {
		file, ok := cfg.PackageFile[path]
		if !ok {
			if cfg.Compiler == "gccgo" && cfg.Standard[path] {
				return nil, nil // fall back to default gccgo lookup
			}
			return nil, fmt.Errorf("no package file for %q", path)
		}
		return os.Open(file)
	})}, nil
}

func readConfig(filename string) (*unitchecker.Config, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := new(unitchecker.Config)
	if err := json.Unmarshal(b, cfg); err != nil {
		return nil, fmt.Errorf("cannot decode JSON config file %s: %w", filename, err)
	}
	if len(cfg.GoFiles) == 0 {
		return nil, fmt.Errorf("package has no files: %s", cfg.ImportPath)
	}
	return cfg, nil
}
