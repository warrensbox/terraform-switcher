### Jenkins setup
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/jenkins_tfswitch.png" alt="drawing" style="width: 370px;"/>

```sh
#!/bin/bash 

echo "Installing tfswitch locally"
wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh 
chmod 755 install.sh
./install.sh -b bin-directory

./bin-directory/tfswitch
```
If you have limited permission, try:

```sh
#!/bin/bash 

echo "Installing tfswitch locally"
wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh 
chmod 755 install.sh
./install.sh -b bin-directory

CUSTOMBIN=`pwd`/bin             #set custom bin path
mkdir $CUSTOMBIN                #create custom bin path
export PATH=$PATH:$CUSTOMBIN    #Add custom bin path to PATH environment

./bin-directory/tfswitch -b $CUSTOMBIN/terraform 0.11.7

terraform -v                    #testing version
```

### Circle CI setup

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/circleci_tfswitch.png" alt="drawing" style="width: 470px;"/>


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

            wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh 
            chmod 755 install.sh
            ./install.sh -b bin-directory

            CUSTOMBIN=`pwd`/bin             #set custom bin path
            mkdir $CUSTOMBIN                #create custom bin path
            export PATH=$PATH:$CUSTOMBIN    #Add custom bin path to PATH environment

            ./bin-directory/tfswitch -b $CUSTOMBIN/terraform 0.11.7

            terraform -v                    #testing version
```