package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/goccy/go-yaml"
	"golang.org/x/tools/go/analysis"
)

const (
	name = "gostyle"
	doc  = "config"
)

var configPath string
var defaultFileNames = []string{
	".gostyle.yml",
	".gostyle.yaml",
}

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
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
	L:
		for {
			if _, err := os.Stat(filepath.Join(wd, ".git", "config")); err == nil {
				for _, n := range defaultFileNames {
					p := filepath.Join(wd, n)
					if _, err := os.Stat(p); err == nil {
						configPath = p
						break L
					}
				}
			}
			if wd == filepath.Dir(wd) {
				break
			}
			wd = filepath.Dir(wd)
		}
		if configPath == "" {
			return c, nil
		}
	}
	if !filepath.IsAbs(configPath) {
		c.err = fmt.Errorf("config file path must be absolute path: %s", configPath)
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
	c.ConfigDir = filepath.Dir(configPath)

	// Set default value
	if c.AnalyzersSettings.Recvnames.Max == 0 {
		c.AnalyzersSettings.Recvnames.Max = DefaultReceiverNameMax
	}
	if c.AnalyzersSettings.Varnames.SmallScopeMax == 0 {
		c.AnalyzersSettings.Varnames.SmallScopeMax = DefaultSmallScopeMax
	}
	if c.AnalyzersSettings.Varnames.SmallVarnameMax == 0 {
		c.AnalyzersSettings.Varnames.SmallVarnameMax = DefaultSmallVarnameMax
	}
	if c.AnalyzersSettings.Varnames.MediumScopeMax == 0 {
		c.AnalyzersSettings.Varnames.MediumScopeMax = DefaultMediumScopeMax
	}
	if c.AnalyzersSettings.Varnames.MediumVarnameMax == 0 {
		c.AnalyzersSettings.Varnames.MediumVarnameMax = DefaultMediumVarnameMax
	}
	if c.AnalyzersSettings.Varnames.LargeScopeMax == 0 {
		c.AnalyzersSettings.Varnames.LargeScopeMax = DefaultLargeScopeMax
	}
	if c.AnalyzersSettings.Varnames.LargeVarnameMax == 0 {
		c.AnalyzersSettings.Varnames.LargeVarnameMax = DefaultLargeVarnameMax
	}
	if c.AnalyzersSettings.Varnames.VeryLargeVarnameMax == 0 {
		c.AnalyzersSettings.Varnames.VeryLargeVarnameMax = DefaultVeryLargeVarnameMax
	}

	c.loaded = true
	return c, nil
}

func SetPath(p string) {
	configPath = p
}

func init() {
	Loader.Flags.StringVar(&configPath, "config", "", "config file path. the config file path must be specified as a full path. if a config file is specified, each analyzer options are ignored")
}
