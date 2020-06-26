package terraform

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
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

// Terraform struct is the representation of the `terraform` command
type Terraform struct {
	Path string
	Output string
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
	Target MultipleOptions
	Vars MultipleOptions
	VarFiles MultipleOptions
}

// Plan runs `terraform plan` with the given arguments
func (t Terraform) Plan() {
	if !isTerraformAvailable(t.Path) {
		panic("Terraform not found! Be sure to put it on your $PATH.")
	}

	plan := t.BuildCommand()
	// cmd := exec.Command(t.Path, plan...)
	fmt.Println(plan)
}

func (t Terraform) BuildCommand() []string {
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

	if t.Parallelism > 0 {
		command = append(command, fmt.Sprintf("-parallelism=%d", t.Parallelism))
	}

	if !t.Refresh {
		command = append(command, "-refresh=false")
	}

	if t.State != "" {
		command = append(command, "-state=" + t.State)
	}

	if !t.Target.Empty() {
		for _, v := range t.Target.Options {
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

func isTerraformAvailable(path string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v " + path)

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
