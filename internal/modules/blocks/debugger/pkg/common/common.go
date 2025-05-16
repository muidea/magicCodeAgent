package common

const DebuggerModule = "/internal/modules/blocks/debugger"

const (
	DebugCommandRun       = "/command/debug/run"
	DebugCommandTerminate = "/command/debug/terminate"
	DebugCommandStatus    = "/command/debug/status"
)

type DebugCommand struct {
	WorkspaceDir string   `json:"workspaceDir"`
	Command      string   `json:"command"`
	Args         []string `json:"args"`
}

type DebugResult struct {
	Workspace string `json:"workspace"`
	Status    string `json:"status"`
}
