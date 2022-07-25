# Installation

cpm has been used on Linux, but nothing was done intentionally to prevent it from working on e.g.
macOS or Windows.

## Dependencies

You need to install [GPG](https://gnupg.org/) if you don't have it already and need to generate at
least a key using:

```sh
gpg --gen-key
```

This will allow cpm to encrypt and decrypt your password database. Don't fear from specifying a
passphrase for your key: you'll only have to enter it once after login (thanks to
[gpg-agent](https://www.gnupg.org/documentation/manuals/gnupg/Invoking-GPG_002dAGENT.html)), and
this way an attacker can't get access to your passwords, even if they steal both your private key &
your cpm database.

Optionally, you can install [oathtool](https://www.nongnu.org/oath-toolkit/) in case you're
interested in TOTP support.

Optionally, you can install [pwgen](http://sourceforge.net/projects/pwgen/) in case you want to have
auto-generated passwords.

## Build from source

You can install cpm using:

```sh
go install vmiklos.hu/go/turtle-cpm@latest
```

If you don't have the original cpm around, you can make a symlink to allow less typing:

```sh
ln -s $(go env GOPATH)/bin/turtle-cpm ~/bin/cpm
```
