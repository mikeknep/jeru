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
However, if we ran the `state` command in `example/move.sh`, a re-run of `plan` should show no changes to make.

Jeru proves this!
From the `example` directory, run:
```sh
../out/jeru plan --changes ./move.sh
```

### Recommending possible refactors

Jeru's `recommend` command suggests possible `terraform state mv` commands to run given the current `terraform plan`.

This implementation is not particularly sophisticated yet.
Currently, if multiple resources of the same type are being renamed, the output will be a Cartesian product,
i.e. each resource being deleted will be matched with _each_ resource being created, instead of _exactly one_.
Ideally Jeru would compare more attributes of the resources under change than simply provider and type, and try to find the "exact match,"
or at the very least a set of commands that could be piped into `jeru plan`.

#### example

As mentioned above, the `local_file` resource in the `example` directory entrypoint has been renamed.
Jeru will recognize this and recommend the appropriate `terraform state mv` command to run.
(It is the same command "proposed" in `example/move.sh` and proven via `jeru plan` to lead to no changes to make.)

From the `example` directory, run:
```sh
../out/jeru recommend
```


## Developing Jeru

- Build Jeru via `make build` and call it via `./out/jeru`.
- Run the tests via `make test`.
