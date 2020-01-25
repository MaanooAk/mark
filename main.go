package main

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
)

// StoreFilename is the file name without the home path
const StoreFilename = ".mark"

func findStorePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return StoreFilename
	}
	return path.Join(home, StoreFilename)
}

func hasInputPipe() bool {
	if info, err := os.Stdin.Stat(); err == nil {
		return info.Mode()&os.ModeDevice == 0
	}
	return false
}

func normalizeItem(item string) string {
	if OptionRelative {
		return item
	}
	if path, err := filepath.Abs(item); err == nil {
		return path
	}
	return item
}

func main() {

	if !ParseArguments() {
		return
	}

	storePath := findStorePath()
	store := &Store{}

	if !OptionClear {
		loadStoreFile(storePath, store)
		store.Changed = false
	} else {
		store.Changed = true
	}

	if OptionReadStdin || hasInputPipe() {
		store.Read(bufio.NewScanner(os.Stdin), OptionOperation, normalizeItem)
	}

	for _, i := range ArgumentItems {
		store.Add(normalizeItem(i), OptionOperation)
	}

	if OptionWriteStdout || (len(os.Args) == 1 && !hasInputPipe()) {
		store.Write(bufio.NewWriter(os.Stdout))
	}

	for _, i := range ArgumentExecs {
		ExecuteCommand(i, &store.List)
	}

	storeStoreFile(storePath, store)

}

func loadStoreFile(storePath string, store *Store) {
	file, err := os.Open(storePath)
	if err != nil {
		return // there is no file to load
	}
	store.Read(bufio.NewScanner(file), Add, nil)
	file.Close()
}

func storeStoreFile(storePath string, store *Store) {
	if !store.Changed {
		return
	}
	file, err := os.Create(storePath)
	if err != nil {
		panic(err)
	}
	store.Write(bufio.NewWriter(file))
	file.Close()
}
