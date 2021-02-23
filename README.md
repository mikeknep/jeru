# Jeru

Jeru helps refactor Terraform code.


## Using Jeru

### Planning state changes

One of Terraform's most valuable features is its `plan` command, which provides a safe preview of the changes it will make to infrastructure.
Unfortunately, this feature does not extend to [`state` commands for moving resources](https://www.terraform.io/docs/cli/state/move.html),
which are often necessary while refactoring Terraform code to prevent Terraform from destroying and recreating resources.

Jeru's `plan` command provides a safe way to preview how `terraform plan` would be affected by moving resources in state.
Jeru makes a copy of the current state, applies your proposed/WIP `state` commands to that copy, and finally runs `terraform plan` against that (now-mutated) copy.

Additional flags for `terraform plan` can be passed through Jeru following a double-dash (`--`).

```sh
jeru plan ./proposed-state-commands.sh -- -var-file staging.tfvars
```

### Rolling back state changes

Even when you plan, mistakes happen.
Often times when prepping `state` commands, I don't pay attention to what I'd need to run to roll back those changes.
Jeru's `rollback` command generates the inverse for as many commands as it can:
- `state mv resource.a resource.b  =>  state mv resource.b resource.a` (restore the original resource address)
- `import resource.a identifier  =>  state rm resource.a` (remove the imported resource from state)
- `rm resource.a  =>  :(` (alas, since resources are imported with such specific identifiers that are not always stored in state, Jeru does not attempt to roll back removals)

```sh
jeru rollback ./proposed-state-commands.sh
```

### Recommending possible refactors

Jeru's `recommend` command suggests possible `terraform state mv` commands to run given the current `terraform plan`.

This implementation is not particularly sophisticated yet.
Currently, if multiple resources of the same type are being renamed, the output will be a Cartesian product,
i.e. each resource being deleted will be matched with _each_ resource being created, instead of _exactly one_.
Ideally Jeru would compare more attributes of the resources under change than simply provider and type, and try to find the "exact match,"
or at the very least a set of commands that could be piped into `jeru plan`.


## Examples

The `example` directory provides a hypothetical Terraform entrypoint in an "unclean" state.
The `local_file` resource had been created with the address `"main"`, but has since been renamed `"test"`.
A standard `terraform plan` shows 1 to add and 1 to destroy.

First, Jeru can suggest some commands to run that may help:
```sh
../out/jeru recommend --out ./move.sh
```
(Note: the `move.sh` file is already checked in to this repo.)

If we ran the `state` command in the `move.sh` script, a subsequent re-run of `terraform plan` would report no changes to make, infrastructure up to date.
Jeru safely proves this in advance, before we change our actual state:
```sh
../out/jeru plan ./move.sh
```

OK, let's run `./move.sh`.
Sure enough, a subsequent `terraform plan` now reports no changes to make... but maybe we don't like that name after all and want to keep the file addressed at `"main"`.
We can change the address back in `main.tf` easily enough (perhaps via `git restore|revert`) and repeat the process above to create a new script of `state mv` commands to run,
or we could have Jeru reverse the changes we just made for us.
Jeru generates and executes a rollback script based on the original change script:
```sh
../out/jeru rollback ./move.sh --out ./rollback.sh
```
(Note: assuming `main.tf` isn't actually changed and you're just running the commands here,
this brings everything in `example/` back to the original "unclean" demo state, i.e. with 1 to add and 1 to destroy.)


## Developing Jeru

- Build Jeru via `make build` and call it via `./out/jeru`.
- Run the tests via `make test`.
