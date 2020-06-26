package terraform

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
)

type MultipleOptions struct {
	Options []string
}

func (v *MultipleOptions) String() string {
	return strings.Join(v.Options, ",")
}

func (v *MultipleOptions) Set(option string) error {
	*&v.Options = append(v.Options, option)
	return nil
}

func (v *MultipleOptions) Empty() bool {
	return len(v.Options) == 0
}

// Plan struct is the representation of the `terraform plan` command
type Plan struct {
	Path string
	Output string
	HasChanges bool
	CompactWarnings bool
	Destroy bool
	Input bool
	Lock bool
	LockTimeout time.Duration
	NoColor bool
	Out string
	Parallelism int
	Refresh bool
	State string
	Targets MultipleOptions
	Vars MultipleOptions
	VarFiles MultipleOptions
}

// Run runs `terraform plan` with the given arguments
func (plan *Plan) Run() error {
	if !IsCommandAvailable(plan.Path) {
		panic("Terraform not found! Be sure to put it on your $PATH.")
	}

	cmd := plan.BuildCommand()
	fmt.Printf("Running command \"%s %s\"\n", plan.Path, strings.Join(cmd, " "))

	stdout, stderr, exitCode := RunCommand(plan.Path, cmd...)

	if exitCode == 1 {
		return fmt.Errorf("Plan failed! Here is the message from Terraform:\n\n%s", stderr)
	}

	plan.HasChanges = exitCode == 2
	plan.Output = stdout

	return nil
}

func (plan *Plan) ProcessedOutput() string {
	// Gets only the plan info
	re := regexp.MustCompile(`(?ms)-+\n\n(.*)\n\n-+`)
	output := re.FindStringSubmatch(plan.Output)[1]

	if !plan.NoColor {
		// We need ro remove all ansi codes from Terraform output
		output = stripansi.Strip(output)
	}

	// Remove exceeded spaces from the beginning of the lines (runs two times)
	re = regexp.MustCompile(`(?m)^ {2}`)
	output = re.ReplaceAllString(output, "")
	output = re.ReplaceAllString(output, "")

	// Moves the change icons (+, -, ~) to the beginning of the line
	re = regexp.MustCompile(`(?m)^( +)([\+|\-|\~])`)
	output = re.ReplaceAllString(output, "$2$1")

	// Gives emphasys on what will hapen to the resource (will be purple and bold) 
	re = regexp.MustCompile(`(?m)^\#(.*)`)
	output = re.ReplaceAllString(output, "@@ #$1 @@")

	// Switches all changing lines (~) for removing and creating lines
	re = regexp.MustCompile(`(?m)^\~(.*) = (.*) -> (.*)`)
	output = re.ReplaceAllString(output, "-$1 = $2\n+$1 = $3")

	// All changing chars (~) can now be replaces by diff changing char (!)
	re = regexp.MustCompile(`(?m)^~`)
	output = re.ReplaceAllString(output, "!")

	// Replaces all replace symbols (-/+) by changing lines (!)
	re = regexp.MustCompile(`(?m)^-/\+`)
	output = re.ReplaceAllString(output, "!")

	return output
}

// BuildCommand builds the `plan` command, with all given flags
func (t Plan) BuildCommand() []string {
	// detailed exit code needed to better parse the plan
	command := []string{"plan", "-detailed-exitcode"}

	if t.CompactWarnings {
		command = append(command, "-compact-warnings")
	}

	if t.Destroy {
		command = append(command, "-destroy")
	}

	if !t.Input {
		command = append(command, "-input=false")
	}

	if t.LockTimeout.String() != "0s" {
		command = append(command, "-lock-timeout=" + t.LockTimeout.String())
	}

	if t.NoColor {
		command = append(command, "-no-color")
	}

	if t.Out != "" {
		command = append(command, "-out=" + t.Out)
	}

	if t.Parallelism != 10 {
		command = append(command, fmt.Sprintf("-parallelism=%d", t.Parallelism))
	}

	if !t.Refresh {
		command = append(command, "-refresh=false")
	}

	if t.State != "" {
		command = append(command, "-state=" + t.State)
	}

	if !t.Targets.Empty() {
		for _, v := range t.Targets.Options {
			command = append(command, "-target=" + v)
		}
	}

	if !t.Vars.Empty() {
		for _, v := range t.Vars.Options {
			command = append(command, "-var '" + v + "'")
		}
	}

	if !t.VarFiles.Empty() {
		for _, v := range t.VarFiles.Options {
			command = append(command, "-var-file=" + v)
		}
	}

	return command
}
