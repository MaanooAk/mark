package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecuteCommand expands the given template and runs it over the given list of items
func ExecuteCommand(template string, unquotedList *[]string) {

	list := quoteList(unquotedList)

	hasSinge := strings.Contains(template, "{}")
	hasMulti := strings.Contains(template, "{ }")

	if hasMulti {
		template = strings.ReplaceAll(template, "{ }", strings.Join(*list, " "))
	} else if !hasSinge && !strings.HasSuffix(template, ";") {
		// no single or multi in the template
		template = template + " {}"
	}

	if hasMulti && !hasSinge {
		// run once
		startCommand(template)

	} else {
		// run for every item
		for _, i := range *list {
			startCommand(strings.ReplaceAll(template, "{}", i))
		}
	}

}

func quoteList(list *[]string) *[]string {

	quotedList := []string{} // TODO optimize
	for _, i := range *list {
		if (i[0] == '"' && i[len(i)-1] == '"') || !strings.Contains(i, " ") {
			quotedList = append(quotedList, i)
		} else {
			quotedList = append(quotedList, "\""+i+"\"")
		}
	}

	return &quotedList
}

func startCommand(commandText string) error {
	if OptionVerbose {
		fmt.Println(commandText)
	}

	command := exec.Command("sh", "-c", commandText)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	command.Start()
	err := command.Wait()

	return err // TODO handle the status code returned
}
