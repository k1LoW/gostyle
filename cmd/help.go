package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/k1LoW/gostyle/analyzer"
	"github.com/spf13/cobra"
)

const (
	rootCommandName = "gostyle"
	helpCommandName = "help"
)

func init() {
	rootCmd.SetHelpFunc(helpCommand)
}

func helpCommand(c *cobra.Command, args []string) {
	w := os.Stdout
	// copy from spf13/cobra/command.go#defaultUsageFunc().
	fmt.Fprint(w, "Usage:")
	if c.Runnable() {
		fmt.Fprintf(w, "\n  %s", c.UseLine())
	}
	if c.HasAvailableSubCommands() {
		fmt.Fprintf(w, "\n  %s [command]", c.CommandPath())
	}
	if len(c.Aliases) > 0 {
		fmt.Fprintf(w, "\n\nAliases:\n")
		fmt.Fprintf(w, "  %s", c.NameAndAliases())
	}
	if c.HasExample() {
		fmt.Fprintf(w, "\n\nExamples:\n")
		fmt.Fprintf(w, "%s", c.Example)
	}
	if c.HasAvailableSubCommands() {
		cmds := c.Commands()
		if len(c.Groups()) == 0 {
			fmt.Fprintf(w, "\n\nAvailable Commands:")
			for _, subcmd := range cmds {
				if subcmd.IsAvailableCommand() || subcmd.Name() == helpCommandName {
					fmt.Fprintf(w, "\n  %s %s", rpad(subcmd.Name(), subcmd.NamePadding()), subcmd.Short)
				}
			}
		} else {
			for _, group := range c.Groups() {
				fmt.Fprintf(w, "\n\n%s", group.Title)
				for _, subcmd := range cmds {
					if subcmd.GroupID == group.ID && (subcmd.IsAvailableCommand() || subcmd.Name() == helpCommandName) {
						fmt.Fprintf(w, "\n  %s %s", rpad(subcmd.Name(), subcmd.NamePadding()), subcmd.Short)
					}
				}
			}
			if !c.AllChildCommandsHaveGroup() {
				fmt.Fprintf(w, "\n\nAdditional Commands:")
				for _, subcmd := range cmds {
					if subcmd.GroupID == "" && (subcmd.IsAvailableCommand() || subcmd.Name() == helpCommandName) {
						fmt.Fprintf(w, "\n  %s %s", rpad(subcmd.Name(), subcmd.NamePadding()), subcmd.Short)
					}
				}
			}
		}
	}
	if c.HasAvailableLocalFlags() {
		fmt.Fprintf(w, "\n\nFlags:\n")
		fmt.Fprint(w, trimRightSpace(c.LocalFlags().FlagUsages()))
	}
	if c.HasAvailableInheritedFlags() {
		fmt.Fprintf(w, "\n\nGlobal Flags:\n")
		fmt.Fprint(w, trimRightSpace(c.InheritedFlags().FlagUsages()))
	}
	if c.Name() == rootCommandName {
		fmt.Fprintf(w, "\n\nRegistered analyzers:")
		analyzers := analyzer.Analyzers
		padding := 0
		sort.Slice(analyzers, func(i, j int) bool {
			if len(analyzers[i].Name) > padding {
				padding = len(analyzers[i].Name)
			}
			return analyzers[i].Name < analyzers[j].Name
		})
		for _, a := range analyzers {
			if a.Name == rootCommandName {
				continue
			}
			title := strings.Split(a.Doc, "\n\n")[0]
			fmt.Fprintf(w, "\n  %s %s", rpad(a.Name, padding+1), title)
		}
	}
	if c.HasHelpSubCommands() {
		fmt.Fprintf(w, "\n\nAdditional help topics:")
		for _, subcmd := range c.Commands() {
			if subcmd.IsAdditionalHelpTopicCommand() {
				fmt.Fprintf(w, "\n  %s %s", rpad(subcmd.CommandPath(), subcmd.CommandPathPadding()), subcmd.Short)
			}
		}
	}
	if c.HasAvailableSubCommands() {
		fmt.Fprintf(w, "\n\nUse \"%s [command] --help\" for more information about a command.", c.CommandPath())
	}
	fmt.Fprintln(w)
}

// copy from spf13/cobra/command.go.
func rpad(s string, padding int) string {
	formattedString := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(formattedString, s)
}

// copy from spf13/cobra/command.go.
func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}
