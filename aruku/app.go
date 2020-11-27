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
	Author              string
	Description         string
	CmdList             []CmdList
	currentCmdList      CmdList
	currentCmdListIndex int
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

// SetCmdList sets which CmdList we want to run based on the description.
// If a list was found, return true, otherwise false
func (a *App) SetCmdList(description string) bool {

	for _, cmdList := range a.CmdList {
		if cmdList.Description == description {
			a.currentCmdList = cmdList
			a.currentCmdListIndex = 0
			return true
		}
	}
	return false
}

// HasNextCmd returns true if there is a command to run.
func (a *App) HasNextCmd() bool {
	if a.currentCmdListIndex < len(a.currentCmdList.Cmds) {
		return true
	}
	return false
}

// GetCurrentCmd returns a copy of the current command
func (a *App) GetCurrentCmd() Command {
	return a.currentCmdList.Cmds[a.currentCmdListIndex]
}

// RunCurrentCmd runs the current command
func (a *App) RunCurrentCmd() {
	a.currentCmdList.Cmds[a.currentCmdListIndex].Run()
}

// PointToNextCmd moves the index to the next command (if available)
func (a *App) PointToNextCmd() {
	if a.currentCmdListIndex < len(a.currentCmdList.Cmds) {
		a.currentCmdListIndex = a.currentCmdListIndex + 1
	}
}
