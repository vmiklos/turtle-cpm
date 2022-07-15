# turtle-cpm

[![tests](https://github.com/vmiklos/turtle-cpm/workflows/tests/badge.svg)](https://github.com/vmiklos/turtle-cpm/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/vmiklos/turtle-cpm)](https://goreportcard.com/report/github.com/vmiklos/turtle-cpm)

The turtle console password manager is a replacement of the now gone [cpm
project](https://www.harry-b.de/dokuwiki/doku.php?id=harry:cpm).

## Description

Notable features:

- simple DB format: encrypted sqlite (via `gpg`), tracking machines -> services -> users -> passwords

- supports plain passwords and also TOTP shared secrets, calculating the actual TOTP code via `oathtool`

- can import the original CPM's XML database

- a little bit better than trivial search: you can search for e.g. `ldap` and show all passwords
  where service is LDAP or search for e.g. `mybank` and search for all machines which contain the
  our bank domain

## The turtle

The turtle-cpm codebase is independent from the original cpm. It's turtle because this project is
not in C, so might be a little bit slower (though not significantly in practice).
