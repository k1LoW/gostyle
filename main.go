package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/k1LoW/gostyle/analyzer/code_review_comments/dontpanic"
	"github.com/k1LoW/gostyle/analyzer/code_review_comments/errorstrings"
	"github.com/k1LoW/gostyle/analyzer/code_review_comments/handlerrors"
	"github.com/k1LoW/gostyle/analyzer/decisions/getters"
	"github.com/k1LoW/gostyle/analyzer/decisions/nilslices"
	"github.com/k1LoW/gostyle/analyzer/decisions/pkgnames"
	"github.com/k1LoW/gostyle/analyzer/decisions/recvnames"
	"github.com/k1LoW/gostyle/analyzer/decisions/recvtype"
	"github.com/k1LoW/gostyle/analyzer/decisions/repetition"
	"github.com/k1LoW/gostyle/analyzer/decisions/underscores"
	"github.com/k1LoW/gostyle/analyzer/decisions/useany"
	"github.com/k1LoW/gostyle/analyzer/decisions/useq"
	"github.com/k1LoW/gostyle/analyzer/decisions/varnames"
	"github.com/k1LoW/gostyle/analyzer/effective/ifacenames"
	"github.com/k1LoW/gostyle/analyzer/guide/mixedcaps"
	"github.com/k1LoW/gostyle/config"
	"golang.org/x/tools/go/analysis/unitchecker"
)

//go:embed .gostyle.yml.init
var defaultConfig []byte

func main() {
	if len(os.Args) == 1 {
		if err := usage(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err) //nostyle:handlerrors
			os.Exit(1)
		}
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "init" {
		if err := generateConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err) //nostyle:handlerrors
			os.Exit(1)
		}
		os.Exit(0)
	}

	unitchecker.Main(
		config.Loader,
		dontpanic.AnalyzerWithConfig,
		errorstrings.AnalyzerWithConfig,
		getters.AnalyzerWithConfig,
		handlerrors.AnalyzerWithConfig,
		ifacenames.AnalyzerWithConfig,
		pkgnames.AnalyzerWithConfig,
		mixedcaps.AnalyzerWithConfig,
		nilslices.AnalyzerWithConfig,
		recvnames.AnalyzerWithConfig,
		recvtype.AnalyzerWithConfig,
		repetition.AnalyzerWithConfig,
		underscores.AnalyzerWithConfig,
		useany.AnalyzerWithConfig,
		useq.AnalyzerWithConfig,
		varnames.AnalyzerWithConfig,
	)
}

func usage() error {
	progname := filepath.Base(os.Args[0])
	u := `%[1]s is a set of analyzers for coding styles.

Usage of %[1]s:
	%.16[1]s init    	# generate .gostyle.yml
	%.16[1]s unit.cfg	# execute analysis specified by config file
	%.16[1]s help    	# general help, including listing analyzers and flags
	%.16[1]s help name	# help on specific analyzer and its flags
`
	if _, err := fmt.Fprintf(os.Stderr, u, progname); err != nil {
		return err
	}
	return nil
}

func generateConfig() error {
	const name = ".gostyle.yml"
	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("%s already exists", name)
	}
	if err := os.WriteFile(name, defaultConfig, os.ModePerm); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(os.Stderr, "%s is generated\n", name); err != nil {
		return err
	}
	return nil
}
