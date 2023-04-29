# Introduction

turtle-cpm is a command line tool to manage passwords, a replacement of the now gone [cpm
project](https://www.harry-b.de/dokuwiki/doku.php?id=harry:cpm).

The latest version is v7.5, released on 2023-02-02.  See the [release notes](news.md).

Notable features:

- Simple database format: encrypted SQLite (via `gpg`), tracking machines, services, users and
  passwords
- Supports plain passwords and also Time-based one-time password (TOTP) shared secrets, calculating
  the actual TOTP code
- Can import the original cpm's XML database
- A little bit better than trivial search: you can search for e.g. `ldap` or `mybank` without
  telling if you are searching for a service type or machine name

Naturally it lacks telemetry, unlike e.g.
[1Password](https://blog.1password.com/privacy-preserving-app-telemetry/).

## Website

Check out the [project's website](https://vmiklos.hu/turtle-cpm/) for a list of features and
installation and usage information.

## Platforms

cpm has been used on Linux, but nothing was done intentionally to prevent it from working on e.g.
macOS or Windows.

## The important bits of the code

- The entry point is the `Main()` function in `commands/root.go`.

- The test code lives under `commands/*_test.go`.

- The documentation is undeer `guide/`.

## The turtle

The turtle-cpm codebase is independent from the original cpm. It's turtle because this project is
not in C, so might be a little bit slower (though not significantly in practice).

## License

Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.
