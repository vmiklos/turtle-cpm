# Introduction

turtle-cpm is a command line tool to manage passwords, a replacement of the now gone [cpm
project](https://www.harry-b.de/dokuwiki/doku.php?id=harry:cpm).

The latest version is v7.0, released on 2022-08-24.  See the [release notes](news.md).

Notable features:

- Simple database format: encrypted SQLite (via `gpg`), tracking machines, services, users and
  passwords
- Supports plain passwords and also Time-based one-time password (TOTP) shared secrets, calculating
  the actual TOTP code via `oathtool`
- Can import the original cpm's XML database
- A little bit better than trivial search: you can search for e.g. `ldap` or `mybank` without
  telling if you are searching for a service type or machine name

## The turtle

The turtle-cpm codebase is independent from the original cpm. It's turtle because this project is
not in C, so might be a little bit slower (though not significantly in practice).

## Contributing

turtle-cpm is free and open source. You can find the source code on
[GitHub](https://github.com/vmiklos/turtle-cpm) and issues and feature requests can be posted on the
issue tracker. If you'd like to contribute, please consider opening a pull request.

## Reporting problems

If you post a bugreport in the issue tracker, please always describe a single problem in one issue.

Here is what to include in a bugreport:

**Reproducer steps**
1. First step
2. Second step
3. Third step

**Actual result**
A description of what happens currently.

**Expected behavior**
A description of what you expected to happen.

**Version information**
The output of `cpm version`.

## License

Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.
