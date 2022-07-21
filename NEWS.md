# Version descriptions

## main

- Added overview documentation to augment the existing automatically generated reference
  documentation. (`cpm -h`, `cpm create -h`, etc.)
- New quite mode during search, to consume the output from scripts.

## 2.0

- While encrypting the database, the `gpg` invocation now defaults to the self recipient, so no need
  to specify a UID manually.
- Search now shows both plain passwords and TOTP shared secrets by default, so in case a site only
  has a TOTP shared secret, there is a search result instead of confusing empty output.
- The decrypted database is removed from disk even in case of application failure.
- More tests: 100% statement coverage.

## 1.0

- Initial release
