package aruku

import (
	"fmt"
	"os/exec"
)

// CommandStatus indicates the status of the command
type CommandStatus string

// Command represents a script, it's description, and status
type Command struct {
	Name             string   `json:"name"`
	Args             []string `json:"args"`
	WorkingDirectory string   `json:"workingDirectory"`
	Description      string   `json:"description"`
	exitStatus       int
	output           string
}

// Run runs the command
func (c *Command) Run() {
	fmt.Printf("Running %v\n", c.Description)
	cmd := exec.Command(c.Name, c.Args...)
	cmd.Dir = c.WorkingDirectory
	combinedOutput, combinedOutputErr := cmd.CombinedOutput()

	if combinedOutputErr != nil {
		c.exitStatus = -1
	} else {
		c.exitStatus = 0
	}

	c.output = string(combinedOutput)
}

// Print prints the command
func (c *Command) Print() {
	fmt.Print(c)
}
