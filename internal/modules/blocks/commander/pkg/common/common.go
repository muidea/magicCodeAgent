package common

const CommanderModule = "/internal/modules/blocks/commander"

const (
	CommandRun        = "/command/run"
	CommandTerminate  = "/command/terminate"
	CommandStatus     = "/command/status"
	ConsoleLogsAppend = "/console/logs"
)

type CommanderParam struct {
	WorkspaceDir string   `json:"workspaceDir"`
	Command      string   `json:"command"`
	Args         []string `json:"args"`
	Timeout      int      `json:"timeout"`
}

type CommanderResult struct {
	Workspace      string `json:"workspace"`
	StatusMessage  string `json:"statusMessage"`
	ErrorMessage   string `json:"errorMessage"`
	ConsoleMessage string `json:"consoleMessage"`
}
