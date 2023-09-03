#!/usr/bin/env bash

# install module dependencies
go get

# install cobra-cli and inject config file
go install github.com/spf13/cobra-cli@latest
ln -s $PWD/.devcontainer/.cobra.yaml $HOME/.cobra.yaml
