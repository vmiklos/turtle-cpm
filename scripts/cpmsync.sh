#!/bin/bash -ex
#
# Copyright 2022 Miklos Vajna. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
#

#
# This script copies the cpm db from a remote machine to the current one.
#

scp cpm:.local/state/cpm/passwords.db ~/.local/state/cpm/passwords.db

# vim:set shiftwidth=4 softtabstop=4 expandtab:
