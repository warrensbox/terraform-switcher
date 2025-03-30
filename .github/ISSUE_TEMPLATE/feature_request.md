---
name: Feature request
about: Suggest an idea for this project
title: ""
labels: enhancement
assignees: warrensbox
---

**Is your feature request related to a problem? Please describe.**
A clear and concise description of what the problem is. Ex. I'm always frustrated when [...]

**Describe the solution you'd like**
A clear and concise description of what you want to happen.

**Describe alternatives you've considered**
A clear and concise description of any alternative solutions or features you've considered.

**Additional context**
Add any other context or screenshots about the feature request here.

### If you would like to contribute to the code, see step-by-step instructions here

## Required version

```sh
go version 1.13
```

### Step 1 - Create workspace

_Skip this step if you already have a GitHub Go workspace_  
Create a GitHub workspace.
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-workspace.gif" alt="drawing" style="width: 600px;"/>

### Step 2 - Set GOPATH

_Skip this step if you already have a GitHub Go workspace_  
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
