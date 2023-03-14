#!/usr/bin/env python3
#
# Copyright 2023 Miklos Vajna
#
# SPDX-License-Identifier: MIT

import io
import json
import subprocess
import sys

status = 0

for file in sys.argv[1:]:
    if file == "commands/context.go":
        continue

    buffer = io.BytesIO()
    with subprocess.Popen(['go-outline', '-f', file], stdout=subprocess.PIPE) as stream:
        buffer.write(stream.stdout.read())
    buffer.seek(0)
    j = json.load(buffer)
    assert j[0]["type"] == "package"
    package = j[0]
    children = package["children"]
    for child in children:
        if child["type"] == "variable":
            print("{}: {} is a global variable".format(file, child["label"]))
            status = 1

sys.exit(status)

# vim:set shiftwidth=4 softtabstop=4 expandtab:
