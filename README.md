# Neptune


![GoBuild](https://github.com/ContaAzul/neptune/workflows/GoBuild/badge.svg) ![GoReleaser](https://github.com/ContaAzul/neptune/workflows/GoReleaser/badge.svg)

Neptune runs plans with Terraform and posts its results on Pull Requests on GitHub (very specific, huh?).

## Installing

Go check the [releases](https://github.com/ContaAzul/neptune/releases) to get the latest version for your OS.

## Using

Neptune accepts all the flags of the [`terraform plan`](https://www.terraform.io/docs/commands/plan.html) command, plus the following:

- `-owner`: Name of the owner of the repo;
- `-path`: Path to the "terraform" binnary. Be sure that it is on your $PATH. (default "terraform");
- `-pr-number`: Pull Request number;
- `-repo`: Repo name.

So, you can simply change the `terraform plan` by `neptune`, adding the flags to post on the Pull Request:

```sh
# From:
terraform plan -out=terraform.tfplan -var-file=terraform.tfvars -input=false
# To:
neptune -out=terraform.tfplan -var-file=terraform.tfvars -input=false -owner=ContaAzul -repo=neptune -pr-number=1
```
