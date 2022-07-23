# Usage

cpm supports creating, reading, updating and deleting passwords. You'll just create and read them
most of the time, though.

## Creating

You can ask cpm to generate a password for you and remember it using:

```sh
cpm create -m mymachine -s myservice -u myuser
```

Or in case you already have a preferred password:

```sh
cpm create -m mymachine -s myservice -u myuser -p mypassword
```

A couple of useful conventions:

- `mymachine` can be e.g. the domain of the website

- `myserivce` can be e.g. just `http` if this is not some special protocol like LDAP or SSH

If you try to insert two passwords for the same machine/service/user combination, you will get an
error. You can update or delete a password, though (see below).

## Reading

You can search in your passwords by entering a search term. You can do this explicitly:

```sh
cpm search -m mymachine -s myservice -u myuser
```

Given that usually you have a single password on a website, you can be much more implicit and just
search using:

```sh
cpm mymachine
```

## TOTP support

TOTP is one from of Two-Factor Authentication (2FA), currently used by many popular websites
(Facebook, Twitter, etc). Once a website asks you to scan a QR code for
[TOTP](https://en.wikipedia.org/wiki/Time-based_one-time_password) purposes, just ask for the TOTP
shared secret and then add it to cpm using:

```sh
cpm create -m mymachine -s myservice -u myuser -p "MY TOTP SHARED SECRET" -t totp
```

When searching, it's a good idea to first narrow down your search results to a single hit, e.g.
first confirm that:

```sh
cpm twitter
```

just returns your password and your TOTP shared secret, and then you can generate the current TOTP
code using:

```sh
cpm --totp twitter
```

## Update and deletion

Update is quite similar to creation:

```sh
cpm update -m mymachine -s myservice -u myuser -p mynewpassword
```

Finally if you really want to delete a password, you can do so by using:

```sh
cpm delete -m mymachine -s myservice -u myuser
```
