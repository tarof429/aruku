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
	//variables        []VariableMap
}

// Run runs the command
func (c *Command) Run(vars chan VariableMap, variables []VariableMap) {

	if c.CommandType == ReadCommandType {
		reader := bufio.NewReader(os.Stdin)

		input, _ := reader.ReadString('\n')

		go func() {
			vars <- VariableMap{c.VariableName, strings.TrimSuffix(input, "\n")}
		}()

	} else {
		var replacedArgs []string

		// for _, variable := range variables {
		// 	fmt.Printf("Variable: %v: %v\n", variable.key, variable.value)
		// }

		for _, arg := range c.Args {
			//fmt.Printf("Evaluating %v\n", arg)
			if strings.HasPrefix(arg, "$") {
				//fmt.Println("It's a variable")
				for _, variable := range variables {
					//fmt.Printf("Checking if %v == %v\n", variable.key, arg[1:])
					if variable.key == arg[1:] {
						replacedArgs = append(replacedArgs, variable.value)
						break
					}
				}
			} else {
				replacedArgs = append(replacedArgs, arg)
			}
		}
		//fmt.Printf("Replaced args: %v\n", replacedArgs)
		cmd := exec.Command(c.Name, replacedArgs...)
		cmd.Dir = c.WorkingDirectory
		combinedOutput, combinedOutputErr := cmd.CombinedOutput()

		if combinedOutputErr != nil {
			c.exitStatus = -1
		} else {
			c.exitStatus = 0
		}

		c.output = string(combinedOutput)
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
