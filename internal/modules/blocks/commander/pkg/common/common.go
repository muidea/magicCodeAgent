package common

const CommanderModule = "/internal/modules/blocks/commander"

const (
	CommandRun       = "/command/run"
	CommandTerminate = "/command/terminate"
	CommandStatus    = "/command/status"
)

type CommanderParam struct {
	WorkspaceDir string   `json:"workspaceDir"`
	Command      string   `json:"command"`
	Args         []string `json:"args"`
	Restart      bool     `json:"restart"`
}

type CommanderResult struct {
	Workspace     string `json:"workspace"`
	StatusMessage string `json:"statusMessage"`
	ErrorMessage  string `json:"errorMessage"`
}
