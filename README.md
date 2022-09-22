# devopstools
Scripts used by DevOps team at Vega

This repository contains shared packages imported and used by scripts in other repositories.

## DevOps goals

We want to minimise usage of `bash` and replace it with `Golang`.

We assume that:
- we have `Go` installed on every machine
- we run our scripts `go run main.go ...` fashion, ineased of: `tag version`->`complie`->`publish`->`download`->`run`
- we use `cobra` to manage all CLI aspects

### Usual use-case

In a repo we create `scripts` directory, where we put DevOps Go scripts.

Example:
```bash
# Script to download latest checkpoint
go run scripts/main.go checkpoint download-latest --network fairground
```
The logic of that specific script is kept in `scripts/cmd/checkpoint.go` in that repository (or similar structure).
Shared functionality used by that script (like list of nodes, ssh to machine or download file) is imported from `devopstools` repo (i.e. here).

## Useful commands

Some useful commands not connected to any particular repo, are kept here.

### Execute command on every node

```bash
go run main.go \
    ops pssh \
    --ssh-user fryderyk \
    --ssh-private-keyfile ~/.ssh/id_ed25519 \
    --network devnet1 \
    --command "pwd"
```

### Vega Network statistics

```bash
# all stats
go run main.go live statistics --network devnet1
# version only
go run main.go live statistics --network devnet1 --version
# block only
go run main.go live statistics --network devnet1 --block
# from all nodes
go run main.go live statistics --network devnet1 --all
```

### Vega Network - list of nodes

```bash
# list all healthy nodes
go run main.go live nodename --network fairground
# list all nodes
go run main.go live nodename --network devnet1 --all
```