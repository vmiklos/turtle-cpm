# Usage

cpm supports creating, reading, updating and deleting passwords. You'll just create and read them
most of the time, though.

## Creating

You can ask cpm to generate a password for you and remember it using:

```sh
$ cpm create
Machine: example.com
User: myuser
Generated password: 7U1FvIzubR95Itg
```

Specifying parameters can be useful if:

- you want to avoid interactive questions
- you want to specify a non-HTTP service or
- you have preferred password

Example for such usage:


```sh
$ cpm create -m example.com -s http -u myuser -p 7U1FvIzubR95Itg
```

When the machine is not yours, it can be e.g. the domain of a website.

If you try to insert a password twice (same machine, service, user and password type), you will get
an error. You can update or delete a password, though (see below).

## Reading

You can search in your passwords by entering a search term. You can do this implicitly:

```sh
$ cpm example.com
machine: example.com, service: http, user: myuser, password type: plain, password: 7U1FvIzubR95Itg
```

In the less likely case when you have multiple passwords on a website or you want to hide TOTP
shared secrets, you can be much more explicit and search using:

```sh
$ cpm search -m example.com -s http -u myuser -t plain
machine: example.com, service: http, user: myuser, password type: plain, password: 7U1FvIzubR95Itg
```

## TOTP support

TOTP is one from of Two-Factor Authentication (2FA), currently used by many popular websites
(Facebook, Twitter, etc). Once a website asks you to scan a QR code for
[TOTP](https://en.wikipedia.org/wiki/Time-based_one-time_password) purposes, just ask for the TOTP
shared secret and then add it to cpm using:

```sh
$ cpm create -m mymachine -u myuser -p "MY TOTP SHARED SECRET" -t totp
```

When searching, it's a good idea to first narrow down your search results to a single hit, e.g.
first confirm that:

```sh
$ cpm twitter
machine: twitter.com, service: http, user: myuser, password type: plain, password: ...
machine: twitter.com, service: http, user: myuser, password type: TOTP shared secret, password: ...
```

just returns your password and your TOTP shared secret, and then you can generate the current TOTP
code using:

```sh
$ cpm --totp twitter
machine: twitter.com, service: http, user: myuser, password type: TOTP code, password: ...
```

## Update and deletion

Update is quite similar to creation. You can generate a new password using:

```sh
$ cpm update -m example.com -u myuser
Generated new password: D95Rx2PlOcPwKbL
Updated 1 password
```

If you want to specify a service, a type or a new password explicitly, you can do that using:

```sh
$ cpm update -m mymachine -s myservice -u myuser -t plain -p mynewpassword
Updated 1 password
```

Finally if you want to delete a password, you can do so by using:

```sh
$ cpm delete -m example.com -u myuser
Deleted 1 password
```

Or if you want to specify the service or type explicitly:

```sh
$ cpm delete -m mymachine -s myservice -u myuser -t plain
Deleted 1 password
```
