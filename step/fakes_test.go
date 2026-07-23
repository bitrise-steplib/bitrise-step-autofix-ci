package step

import (
	"strings"

	"github.com/bitrise-io/go-utils/v2/command"
)

// fakeEnvRepo is a simple map-backed env.Repository for tests.
type fakeEnvRepo map[string]string

func (m fakeEnvRepo) Get(key string) string { return m[key] }
func (m fakeEnvRepo) List() []string        { return nil }
func (m fakeEnvRepo) Set(k, v string) error { return nil }
func (m fakeEnvRepo) Unset(k string) error  { return nil }

type capturedCall struct {
	name string
	args []string
	opts *command.Opts
}

type fakeCommandFactory struct {
	calls     []capturedCall
	responses map[string]string // maps git arg keyword to output returned by RunAndReturnTrimmedCombinedOutput
	errors    map[string]error  // maps git arg keyword to error returned by the run methods
}

func (f *fakeCommandFactory) Create(name string, args []string, opts *command.Opts) command.Command {
	f.calls = append(f.calls, capturedCall{name, args, opts})
	if f.errors != nil {
		for _, arg := range args {
			if err, ok := f.errors[arg]; ok {
				return &noopCommand{err: err}
			}
		}
	}
	if f.responses != nil {
		for _, arg := range args {
			if resp, ok := f.responses[arg]; ok {
				return &noopCommand{output: resp}
			}
		}
	}
	return &noopCommand{}
}

// findCall returns the first recorded Create call whose args contain gitSubcmd.
func (f *fakeCommandFactory) findCall(gitSubcmd string) (capturedCall, bool) {
	for _, c := range f.calls {
		for _, arg := range c.args {
			if arg == gitSubcmd {
				return c, true
			}
		}
	}
	return capturedCall{}, false
}

type noopCommand struct {
	output string
	err    error
}

func (c *noopCommand) PrintableCommandArgs() string { return "" }
func (c *noopCommand) Run() error                   { return c.err }
func (c *noopCommand) RunAndReturnExitCode() (int, error) {
	if c.err != nil {
		return 1, c.err
	}
	return 0, nil
}
func (c *noopCommand) RunAndReturnTrimmedOutput() (string, error) { return c.output, c.err }
func (c *noopCommand) RunAndReturnTrimmedCombinedOutput() (string, error) {
	return c.output, c.err
}
func (c *noopCommand) Start() error { return nil }
func (c *noopCommand) Wait() error  { return nil }

// credentialHelperArg returns the value of the first "-c credential.helper=..." pair
// found in args, or empty string if none is present.
func credentialHelperArg(args []string) string {
	for i, arg := range args {
		if arg == "-c" && i+1 < len(args) && strings.HasPrefix(args[i+1], "credential.helper=") {
			return args[i+1]
		}
	}
	return ""
}

func envContainsPrefix(envs []string, prefix string) bool {
	for _, e := range envs {
		if strings.HasPrefix(e, prefix) {
			return true
		}
	}
	return false
}
