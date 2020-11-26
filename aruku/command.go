package aruku

// CommandStatus indicates the status of the command
type CommandStatus string

// Command represents a script, it's description, and status
type Command struct {
	Args             []string `json:"args"`
	WorkingDirectory string   `json:"workingDirectory"`
	Description      string   `json:"description"`
	exitStatus       int
	output           string
}
