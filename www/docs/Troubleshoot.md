<!-- markdownlint-disable MD041 -->

Problem:

```sh
install: can't change permissions of /usr/local/bin: Operation not permitted
```

```sh
"Unable to remove symlink. You must have SUDO privileges"
```

```sh
"Unable to create symlink. You must have SUDO privileges"
```

```sh
install: cannot create regular file '/usr/local/bin/tfswitch': Permission denied
```

Solution: You probably need to have privileges to install _tfswitch_ at /usr/local/bin.

Try the following:

```sh
wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/master/install.sh  #Get the installer on to your machine:

chmod 755 install.sh #Make installer executable

./install.sh -b $HOME/.bin #Install tfswitch in a location you have permission:

$HOME/.bin/tfswitch #test

export PATH=$PATH:$HOME/.bin #Export your .bin into your path

#You should probably add step 4 in your `.bash_profile` in your $HOME directory.

#Next, try:
`tfswitch -b $HOME/.bin/terraform 0.11.7`

#or simply

`tfswitch -b $HOME/.bin/terraform`


```

See the custom directory option `-b`:  
![custom directory](static/tfswitch-v7.gif "Custom binary path")

---

Problem:

When installing Terraform 1.14.9 or later using tfswitch v1.16.0 or earlier, or after upgrading to v1.17.0+ without clearing the locally cached HashiCorp PGP key:

```sh
ERROR Could not verify PGP signature using key №1 (out of 1): Signature Verification Error: Invalid signature caused by openpgp: key expired
FATAL Error downloading: Unable to verify checksum signature against PGP key
```

Solution: Your locally cached copy of HashiCorp's PGP public key is stale. Starting with Terraform 1.14.9, HashiCorp began signing releases with a refreshed key block after the original key block expired. This was fixed in tfswitch v1.17.0 (#747), but upgrading alone is not enough — if tfswitch had already cached the key before that point, the stale single-block file remains on disk and the same error persists. Delete the cached file so tfswitch re-downloads the current key:

```sh
rm ~/.terraform.versions/terraform_72D7468F.asc
tfswitch <version>
```

