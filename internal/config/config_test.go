package config

import (
	"path/filepath"
	"testing"

	cd "github.com/muidea/magicCommon/def"

	"github.com/muidea/magicCodeAgent/pkg/common"
)

func TestVerifyModuleWorkPath(t *testing.T) {
	// 获取当前系统的路径分隔符
	sep := string(filepath.Separator)

	// 定义跨平台测试用例
	tests := []struct {
		name            string
		input           string
		wantProjectName string
		wantPrefixPath  string
		wantError       *cd.Error
	}{
		{
			name:      "empty path",
			input:     "",
			wantError: cd.NewError(cd.IllegalParam, ""),
		},
		{
			name:      "relative path",
			input:     "internal/modules",
			wantError: cd.NewError(cd.IllegalParam, ""),
		},
		{
			name:            "valid modules path",
			input:           filepath.Join(sep+"abs", "internal", "modules"),
			wantProjectName: "abs",
			wantPrefixPath:  "internal" + sep + "modules",
			wantError:       nil,
		},
		{
			name:            "valid blocks path",
			input:           filepath.Join(sep, "a", "b", "internal", "modules", "blocks"),
			wantProjectName: "b",
			wantPrefixPath:  "internal" + sep + "modules" + sep + "blocks",
			wantError:       nil,
		},
		{
			name:            "valid kernel path",
			input:           filepath.Join(sep+"x", "y", "z", "internal", "modules", "kernel"),
			wantProjectName: "z",
			wantPrefixPath:  "internal" + sep + "modules" + sep + "kernel",
			wantError:       nil,
		},
		{
			name:      "invalid parent path",
			input:     filepath.Join(sep+"invalid", "path"),
			wantError: cd.NewError(cd.IllegalParam, ""),
		},
		{
			name:            "normalized path with slashes",
			input:           sep + filepath.Join("with", "extra", "slashes", "") + sep,
			wantProjectName: sep + filepath.Join("with", "extra") + sep,
			wantError:       cd.NewError(cd.IllegalParam, ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getPrefix, gotErr := VerifyModuleWorkPath(common.Service, tt.input)

			// 验证返回路径
			if tt.wantError == nil {
				if getPrefix != tt.wantPrefixPath {
					t.Errorf("unexpected result: getPrefix=%v, gotErr=%v", getPrefix, gotErr)
				}
				if gotErr != nil {
					t.Errorf("unexpected error: %v", gotErr)
				}
			} else {
				// 验证错误代码
				if gotErr == nil {
					t.Fatal("expected error, got nil")
				}
				if gotErr.Code != tt.wantError.Code {
					t.Errorf("expected error code: %v, got %v", tt.wantError.Code, gotErr.Code)
				}
			}
		})
	}
}
