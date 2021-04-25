# Jeru

Jeru helps refactor Terraform code.

One of Terraform's most valuable features is its `plan` command, which provides a safe preview of the changes it will make to infrastructure.
Unfortunately, this feature does not extend to [`state` commands for moving resources](https://www.terraform.io/docs/cli/state/move.html),
which are often necessary while refactoring Terraform code to prevent Terraform from destroying and recreating resources that are only changing addresses.


## Using Jeru

Jeru will run `terraform plan` as part of its operation in order to get an up to date plan.
Additional flags for that underlying `terraform plan` can be passed through Jeru following a double-dash (`--`).
(This is most useful if you use a variables file, ex. `terraform plan -var-file staging.tfvars`.)

### Finding possible refactors

Jeru's `find` command suggests possible `terraform state mv` commands to run given the current `terraform plan`.
In normal mode, `find` will consider all possible valid refactors and return the "best" set of commands.
In interactive mode (`--i`), `find` will walk you through the plan, providing the opportunity to match resources being deleted with resources being created.

### Planning state changes

Jeru's `plan` command provides a safe way to preview how `terraform plan` would be affected by moving resources in state.
Jeru makes a copy of the current state, applies your proposed/WIP `state` commands to that copy, and finally runs `terraform plan` against that (now-mutated) copy.

### Rolling back state changes

Even when you plan, mistakes happen.
In case you only considered the happy path, Jeru's `rollback` command generates the inverse for as many commands as it can:
- `state mv resource.a resource.b  =>  state mv resource.b resource.a` (restore the original resource address)
- `import resource.a identifier  =>  state rm resource.a` (remove the imported resource from state)
- `rm resource.a  =>  :(` (alas, since resources are imported with such specific identifiers that are not always stored in state, Jeru does not attempt to roll back removals)


## Examples

The `example` directory provides a Terraform entrypoint with a module named `original`.
Running `terraform init && terraform plan` shows the local state is up to date with no changes to make.

Try changing the module name from `original` to `new` and re-running `terraform init && terraform plan`.
Terraform now reports 8 to add and 8 to destroy.

First, Jeru can suggest some commands to run that may help:
```sh
../out/jeru find
```
You can also step through this resource-by-resource by adding the `--i` flag.
Let's write this to a file for use with the other commands:
```sh
../out/jeru find > ./move.sh
```

If we ran the `state` commands in `move.sh`, a subsequent re-run of `terraform plan` would report no changes to make, infrastructure up to date.
Before changing the actual state, though, Jeru can safely test these changes in advance:
```sh
../out/jeru plan ./move.sh
```

OK, let's run `./move.sh`.
Sure enough, a subsequent `terraform plan` now reports no changes to make... but maybe we don't like that name after all and want to keep the module named `"original"`.
We can change the address back in `main.tf` easily enough (perhaps via `git restore|revert`) and repeat the process above to create a new script of `state mv` commands to run,
or we could have Jeru reverse the changes we just made for us.
Jeru generates and executes a rollback script based on the original change script:
```sh
../out/jeru rollback ./move.sh
```


## Developing Jeru

- Build Jeru via `make build` and call it via `./out/jeru`.
- Run the tests via `make test`.
