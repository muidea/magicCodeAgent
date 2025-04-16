package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tool/biz"
	engine "github.com/muidea/magicEngine/http"
)

type Tool struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Tool
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Tool) *Tool {
	ptr := &Tool{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Tool) RegisterRoute() {
	// 注册读取文件路由
	s.routeRegistry.SetRoute("/tool/read-file", s.readFileHandler)

	// 注册搜索文件路由
	s.routeRegistry.SetRoute("/tool/search-files", s.searchFilesHandler)

	// 注册列出文件路由
	s.routeRegistry.SetRoute("/tool/list-files", s.listFilesHandler)

	// 注册应用差异路由
	s.routeRegistry.SetRoute("/tool/apply-diff", s.applyDiffHandler)

	// 注册写入文件路由
	s.routeRegistry.SetRoute("/tool/write-to-file", s.writeToFileHandler)

	// 注册插入内容路由
	s.routeRegistry.SetRoute("/tool/insert-content", s.insertContentHandler)

	// 注册搜索替换路由
	s.routeRegistry.SetRoute("/tool/search-and-replace", s.searchAndReplaceHandler)

	// 注册执行命令路由
	s.routeRegistry.SetRoute("/tool/execute-command", s.executeCommandHandler)

	// 注册获取指令路由
	s.routeRegistry.SetRoute("/tool/fetch-instructions", s.fetchInstructionsHandler)

	// 注册使用MCP工具路由
	s.routeRegistry.SetRoute("/tool/use-mcp-tool", s.useMCPToolHandler)

	// 注册访问MCP资源路由
	s.routeRegistry.SetRoute("/tool/access-mcp-resource", s.accessMCPResourceHandler)
}

// 读取文件处理函数
func (s *Tool) readFileHandler(ctx engine.RequestContext) {
	path := ctx.Request().FormValue("path")
	startLineStr := ctx.Request().FormValue("start_line")
	endLineStr := ctx.Request().FormValue("end_line")

	startLine := 1
	endLine := 0

	if startLineStr != "" {
		val, err := strconv.Atoi(startLineStr)
		if err == nil && val > 0 {
			startLine = val
		}
	}

	if endLineStr != "" {
		val, err := strconv.Atoi(endLineStr)
		if err == nil && val > 0 {
			endLine = val
		}
	}

	result, err := s.bizPtr.ReadFile(path, startLine, endLine)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte(result))
}

// 搜索文件处理函数
func (s *Tool) searchFilesHandler(ctx engine.RequestContext) {
	path := ctx.Request().FormValue("path")
	regex := ctx.Request().FormValue("regex")
	filePattern := ctx.Request().FormValue("file_pattern")

	result, err := s.bizPtr.SearchFiles(path, regex, filePattern)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte(result))
}

// 列出文件处理函数
func (s *Tool) listFilesHandler(ctx engine.RequestContext) {
	path := ctx.Request().FormValue("path")
	recursiveStr := ctx.Request().FormValue("recursive")

	recursive := false
	if recursiveStr == "true" {
		recursive = true
	}

	result, err := s.bizPtr.ListFiles(path, recursive)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte(result))
}

// 应用差异处理函数
func (s *Tool) applyDiffHandler(ctx engine.RequestContext) {
	path := ctx.Request().FormValue("path")
	diff := ctx.Request().FormValue("diff")

	err := s.bizPtr.ApplyDiff(path, diff)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte("Diff applied successfully"))
}

// 写入文件处理函数
func (s *Tool) writeToFileHandler(ctx engine.RequestContext) {
	path := ctx.Request().FormValue("path")
	content := ctx.Request().FormValue("content")

	err := s.bizPtr.WriteToFile(path, content)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte("File written successfully"))
}

// 插入内容处理函数
func (s *Tool) insertContentHandler(ctx engine.RequestContext) {
	path := ctx.Request().FormValue("path")
	operationsStr := ctx.Request().FormValue("operations")

	var operations []biz.InsertOperation
	err := json.Unmarshal([]byte(operationsStr), &operations)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusBadRequest)
		ctx.Response().Write([]byte(fmt.Sprintf("Invalid operations format: %v", err)))
		return
	}

	err = s.bizPtr.InsertContent(path, operations)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte("Content inserted successfully"))
}

// 搜索替换处理函数
func (s *Tool) searchAndReplaceHandler(ctx engine.RequestContext) {
	path := ctx.Request().FormValue("path")
	operationsStr := ctx.Request().FormValue("operations")

	var operations []biz.SearchReplaceOperation
	err := json.Unmarshal([]byte(operationsStr), &operations)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusBadRequest)
		ctx.Response().Write([]byte(fmt.Sprintf("Invalid operations format: %v", err)))
		return
	}

	err = s.bizPtr.SearchAndReplace(path, operations)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte("Search and replace completed successfully"))
}

// 执行命令处理函数
func (s *Tool) executeCommandHandler(ctx engine.RequestContext) {
	command := ctx.Request().FormValue("command")
	cwd := ctx.Request().FormValue("cwd")

	result, err := s.bizPtr.ExecuteCommand(command, cwd)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte(result))
}

// 获取指令处理函数
func (s *Tool) fetchInstructionsHandler(ctx engine.RequestContext) {
	task := ctx.Request().FormValue("task")

	result, err := s.bizPtr.FetchInstructions(task)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte(result))
}

// 使用MCP工具处理函数
func (s *Tool) useMCPToolHandler(ctx engine.RequestContext) {
	serverName := ctx.Request().FormValue("server_name")
	toolName := ctx.Request().FormValue("tool_name")
	argumentsStr := ctx.Request().FormValue("arguments")

	var arguments json.RawMessage
	err := json.Unmarshal([]byte(argumentsStr), &arguments)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusBadRequest)
		ctx.Response().Write([]byte(fmt.Sprintf("Invalid arguments format: %v", err)))
		return
	}

	result, err := s.bizPtr.UseMCPTool(serverName, toolName, arguments)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte(result))
}

// 访问MCP资源处理函数
func (s *Tool) accessMCPResourceHandler(ctx engine.RequestContext) {
	serverName := ctx.Request().FormValue("server_name")
	uri := ctx.Request().FormValue("uri")

	result, err := s.bizPtr.AccessMCPResource(serverName, uri)
	if err != nil {
		ctx.Response().WriteHeader(http.StatusInternalServerError)
		ctx.Response().Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	ctx.Response().Write([]byte(result))
}
