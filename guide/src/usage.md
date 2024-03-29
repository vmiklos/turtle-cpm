# Usage

cpm supports creating, reading, updating and deleting passwords. You'll just create and read them
most of the time, though.

## Creating

You can ask cpm to generate a password for you and remember it using:

```console
cpm create
```

You'll have to provide the machine and the user:

```console
Machine: example.com
User: myuser
Generated password: 7U1FvIzubR95Itg
```

Specifying parameters can be useful if:

- you want to avoid an interactive question for the machine
- you want to specify a non-HTTP service
- you want to avoid an interactive question for the user
- you want to specify a password type
- you have preferred password

Example for such usage:

```console
cpm create -m example.com -s http -u myuser -t plain -p 7U1FvIzubR95Itg
```

When the machine is not yours, it can be e.g. the domain of a website.

If you try to insert a password twice (same machine, service, user and password type), you will get
an error. You can update or delete a password, though (see below).

## Reading

You can search in your passwords by entering a search term. You can do this interactively:

```console
cpm
```

You'll have to provide a search term:

```console
Search term: example.com
id:        1, machine: example.com, service: http, user: myuser, password type: plain, password: 7U1FvIzubR95Itg
```

The search term can also be specified as an argument if non-interactive mode is wanted.

Or you can specify parameters to create additional filters for the search:

```console
cpm search -m example.com -s http -u myuser -t plain
```

The search term is already specified in this case:

```
id:        1, machine: example.com, service: http, user: myuser, password type: plain, password: 7U1FvIzubR95Itg
```

Archived passwords are not shown, unless `-v` or `--verbose` is used.

## TOTP support

TOTP is one from of Two-Factor Authentication (2FA), currently used by many popular websites
(Facebook, Mastodon, etc). Once a website asks you to scan a QR code for
[TOTP](https://en.wikipedia.org/wiki/Time-based_one-time_password) purposes, just ask for the TOTP
shared secret and then add it to cpm using:

```console
cpm create -m mymachine -u myuser -p "MY TOTP SHARED SECRET" -t totp
```

When searching, only the TOTP shared secret is shown by default:

```console
cpm facebook
id:        1, machine: facebook.com, service: http, user: myuser, password type: plain, password: 7U1FvIzubR95Itg
id:        2, machine: facebook.com, service: http, user: myuser, password type: TOTP shared secret, password: ...
```

You can generate the current TOTP code using:

```console
cpm --totp facebook
id:        2, machine: facebook.com, service: http, user: myuser, password type: TOTP code, password: ...
```

You can make this interaction easier using:

```console
alias 2fa='cpm --totp'
```

And then you can generate the current TOTP code just by:

```console
2fa facebook
id:        2, machine: facebook.com, service: http, user: myuser, password type: TOTP code, password: ...
```

## Update and deletion

Update is quite similar to creation. If you want to update a password to a new, generated value, you
can do so by using:

```console
cpm update -p -
```

Notice the trailing hyphen.

You'll have to specify the ID:

```console:
Id: 2
Updated 1 password
Generated password: aDu3WwGlVP60HEn
```

You can also specify more parameters for `cpm update`:

```console
cpm update -i 2 -p -
```

In which case the command is not interactive:

```console
Updated 1 password
Generated password: Ilsd08zGov5JyBR
```

You can use `cpm search` to find the password ID.

The rest of the `cpm update` parameters allow explicitly setting the
machine/service/user/type/password of an ID to a new, specified value.

Finally if you want to delete a password, you can do so by using:

```console
cpm delete
```

You'll have to specify the ID:

```console:
Id: 2
Deleted 1 password
```

You can also specify a parameter for `cpm delete`:

```console
cpm delete -i 2
```

In which case the command is not interactive:

```console
Deleted 1 password
```

Again, you can use `cpm search` to find the password ID.

An alternative for deletion is to just mark the password as archived:

```console
cpm update -i ... -a 1
```
