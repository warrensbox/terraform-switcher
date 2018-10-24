package lib

import (
	"fmt"
	"log"
	"os"
)

//CreateSymlink : create symlink
func CreateSymlink(cwd string, dir string) {

	err := os.Symlink(cwd, dir)
	if err != nil {
		log.Fatal("Unable to create symlink. You must have SUDO privileges")
		panic(err)
	}
}

//RemoveSymlink : remove symlink
func RemoveSymlink(symlinkPath string) {

	_, err := os.Lstat(symlinkPath)
	if err != nil {
		log.Fatalf("Unable to find symlink. You must have SUDO privileges - %v \n", err)
		panic(err)
	} else {
		errRemove := os.Remove(symlinkPath)
		if errRemove != nil {
			log.Fatalf("Unable to remove symlink. You must have SUDO privileges - %v \n", err)
			panic(errRemove)
		}
	}
}

// CheckSymlink : check file is symlink
func CheckSymlink(symlinkPath string) bool {

	//symlink := false
	//fmt.Println("Checking symlink")

	fi, err := os.Lstat(symlinkPath)
	if err != nil {
		fmt.Println(err)
		// symlink = false
		return false
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		//symlink = true
		return true
	}

	return false
}
