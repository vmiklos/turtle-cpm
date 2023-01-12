# Advanced topics

## Distributed usage

One big downside of mobile-based TOTP apps is that you can't login to machines if you loose your
phone, and the TOTP shared secrets on those devices are not backed up.  `cpm` supports distributed
storage of passwords in a way that's similar to single-master replication in DBMS terms.

One possible setup is to have a central machine where you edit your cpm database, e.g. your home
router which is usually available or your VPS in some hosting. Then you can replicate the `cpm`
database to your other devices by configuring a virtual `cpm` remote machine in your `.ssh/config`
on the other machines:

```
Host cpm
        Hostname myserver.example.com
```

Finally pull the remote database to your local one, using:

```console
cpm pull
```

This allows searching in your passwords even when you're offline. Keep in mind that editing the
database on the slaves is not a good idea as the next pull will overwrite your local changes.

## Toolkit integration

In case you have scripts to generate your local configuration files containing passwords from
templates, `cpm` can be integrated into such a workflow, using the quiet mode of the search
subcommand. For example, if you have an app password at your mail provider, and you want to generate
your mutt configuration, you can query just the password from `cpm` using:

```console
cpm -q -m accounts.example.com -u $USER-mail-$HOSTNAME
```

## Importing the old CPM XML database

In case you used the old `cpm` tool, it used to store its data at `~/.cpmdb` as an XML file,
compressed and encrypted. If you want to import that into turtle-cpm's database, you can do so
using:

```console
cpm import
```

## Inspecting the encrypted database manually

In case you want to inspect the SQLite database of `cpm` manually, you need to decrypt it yourself,
using (assuming an empty `XDG_STATE_HOME` environment variable):

```console
gpg --decrypt -a -o decrypted.db ~/.local/state/cpm/passwords.db
```

After this, you can inspect the database using a GUI like:

```console
sqlitebrowser decrypted.db
```

Don't forget to delete the decrypted database after you're done with your investigation.

## Reference documentation

Apart from this guide, reference documentation is available in `cpm` itself. You can learn about the
possible subcommands using:

```console
cpm -h
```

You can also check all the available options for one given subcommand using e.g.:

```console
cpm create -h
```

An alternative to this is the manual pages under `man/`, which provide the same information.

## Re-sharing TOTP shared secrets

TOTP shared secrets are typically transferred as QR codes, though there is usually a fallback option
to get the shared secret string itself, which is what `cpm` can manage. However, the QR code also
contains other information about the shared secret, and there are tools like
[2fa-qr](https://stefansundin.github.io/2fa-qr/) that allow obtaining the full `otpauth://` URL from the
QR code image. `cpm` supports storing these full URLs as well, they look something like this:

otpauth://totp/Myserver:myuser?secret=...&digits=6&algorithm=SHA1&issuer=Myserver&period=30

Where Myserver is some server-side app name and myuser is your user name.

The benefit of storing the full URL in the `cpm` database is that later you can re-share them as QR
codes using e.g.:

```console
cpm -t totp --qrcode twitter
machine: twitter.com, service: http, user: myuser, password type: TOTP shared secret, password:
...
```

The following lines will be a QR code you can scan with a mobile app.
