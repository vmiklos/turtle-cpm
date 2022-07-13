# turtle-cpm

The turtle console password manager is a replacement of the now gone
<https://www.harry-b.de/dokuwiki/doku.php?id=harry:cpm>.

Notable features:

- simple DB format: encrypted sqlite (via `gpg`), tracking machines -> services -> users -> passwords

- supports plain passwords and also TOTP secrets, via `oathtool`

- can import the original CPM's XML database

- a little bit better than trivial search: you can search for e.g. `ldap` and show all passwords
  where service is LDAP or search for e.g. `mybank` and search for all machines which contain the
  our bank domain
