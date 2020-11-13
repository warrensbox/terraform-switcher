## Step by step instructions

An open source project becomes meaningful when people collaborate to improve the code. 

Feel free to look at the code, critique and make suggestions. Lets make `tfswitch` better!

## Required version
```sh
go version 1.13
```

### Step 1 - Create workspace
*Skip this step if you already have a github go workspace*   
Create a github workspace.
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-workspace.gif" alt="drawing" style="width: 600px;"/>   

### Step 2 - Set GOPATH
*Skip this step if you already have a github go workspace*    
Export your GOPATH environment variable in your `go` directory.   
```sh
export GOPATH=`pwd`
```
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-gopath.gif" alt="drawing" style="width: 600px;"/>   

### Step 3 - Clone repository
Git clone this repository.  
```sh 
git clone git@github.com:warrensbox/terraform-switcher.git
```
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-git-clone.gif" alt="drawing" style="width: 600px;"/>  

### Step 4 - Get dependencies
Go get all the dependencies.   

```sh 
go mod download
```
```sh 
go get -v -t -d ./...
```
Test the code (optional).
```sh  
go vet -tests=false ./...
```
```sh 
go test -v ./...
```
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-go-get.gif" alt="drawing" style="width: 600px;"/>  

### Step 5 - Build executable
Create a new branch.   
```sh 
git checkout -b feature/put-your-branch-name-here
```
Refactor and add new features to the code.  
Go build the code.   
```sh 
go build -o test-tfswitch
```
Test the code and create a new pull request!

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-build.gif" alt="drawing" style="width: 600px;"/>  

### Contributors

<img style="text-allign:center" src="https://avatars3.githubusercontent.com/u/38867521?s=64&v=4" alt="drawing" width="42" height="42"/> <img style="text-allign:center" src="https://avatars2.githubusercontent.com/u/10674287?s=64&v=4" alt="drawing" width="42" height="42"/> <img style="text-allign:center" src="https://avatars1.githubusercontent.com/u/9209870?s=64&v=4" alt="drawing" width="42" height="42"/> <img style="text-allign:center" src="
https://avatars0.githubusercontent.com/u/49199497?s=64&v=4" alt="drawing" width="42" height="42"/> <img style="text-allign:center" src="https://avatars1.githubusercontent.com/u/435832?s=64&v=4" alt="drawing" width="42" height="42"/> <img style="text-allign:center" src="https://avatars1.githubusercontent.com/u/1022296?s=64&v=4" alt="drawing" width="42" height="42"/> <img style="text-allign:center" src="https://avatars2.githubusercontent.com/u/1111441?s=64&v=4" alt="drawing" width="42" height="42"/> <img style="text-allign:center" src="https://avatars1.githubusercontent.com/u/1266467?s=64&v=4" alt="drawing" width="42" height="42"/> <img style="text-allign:center" src="https://avatars3.githubusercontent.com/u/2305030?s=64&v=4" alt="drawing" width="42" height="42"/> <img style="text-allign:center" src="https://avatars1.githubusercontent.com/u/4919969?s=64&v=4" alt="drawing" width="42" height="42"/><img style="text-allign:center" src="https://avatars0.githubusercontent.com/u/12174752?s=64&v=4" alt="drawing" width="42" height="42"/>













