package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"gopkg.in/yaml.v3"
)

const (
	name = "gostyle"
	doc  = "config"
)

var configPath string

var Loader = &analysis.Analyzer{
	Name:       name,
	Doc:        doc,
	Run:        run,
	Requires:   []*analysis.Analyzer{},
	ResultType: reflect.TypeOf((*Config)(nil)),
}

func run(pass *analysis.Pass) (any, error) {
	c := &Config{}
	if configPath == "" {
		return c, nil
	}
	if !strings.HasPrefix(configPath, "/") {
		c.err = fmt.Errorf("config file path must be specified as a full path")
		return c, nil
	}
	f, err := os.Open(configPath)
	if err != nil {
		c.err = fmt.Errorf("failed to open config file: %w", err)
		return c, nil
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		c.err = fmt.Errorf("failed to decode config file: %w", err)
		return c, nil
	}
	c.loaded = true
	return c, nil
}

func init() {
	Loader.Flags.StringVar(&configPath, "config", "", "config file path. the config file path must be specified as a full path. if a config file is specified, each analyzer options are ignored")
}
