package main

import (
	"flag"
	"time"

	"github.com/igorbrites/neptune/terraform"
)

func main() {
	var terraform terraform.Terraform

	flag.StringVar(&terraform.Path, "path", "terraform", "Path to the \"terraform\" binnary. Be sure that it is on your $PATH.")
	flag.BoolVar(&terraform.CompactWarnings, "compact-warnings", false, "If Terraform produces any warnings that are not accompanied by errors, show them in a more compact form that includes only the summary messages.")
	flag.BoolVar(&terraform.Destroy, "destroy", false, "If set, a plan will be generated to destroy all resources managed by the given configuration and state.")
	flag.BoolVar(&terraform.Input, "input", true, "Ask for input for variables if not directly set.")
	flag.BoolVar(&terraform.Lock, "lock", true, "Lock the state file when locking is supported.")
	flag.DurationVar(&terraform.LockTimeout, "lock-timeout", time.Duration(0), "Lock the state file when locking is supported.")
	flag.BoolVar(&terraform.NoColor, "no-color", false, "If specified, output won't contain any color.")
	flag.StringVar(&terraform.Out, "out", "", "Write a plan file to the given path. This can be used as input to the \"apply\" command.")
	flag.IntVar(&terraform.Parallelism, "parallelism", 10, "Limit the number of concurrent operations. Defaults to 10.")
	flag.BoolVar(&terraform.Refresh, "refresh", true, "Resource to target. Operation will be limited to this resource and its dependencies. This flag can be used multiple times.")
	flag.StringVar(&terraform.State, "state", "", "Path to a Terraform state file to use to look up Terraform-managed resources. By default it will use the state \"terraform.tfstate\" if it exists.")
	flag.Var(&terraform.Target, "target", "Resource to target. Operation will be limited to this resource and its dependencies. This flag can be used multiple times.")
	flag.Var(&terraform.Vars, "var", "Set a variable in the Terraform configuration. This flag can be set multiple times.")
	flag.Var(&terraform.VarFiles, "var-file", "Set variables in the Terraform configuration from a file. If \"terraform.tfvars\" or any \".auto.tfvars\" files are present, they will be automatically loaded.")

	terraform.Plan()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
