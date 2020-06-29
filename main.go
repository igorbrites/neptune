package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/igorbrites/neptune/github"
	"github.com/igorbrites/neptune/terraform"
)

func main() {
	var plan terraform.Plan
	var pr github.PullRequest
	parseFlags(&plan, &pr)

	plan.Run()

	fmt.Println(plan.Output)

	pr.Comment(plan)

	if plan.Type == terraform.Error {
		os.Exit(1)
	}
}

func parseFlags(plan *terraform.Plan, pr *github.PullRequest) {
	flag.StringVar(&plan.Path, "path", "terraform", "Path to the \"terraform\" binnary. Be sure that it is on your $PATH.")
	flag.BoolVar(&plan.CompactWarnings, "compact-warnings", false, "If Terraform produces any warnings that are not accompanied by errors, show them in a more compact form that includes only the summary messages.")
	flag.BoolVar(&plan.Destroy, "destroy", false, "If set, a plan will be generated to destroy all resources managed by the given configuration and state.")
	flag.BoolVar(&plan.Input, "input", true, "Ask for input for variables if not directly set.")
	flag.BoolVar(&plan.Lock, "lock", true, "Lock the state file when locking is supported.")
	flag.DurationVar(&plan.LockTimeout, "lock-timeout", time.Duration(0), "Lock the state file when locking is supported.")
	flag.BoolVar(&plan.NoColor, "no-color", false, "If specified, output won't contain any color.")
	flag.StringVar(&plan.Out, "out", "", "Write a plan file to the given path. This can be used as input to the \"apply\" command.")
	flag.IntVar(&plan.Parallelism, "parallelism", 10, "Limit the number of concurrent operations. Defaults to 10.")
	flag.BoolVar(&plan.Refresh, "refresh", true, "Resource to target. Operation will be limited to this resource and its dependencies. This flag can be used multiple times.")
	flag.StringVar(&plan.State, "state", "", "Path to a Terraform state file to use to look up Terraform-managed resources. By default it will use the state \"terraform.tfstate\" if it exists.")
	flag.Var(&plan.Targets, "target", "Resource to target. Operation will be limited to this resource and its dependencies. This flag can be used multiple times.")
	flag.Var(&plan.Vars, "var", "Set a variable in the Terraform configuration. This flag can be set multiple times.")
	flag.Var(&plan.VarFiles, "var-file", "Set variables in the Terraform configuration from a file. If \"terraform.tfvars\" or any \".auto.tfvars\" files are present, they will be automatically loaded.")

	flag.StringVar(&pr.Owner, "owner", "", "Name of the owner of the repo")
	flag.StringVar(&pr.Repo, "repo", "", "Repo name")
	flag.IntVar(&pr.Number, "pr-number", 0, "Pull Request number")

	flag.Parse()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
