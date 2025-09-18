/*
Copyright © 2025 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/k1LoW/gostyle/analyzer"
	"github.com/k1LoW/gostyle/config"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/analysis/multichecker"
)

var configPath string

var runCmd = &cobra.Command{
	Use:   "run [packages]",
	Short: "Run analyzers",
	Long:  `Run analyzers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if configPath != "" && !filepath.IsAbs(configPath) {
			abs, err := filepath.Abs(configPath)
			if err != nil {
				return err
			}
			configPath = abs
		}
		config.SetPath(configPath)
		if len(args) == 0 {
			args = []string{"."}
		}
		os.Args = append([]string{"gostlye"}, args...)
		multichecker.Main(analyzer.Analyzers...)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&configPath, "config", "c", "", "path of config file")
}
