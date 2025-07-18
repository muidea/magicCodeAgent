package biz

import (
	"fmt"
	"os"
	"path/filepath"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicCodeAgent/internal/config"
	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

func (s *Workspace) filterWorkspaceFile(ev event.Event, re event.Result) {
	val := ev.Data()
	if val == nil {
		log.Errorf("filterWorkspaceFile failed, invalid event value")
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal event value"))
		}
		return
	}
	workspaceName, workspaceOK := val.(string)
	if !workspaceOK {
		log.Errorf("filterWorkspaceFile failed, invalid event value:%v", val)
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal event value"))
		}
		return
	}
	fileList, fileErr := s.FilterWorkspaceFileList(workspaceName, nil)
	if re != nil {
		re.Set(fileList, fileErr)
	}
}

func (s *Workspace) FilterWorkspaceFileList(workspaceName string, filter fu.Filter) (ret []string, err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(workspaceName)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", workspaceName))
		return
	}
	if !curWorkspace.(*Station).Available() {
		err = cd.NewError(cd.InvalidOperation, fmt.Sprintf("workspace:%s not available", workspaceName))
		log.Errorf("workspace:%s not available", workspaceName)
		return
	}

	if !fp.Exist(curWorkspace.(*Station).workspaceConfig.WorkspaceDir) {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s, path:%s not exists", workspaceName, curWorkspace.(*Station).workspaceConfig.WorkspaceDir))
		return
	}

	ret = curWorkspace.(*Station).Filter(filter)
	return
}

func (s *Workspace) FilterWorkspaceFileTree(workspaceName, basePath string) (ret []*imkwc.DataNode, err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(workspaceName)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", workspaceName))
		return
	}
	if !curWorkspace.(*Station).Available() {
		err = cd.NewError(cd.InvalidOperation, fmt.Sprintf("workspace:%s not available", workspaceName))
		log.Errorf("workspace:%s not available", workspaceName)
		return
	}

	workingRootPath := curWorkspace.(*Station).workspaceConfig.GetWorkingPath()
	filter := NewBasePathFilter(workingRootPath, basePath)

	fileList := curWorkspace.(*Station).Filter(filter)
	return ConvertPathsToTree(workingRootPath, fileList), nil
}

func (s *Workspace) QueryFileContent(workspaceName string, filePath string) (ret *imkwc.FileContent, err *cd.Error) {
	workspacePtr := config.GetWorkspace(workspaceName)
	if workspacePtr == nil {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", workspaceName))
		log.Errorf("QueryFileContent failed, workspace not found:%s", workspaceName)
		return
	}
	if !workspacePtr.Available() {
		err = cd.NewError(cd.InvalidOperation, fmt.Sprintf("workspace:%s not available", workspaceName))
		log.Errorf("QueryFileContent failed, workspace not available:%s", workspaceName)
		return
	}

	fullFilePath := filepath.Join(workspacePtr.WorkspaceDir, filePath)
	if !fp.Exist(fullFilePath) {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s, path:%s not exists", workspaceName, fullFilePath))
		log.Errorf("QueryFileContent failed, file not exists:%s", fullFilePath)
		return
	}

	content, readErr := os.ReadFile(fullFilePath)
	if readErr != nil {
		err = cd.NewError(cd.Unexpected, readErr.Error())
		log.Errorf("QueryFileContent failed, read file error:%s, err:%v", fullFilePath, readErr)
		return
	}

	fileExt := filepath.Ext(fullFilePath)
	if len(fileExt) > 0 && fileExt[0] == '.' {
		fileExt = fileExt[1:]
	}

	ret = &imkwc.FileContent{
		Type:    fileExt,
		Content: string(content),
	}
	return
}
