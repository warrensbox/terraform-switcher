package lib

import (
	"os"
	"path/filepath"
)

func checkFileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}

func createFile(path string) {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			logger.Error(err)
			return
		}
		defer file.Close()
	}

	logger.Infof("==> done creating %q file", path)
}

func cleanUp(path string) {
	removeContents(path)
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
