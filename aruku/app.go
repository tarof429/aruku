package aruku

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
)

const (
	tmpDir        = "/tmp/aruku"
	arukuDataFile = "aruku.yaml"
)

var errorKeywords = []string{"usage", "inactive", "disabled", "dead", "error", "fail"}

// App represents this application
type App struct {
	Author               string
	Description          string
	CmdList              []CmdList
	currentCmdList       CmdList
	currentCmdListIndex  int
	previousCmdListIndex int
	variableMapChan      chan VariableMap
	variables            []VariableMap
}

type VariableMap struct {
	key   string
	value string
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

	fmt.Printf("Creating: %v\n", path)
	os.MkdirAll(path, os.FileMode(0755))

	mode := int(0644)

	updatedData, _ := json.MarshalIndent(a, "", "\t")

	err := ioutil.WriteFile(filepath.Join(path, arukuDataFile), updatedData, os.FileMode(mode))

	return err
}

// Load loads data from path
func (a *App) Load(path string) error {

	err := a.Read(filepath.Join(path, arukuDataFile))

	if err != nil {
		return err
	}
	return nil
}

// SetCmdList sets which CmdList we want to run based on the description.
// If a list was found, return true, otherwise false
func (a *App) SetCmdList(description string) bool {

	for _, cmdList := range a.CmdList {
		if cmdList.Description == description {
			a.currentCmdList = cmdList
			a.variableMapChan = make(chan VariableMap)

			go func() {
				for readVariable := range a.variableMapChan {
					a.variables = append(a.variables, readVariable)
				}

			}()

			a.currentCmdListIndex = 0
			return true
		}
	}
	return false
}

// HasCmd returns true if the index is less than the length of the command list
func (a *App) HasCmd() bool {
	return a.currentCmdListIndex < len(a.currentCmdList.Cmds)
}

// HasNextCmd returns true if there is a command to run.
func (a *App) HasNextCmd() bool {
	if a.currentCmdListIndex < len(a.currentCmdList.Cmds) {
		return true
	}
	return false
}

// TotalCmds returns the total number of commands
func (a *App) TotalCmds() int {
	return len(a.currentCmdList.Cmds)
}

// GetCurrentCmd returns a copy of the current command
func (a *App) GetCurrentCmd() Command {
	return a.currentCmdList.Cmds[a.currentCmdListIndex]
}

// RunCurrentCmd runs the current command
func (a *App) RunCurrentCmd() bool {
	return a.currentCmdList.Cmds[a.currentCmdListIndex].Run(a.variableMapChan, a.variables)
}

// RunPreviousCmd runs the current command
func (a *App) RunPreviousCmd() bool {
	return a.currentCmdList.Cmds[a.previousCmdListIndex].Run(a.variableMapChan, a.variables)
}

// PointToPreviousCmd moves currentCmdListIndex to the previous location
func (a *App) PointToPreviousCmd() {
	if a.currentCmdListIndex-1 >= 0 {
		a.currentCmdListIndex = a.currentCmdListIndex - 1
	}
}

// PointToNextCmd moves the index to the next command (if available)
func (a *App) PointToNextCmd() {
	if a.currentCmdListIndex < len(a.currentCmdList.Cmds) {
		a.currentCmdListIndex = a.currentCmdListIndex + 1
	}
}

// Run runs the current command
func (a *App) Run() bool {
	return a.RunCurrentCmd()
}

// ShowCurrentCommandOutput shows the command output
func (a *App) ShowCurrentCommandOutput() {
	cmd := a.GetCurrentCmd()

	pterm.Println()

	output := cmd.GetOutput()
	lowerOutput := strings.ToLower(output)

	errorFlag := false

	for _, keyword := range errorKeywords {
		index := strings.Index(lowerOutput, keyword)
		if index >= 0 {
			errorFlag = true
			break
		}
	}

	if errorFlag {
		pterm.Println(pterm.LightRed(output))
	} else {
		pterm.Println(output)
	}

}
