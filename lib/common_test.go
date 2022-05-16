package lib_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func checkFileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func createFile(path string) {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		defer file.Close()
	}

	fmt.Println("==> done creating file", path)
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Creating directory for terraform: %v", dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Unable to create directory for terraform: %v", dir)
			panic(err)
		}
	}
}

func cleanUp(path string) {
	err := removeContents(path)
	if err != nil {
		log.Fatalf("Error during cleanup")
	}
	removeFiles(path)
}

func removeFiles(src string) {
	files, err := filepath.Glob(src)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
