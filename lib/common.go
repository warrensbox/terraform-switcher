package lib

import (
	"os"
	"path/filepath"

	"github.com/pborman/getopt/v2"
	"golang.org/x/term"
)

// Default value used as fallback in `getTerminalWidth()`
const defaultTerminalWidth = 80

func checkFileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func createFile(path string) {
	// detect if file exists
	_, err := os.Stat(path)

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
	err := removeContents(path)
	if err != nil {
		logger.Error(err)
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

// Returns the width of the terminal backing the given file (typically os.Stdout or os.Stderr),
// or a fallback value if it cannot be determined
// (or if the terminal is too narrow for `getopt` to wrap help text in a readable way)
// Callers must pass the same file they intend to write the width-formatted output to,
// since stdout and stderr can be redirected independently and may not share a terminal.
// For user convenience, the fallback value can be specified as an argument (or otherwise defaults to `defaultTerminalWidth`)
// Errors out if more than one fallback value is specified
func getTerminalWidth(file *os.File, fallback ...int) int {
	if len(fallback) > 1 {
		logger.Warnf("getTerminalWidth received %d fallback arguments, ignoring extras", len(fallback))
	}

	indentCols := 5 // Indent 5 columns from right side of the window for the readability's sake
	if width, _, err := term.GetSize(int(file.Fd())); err == nil && width > indentCols {
		// Below `getopt.HelpColumn`, there's no room left for help text, so `getopt` wraps it
		// one character per line. Treat that as if detection failed and fall back below:
		// a too-wide-for-the-terminal but readable rendering beats a vertically shredded one.
		if usableWidth := width - indentCols; usableWidth > getopt.HelpColumn {
			return usableWidth
		}
	}
	if len(fallback) > 0 && fallback[0] > 0 {
		return fallback[0]
	}
	return defaultTerminalWidth
}
