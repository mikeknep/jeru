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
jeru plan --changes ./proposed-state-commands.sh -- -var-file=foo
```

#### example

The `example` directory provides a hypothetical Terraform entrypoint in an unclean state.
The `local_file` resource had been created with the name `"main"`, but has been renamed `"test"`.
A standard `terraform plan` run will show 1 to add and 1 to destroy.
However, if we ran the `state` command in `example/move.sh`, a re-ran of `plan` should show no changes to make.

Jeru proves this!
From the `example` directory, run:
```sh
../out/jeru plan --changes ./move.sh
```


## Developing Jeru

Build Jeru via `make build` and call it via `./out/jeru`.
