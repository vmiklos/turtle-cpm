# Changelog

## 7.6

- `create` / `create -n` now confirms if the password was actually created or it only would be
  created.
- [pwgen](http://sourceforge.net/projects/pwgen/) installation is not necessary, `cpm` now uses
  [go-password](https://github.com/sethvargo/go-password) instead.
- `update` now has only a single way to select which password to update: the ID has to be specified
  using `-i` or interactively. This allows updating the password or type to a specified value while
  only specifying an ID.
- `delete` now has only a single way to select which password to delete: the ID has to be specified
  using `-i` or interactively.
- [oathtool](https://www.nongnu.org/oath-toolkit/) installation is no longer necessary, `cpm` now
  uses [otp](https://github.com/pquerna/otp) instead.

## 7.5

- update: new `-i` switch to edit the `machine`, `service` or `user` of a password without changing
  the password itself
- `sync` is renamed to `pull`, since it's a remote -> local copy only
- delete: new `-i` switch to delete a password based on the search result

## 7.4

- added shell completion support

## 7.3

- create: new `-n` switch, introducing a dry run mode to see what style of password would be
  generated
- update and delete: new `-n` switch, showing how many passwords would be updated/deleted (typically
  0 or 1), without writing the database

## 7.2

- in case the `XDG_STATE_HOME` environment variable is set to a custom value, it is now respected

## 7.1

- when specifying TOTP shared secrets, it is now supported to specify `otpauth://` URLs
- new `cpm --qrcode` option to show the TOTP shared secret as a qr code

## 7.0

- new `version` subcommand that shows the version number

## 6.0

- create, update and delete: add interactive mode in case `--machine` or `--user` is not specified
- search: add interactive mode when no search terms are specified

## 5.0

- `scripts/cpmsync.sh` is now a built-in subcommand, `cpm sync`. No need to install it manually
  anymore.
- `cpm search` no longer encrypts the database at the end, to be a little bit faster.
- create/search/update/delete's -t flag now has a proper type, so it errors on invalid values (other
  than `plain` or `totp`)

## 4.0

- New password generation during update, the password flag is now optional.
- The service flag is now optional in general, and defaults to "http".

## 3.0

- Added overview documentation to augment the existing automatically generated reference
  documentation. (`cpm -h`, `cpm create -h`, etc.)
- New quiet mode during search, to consume the output from scripts
- Added manpages
- New password generation during create, the password flag is now optional.

## 2.0

- While encrypting the database, the `gpg` invocation now defaults to the self recipient, so no need
  to specify a UID manually.
- Search now shows both plain passwords and TOTP shared secrets by default, so in case a site only
  has a TOTP shared secret, there is a search result instead of confusing empty output.
- The decrypted database is removed from disk even in case of application failure.
- More tests: 100% statement coverage.

## 1.0

- Initial release
