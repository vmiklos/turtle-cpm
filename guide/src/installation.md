# Installation

## Dependencies

You need to install [GPG](https://gnupg.org/) if you don't have it already and need to generate at
least a key using:

```console
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

```console
go install vmiklos.hu/go/cpm@latest
```

If `$(go env GOPATH)/bin` is not in your `PATH` yet, you may want to add it, so typing `cpm` will
invoke the installed executable.

## Optional shell completion

Optionally, you can install shell completion for cpm, example for bash:

```console
cpm completion bash > ~/.local/share/bash-completion/completions/cpm
```

You can test if it works in a new shell using:

```console
cpm <tab><tab>
```
