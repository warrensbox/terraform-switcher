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
