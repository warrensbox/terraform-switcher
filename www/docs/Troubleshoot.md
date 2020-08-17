
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

Solution: You probably need to have privileges to install *tfswitch* at /usr/local/bin.

Try the following.

```sh
wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh 
chmod 755 install.sh
./install.sh -b bin-directory

./bin-directory/tfswitch
```

Use custom directory option `-b`:    
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v7.gif" alt="drawing" style="width: 370px;"/>    


