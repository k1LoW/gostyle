package config

import (
	"reflect"

	"golang.org/x/tools/go/analysis"
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
	return c, nil
}

func init() {
	Loader.Flags.StringVar(&configPath, "config", "", "config file path")
}

type Config struct{}
