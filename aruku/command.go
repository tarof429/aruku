package aruku

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type CommandType string

const (
	ExecuteCommandType CommandType = "execute"
	ReadCommandType    CommandType = "read"
)

// Command represents a script, it's description, and status
type Command struct {
	Name             string   `json:"name"`
	Args             []string `json:"args"`
	WorkingDirectory string   `json:"workingDirectory"`
	Description      string   `json:"description"`
	exitStatus       int
	output           string
	CommandType      `json:"type"`
	VariableName     string `json:"variable"`
}

// Run runs the command
func (c *Command) Run(vars chan VariableMap, variables []VariableMap) bool {

	if c.CommandType == ReadCommandType {
		reader := bufio.NewReader(os.Stdin)

		var input string

		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if len(input) == 0 {
			return false
		}

		go func() {
			vars <- VariableMap{c.VariableName, strings.TrimSuffix(input, "\n")}
		}()
		return true

	} else {
		var replacedArgs []string

		for _, arg := range c.Args {
			if strings.HasPrefix(arg, "$") {
				for _, variable := range variables {
					if variable.key == arg[1:] {
						replacedArgs = append(replacedArgs, variable.value)
						break
					}
				}
			} else {
				replacedArgs = append(replacedArgs, arg)
			}
		}

		cmd := exec.Command(c.Name, replacedArgs...)
		cmd.Dir = c.WorkingDirectory
		combinedOutput, combinedOutputErr := cmd.CombinedOutput()

		if combinedOutputErr != nil {
			c.exitStatus = -1
		} else {
			c.exitStatus = 0
		}

		c.output = string(combinedOutput)
		return c.exitStatus == 0
	}

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
