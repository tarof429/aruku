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

func (c *Command) readCommand(vars chan VariableMap) bool {
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
}

func ParseCommand(cmd string, args []string, variableMap []VariableMap) (string, []string) {

	var replacedArgs []string

	// Iterate through the array of args
	for _, arg := range args {
		// Split each arg by whitespace
		miniargs := strings.Split(arg, " ")
		var updatedArg string // the updated value of 'arg' with variables replaced
		// iterate through each word in the arg, doing variable replacement
		for _, miniarg := range miniargs {
			if strings.Index(miniarg, "$") == 0 {
				for _, vmap := range variableMap {
					if strings.Contains(miniarg, vmap.key) {
						updatedArg = updatedArg + " " + vmap.value
						if miniarg[len(miniarg)-1] == '"' {
							updatedArg = updatedArg + "\""
						}
						break
					}
				}
			} else {
				updatedArg = updatedArg + " " + miniarg
			}
			updatedArg = strings.Trim(updatedArg, " ")
		}
		// finally, append the updated arg to our list of args
		replacedArgs = append(replacedArgs, updatedArg)
	}
	return cmd, replacedArgs

}

func (c *Command) execCommand(variableMap []VariableMap) bool {

	name, replacedArgs := ParseCommand(c.Name, c.Args, variableMap)

	cmd := exec.Command(name, replacedArgs...)
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

// Run either sets a variable or executes a command
func (c *Command) Run(vars chan VariableMap, variableMap []VariableMap) bool {

	if c.CommandType == ReadCommandType {
		return c.readCommand(vars)
	}
	return c.execCommand(variableMap)
}

// Print prints the command
func (c *Command) Print() {
	fmt.Printf("Name: %v Description: %v Exit Status: %v\n", c.Name, c.Description, c.exitStatus)
	fmt.Println(c.output)
}

// GetOutput returns the output from the command
func (c *Command) GetOutput() string {
	return c.output
}
