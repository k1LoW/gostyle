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
	"fmt"
	"os"

	"github.com/k1LoW/gostyle/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate .gostyle.yml",
	Long:  `Generate .gostyle.yml.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := generateConfig(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func generateConfig() error {
	const name = ".gostyle.yml"
	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("%s already exists", name)
	}
	if err := os.WriteFile(name, config.Default, os.ModePerm); err != nil { //nolint:gosec
		return err
	}
	if _, err := fmt.Fprintf(os.Stderr, "%s is generated\n", name); err != nil {
		return err
	}
	return nil
}
