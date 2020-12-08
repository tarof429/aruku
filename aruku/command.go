package aruku

import (
	"fmt"
	"os/exec"
)

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
	fmt.Printf("Name: %v Description: %v Exit Status: %v\n", c.Name, c.Description, c.exitStatus)
	fmt.Println(c.output)
}

func (c *Command) GetExitStatus() int {
	return c.exitStatus
}

func (c *Command) GetOutput() string {
	return c.output
}
