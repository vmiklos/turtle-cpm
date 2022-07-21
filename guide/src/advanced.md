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

Finally put `scripts/cpmsync.sh` somewhere in your `$PATH`, so you can copy from your master to your
slave using:

```sh
cpmsync
```

This allows searching in your passwords even when you're offline. Keep in mind that editing the
database on the slaves is not a good idea as the next sync will overwrite your local changes.

## Toolkit integration

In case you have scripts to generate your local configuration files containing passwords from
templates, `cpm` can be integrated into such a workflow, using the quiet mode of the search
subcommand. For example, if you have an app password at your mail provider, and you want to generate
your mutt configuration, you can query just the password from `cpm` using:

```sh
cpm -q -m accounts.example.com -u $USER-mail-$HOSTNAME
```

## Importing the old CPM XML database

In case you used the old `cpm` tool, it used to store its data at `~/.cpmdb` as an XML file,
compressed and encrypted. If you want to import that into turtle-cpm's database, you can do so
using:

```sh
cpm import
```

## Reference documentation

Apart from this guide, reference documentation is available in `cpm` itself. You can learn about the
possible subcommands using:

```sh
cpm -h
```

You can also check all the available options for one given subcommand using e.g.:

```sh
cpm create -h
```
