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
