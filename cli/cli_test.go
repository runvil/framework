package cli

import (
	"flag"
	"io"
	"testing"

	"github.com/runvil/runvil-libs/core"
)

func successHandler(*flag.FlagSet) core.ExitCode {
	return core.ExitCodeSuccess
}

func testApp() *App {
	app := NewApp("demo", "0.1.0").
		Command(NewCommand("ping", "ping the app", nil, successHandler))
	app.stderr = io.Discard
	return app
}

func TestRunDispatchesToKnownCommand(t *testing.T) {
	app := testApp()
	if got := app.Run([]string{"ping"}); got != core.ExitCodeSuccess {
		t.Errorf("Run() = %v, want %v", got, core.ExitCodeSuccess)
	}
}

func TestRunEmptyIsUsage(t *testing.T) {
	app := testApp()
	if got := app.Run(nil); got != core.ExitCodeUsage {
		t.Errorf("Run() = %v, want %v", got, core.ExitCodeUsage)
	}
}

func TestRunUnknownCommandIsUsage(t *testing.T) {
	app := testApp()
	if got := app.Run([]string{"nope"}); got != core.ExitCodeUsage {
		t.Errorf("Run() = %v, want %v", got, core.ExitCodeUsage)
	}
}

func TestRunParseErrorIsUsage(t *testing.T) {
	app := NewApp("demo", "0.1.0").
		Command(NewCommand("ping", "ping the app", func(fs *flag.FlagSet) {
			fs.String("name", "", "name of the target")
		}, successHandler))
	app.stderr = io.Discard
	if got := app.Run([]string{"ping", "--nope"}); got != core.ExitCodeUsage {
		t.Errorf("Run() = %v, want %v", got, core.ExitCodeUsage)
	}
}
