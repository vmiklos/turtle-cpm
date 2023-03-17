#!/bin/bash -ex
#
# Copyright 2022 Miklos Vajna
#
# SPDX-License-Identifier: MIT
#

#
# This script runs all the tests for CI purposes.
#

if [ -n "${GITHUB_WORKFLOW}" ]; then
    go install golang.org/x/lint/golint@latest
    go install github.com/dave/courtney@latest
    go install github.com/google/addlicense@latest
fi

make check

# vim:set shiftwidth=4 softtabstop=4 expandtab:
