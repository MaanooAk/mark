package main

import (
	"bufio"
	"errors"
)

// Store is a collection of all the marked items with meta data
type Store struct {
	List    []string
	Changed bool
}

// Operation is an operation over the collection
type Operation int

const (
	// Add Add if missing
	Add Operation = iota
	// Remove Remove if exists
	Remove
	// Xor Add if missing or remove if exists
	Xor
)

func (store *Store) has(item string) (exists bool, index int) {

	for index, i := range store.List {
		if i == item {
			return true, index
		}
	}
	return false, 0
}

// Append appends the given item to the end of the list with no checks
func (store *Store) append(item string) {

	store.List = append(store.List, item)
	store.Changed = true
}

// Remove removes the given index
func (store *Store) remove(index int) {

	store.List = append(store.List[:index], store.List[index+1:]...) // TODO check
	store.Changed = true
}

// Clear removes all the items from the list
func (store *Store) Clear() {
	if store.List == nil || len(store.List) == 0 {
		return
	}

	store.List = nil // TODO check
	store.Changed = true
}

// Add adds a single item to the stores list
func (store *Store) Add(item string, op Operation) {

	exists, index := store.has(item)

	if op == Add {

		if !exists {
			store.append(item)
		}

	} else if op == Remove {

		if exists {
			store.remove(index)
		}

	} else if op == Xor {

		if exists {
			store.remove(index)
		} else {
			store.append(item)
		}

	} else {
		panic(errors.New("Unknown operation"))
	}

}

func (store *Store) Read(scanner *bufio.Scanner, op Operation, transform func(string) string) {
	for scanner.Scan() {
		item := scanner.Text()
		if item == "" {
			continue
		}
		if transform != nil {
			item = transform(item)
		}
		store.Add(item, op)
	}
}

func (store *Store) Write(writer *bufio.Writer) {
	for _, i := range store.List {
		writer.WriteString(i)
		writer.WriteRune('\n')
	}
	writer.Flush()
}
