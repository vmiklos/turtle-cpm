# Version descriptions

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
