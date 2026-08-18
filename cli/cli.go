// Package cli provides the integrated command-line application model for the
// Runvil meta-framework.
//
// It composes the Go standard library's flag and log/slog packages with the
// ecosystem's term package into a cohesive CLI application model.
package cli

import (
	"flag"
	"io"
	"log/slog"
	"os"

	"github.com/runvil/libs/core"
	"github.com/runvil/libs/term"
)

// Command is a runnable sub-command.
type Command struct {
	// Name is the sub-command name as typed on the command line.
	Name string
	// About is a one-line description shown in help output.
	About string
	// Register defines the command's flags on the FlagSet.
	Register func(*flag.FlagSet)
	// Run executes the command against the parsed FlagSet and returns an
	// exit code.
	Run func(*flag.FlagSet) core.ExitCode
}

// NewCommand creates a Command with the given name, about, register and run
// functions.
func NewCommand(name, about string, register func(*flag.FlagSet), run func(*flag.FlagSet) core.ExitCode) *Command {
	return &Command{Name: name, About: about, Register: register, Run: run}
}

// App is the framework's CLI application model.
type App struct {
	name     string
	version  string
	commands []*Command
	logger   *slog.Logger
	terminal *term.Terminal
	stderr   io.Writer
}

// NewApp creates an App with the given name and version.
func NewApp(name, version string) *App {
	return &App{
		name:     name,
		version:  version,
		logger:   slog.New(slog.NewTextHandler(os.Stderr, nil)),
		terminal: term.NewTerminal(),
		stderr:   os.Stderr,
	}
}

// Command registers a Command and returns the App for chaining.
func (a *App) Command(cmd *Command) *App {
	a.commands = append(a.commands, cmd)
	return a
}

// Name returns the application name.
func (a *App) Name() string {
	return a.name
}

// Run parses tokens, dispatches to the matching command, and returns the
// resulting exit code.
func (a *App) Run(tokens []string) core.ExitCode {
	if len(tokens) == 0 {
		a.printHelp()
		return core.ExitCodeUsage
	}

	name := tokens[0]
	rest := tokens[1:]

	var cmd *Command
	for _, c := range a.commands {
		if c.Name == name {
			cmd = c
			break
		}
	}
	if cmd == nil {
		a.logger.Warn("unknown command", "command", name)
		a.printHelp()
		return core.ExitCodeUsage
	}

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	if cmd.Register != nil {
		cmd.Register(fs)
	}

	a.logger.Debug("dispatching command", "command", cmd.Name)
	if err := fs.Parse(rest); err != nil {
		return core.ExitCodeUsage
	}
	return cmd.Run(fs)
}

func (a *App) printHelp() {
	title := a.terminal.Paint(a.name, term.ColorCyan, []term.Style{term.StyleBold})
	a.stderr.Write([]byte(title + " v" + a.version + "\n"))
	a.stderr.Write([]byte("Usage: " + a.name + " <command> [options]\n\n"))
	a.stderr.Write([]byte("Commands:\n"))
	for _, cmd := range a.commands {
		a.stderr.Write([]byte("  " + cmd.Name + "  " + cmd.About + "\n"))
	}
}
