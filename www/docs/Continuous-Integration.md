## Jenkins setup
![jenkins_tfswitch](static/jenkins_tfswitch.png)

```sh
#!/bin/bash 

echo "Installing tfswitch locally"
wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/master/install.sh  #Get the installer on to your machine

chmod 755 install.sh #Make installer executable

./install.sh -b `pwd`/.bin      #Install tfswitch in a location you have permission

CUSTOMBIN=`pwd`/.bin            #set custom bin path

export PATH=$PATH:$CUSTOMBIN    #Add custom bin path to PATH environment

$CUSTOMBIN/tfswitch -b $CUSTOMBIN/terraform 0.11.7 #or simply tfswitch -b $CUSTOMBIN/terraform 0.11.7

#OR 
$CUSTOMBIN/tfswitch -d 0.11.7 -b $CUSTOMBIN/terraform  #or simply tfswitch -d 0.11.7 -b $CUSTOMBIN/terraform 

terraform -v                    #testing version
```

## Circle CI setup

![cirecleci_tfswitch](static/circleci_tfswitch.png)

Example config yaml
```yaml
version: 2
jobs:
  build:
    docker:
      - image: ubuntu

    working_directory: /go/src/github.com/warrensbox/terraform-switcher

    steps:
      - checkout
      - run: 
          command: |    
            set +e   
            apt-get update 
            apt-get install -y wget 
            rm -rf /var/lib/apt/lists/*

            echo "Installing tfswitch locally"

            wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/master/install.sh  #Get the installer on to your machine

            chmod 755 install.sh            #Make installer executable

            ./install.sh -b `pwd`/.bin      #Install tfswitch in a location you have permission

            CUSTOMBIN=`pwd`/.bin            #set custom bin path

            export PATH=$PATH:$CUSTOMBIN    #Add custom bin path to PATH environment

            $CUSTOMBIN/tfswitch -b $CUSTOMBIN/terraform 0.11.7 #or simply tfswitch -b $CUSTOMBIN/terraform 0.11.7
            
            #OR 
            $CUSTOMBIN/tfswitch -d 0.11.7 -b $CUSTOMBIN/terraform  #or simply tfswitch -d 0.11.7 -b $CUSTOMBIN/terraform 

            terraform -v                    #testing version
```
