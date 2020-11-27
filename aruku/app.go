package aruku

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mholt/archiver"
)

// App represents this application
type App struct {
	Author      string
	Description string
	CmdList     []CmdList
}

// CmdList is a list of commands
type CmdList struct {
	Description string    `json:"Description"`
	Cmds        []Command `json:"Commands"`
}

// Read reads the file and returns a list of Commands
func (a *App) Read(path string) error {

	data, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while reading %v: %v\n", path, err)
		return err
	}

	err = json.Unmarshal(data, &a)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while unmarshalling: %v\n", err)
	}

	return err
}

// Write creates the data file
func (a *App) Write(path string) error {

	mode := int(0644)

	updatedData, _ := json.MarshalIndent(a, "", "\t")

	err := ioutil.WriteFile(path, updatedData, os.FileMode(mode))

	return err
}

// Import unzips file containing list of commands and any scripts
// and reads the list of commands
func (a *App) Import(path, destination string) error {
	err := archiver.Unarchive(path, destination)

	if err != nil {
		return err
	}

	err = a.Read(destination)

	if err != nil {
		return err
	}
	return nil
}

// Export basically zips up the file containing the list of commands
// and any scripts
func (a *App) Export(path string, archiveName string) error {

	return archiver.Archive([]string{path}, archiveName)
}

func (a *App) getCmdListIndex(description string) int {
	cmdListIndex := -1

	for index, cmdList := range a.CmdList {
		if cmdList.Description == description {
			cmdListIndex = index
			break
		}
	}
	return cmdListIndex
}

// Run runs a cmdList with a matching description
func (a *App) Run(description string) {

	cmdListIndex := a.getCmdListIndex(description)

	if cmdListIndex > -1 {
		cmds := a.CmdList[cmdListIndex].Cmds

		for index := range cmds {
			cmds[index].Run()
		}
	}
}

// Print the commands
func (a *App) Print(description string) {

	cmdListIndex := a.getCmdListIndex(description)

	if cmdListIndex > -1 {
		cmds := &a.CmdList[cmdListIndex].Cmds

		for _, cmd := range *cmds {
			cmd.Print()
		}
	}
}
