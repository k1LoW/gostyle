package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/k1LoW/gostyle/analyzer/decisions/pkgnames"
	"github.com/k1LoW/gostyle/analyzer/decisions/recvnames"
	"github.com/k1LoW/gostyle/analyzer/decisions/underscores"
	"github.com/k1LoW/gostyle/analyzer/effective/ifacenames"
	"github.com/k1LoW/gostyle/analyzer/guide/mixedcaps"
	"github.com/k1LoW/gostyle/config"
	"golang.org/x/tools/go/analysis/unitchecker"

	_ "embed"
)

const configPath = ".gostyle.yml"

//go:embed .gostyle.yml
var defaultConfig []byte

func main() {
	progname := filepath.Base(os.Args[0])

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `%[1]s is a set of analyzers for coding styles.

Usage of %[1]s:
	%.16[1]s unit.cfg	# execute analysis specified by config file
	%.16[1]s help    	# general help, including listing analyzers and flags
	%.16[1]s help name	# help on specific analyzer and its flags
	%.16[1]s init    	# generate config file for -gostyle.config flag
`, progname)
		os.Exit(1)
	}

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
	}

	if args[0] == "init" {
		if _, err := os.Stat(configPath); err == nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s already exists\n", configPath)
			os.Exit(1)
		}
		fmt.Printf("%s is generated\n", configPath)
		os.WriteFile(configPath, defaultConfig, os.ModePerm)
		os.Exit(0)
	}

	unitchecker.Main(
		config.Loader,
		ifacenames.Analyzer,
		pkgnames.Analyzer,
		mixedcaps.Analyzer,
		recvnames.Analyzer,
		underscores.Analyzer,
	)
}
