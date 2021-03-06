# Usage

cpm supports creating, reading, updating and deleting passwords. You'll just create and read them
most of the time, though.

## Creating

You can ask cpm to generate a password for you and remember it using:

```sh
cpm create -m mymachine -u myuser
```

Or in case you already have non-HTTP service or a preferred password:

```sh
cpm create -m mymachine -s myservice -u myuser -p mypassword
```

When the machine is not yours, it can be e.g. the domain of a website.

If you try to insert two passwords for the same machine/service/user combination, you will get an
error. You can update or delete a password, though (see below).

## Reading

You can search in your passwords by entering a search term. You can do this implicitly:

```sh
cpm mymachine
```

In the less likely case when you have multiple passwords on a website, you can be much more explicit
and search using:


```sh
cpm search -m mymachine -s myservice -u myuser
```

## TOTP support

TOTP is one from of Two-Factor Authentication (2FA), currently used by many popular websites
(Facebook, Twitter, etc). Once a website asks you to scan a QR code for
[TOTP](https://en.wikipedia.org/wiki/Time-based_one-time_password) purposes, just ask for the TOTP
shared secret and then add it to cpm using:

```sh
cpm create -m mymachine -u myuser -p "MY TOTP SHARED SECRET" -t totp
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

Update is quite similar to creation. You can generate a new password using:

```sh
cpm update -m mymachine -u myuser
```

If you want to specify a service or a new password explicitly, you can do that using:

```sh
cpm update -m mymachine -s myservice -u myuser -p mynewpassword
```

Finally if you want to delete a password, you can do so by using:

```sh
cpm delete -m mymachine -u myuser
```

Or if you want to specify the service explicitly:

```sh
cpm delete -m mymachine -s myservice -u myuser
```
