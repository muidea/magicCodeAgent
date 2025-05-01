package biz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwitchMode(t *testing.T) {
	// 测试用例1：有理由的模式切换
	modeSlug1 := "code"
	reason1 := "Need to make code changes"
	expectedResult1 := "Switching to mode 'code' because: Need to make code changes"

	result1, err := SwitchMode(modeSlug1, reason1)

	// 验证结果
	assert.NoError(t, err, "模式切换应该成功")
	assert.Equal(t, expectedResult1, result1, "返回结果应该符合预期")

	// 测试用例2：无理由的模式切换
	modeSlug2 := "ask"
	reason2 := ""
	expectedResult2 := "Switching to mode 'ask'"

	result2, err := SwitchMode(modeSlug2, reason2)

	// 验证结果
	assert.NoError(t, err, "无理由的模式切换应该成功")
	assert.Equal(t, expectedResult2, result2, "返回结果应该符合预期")

	// 测试用例3：空模式名
	modeSlug3 := ""
	reason3 := "Testing empty mode"
	expectedResult3 := "Switching to mode '' because: Testing empty mode"

	result3, err := SwitchMode(modeSlug3, reason3)

	// 验证结果
	assert.NoError(t, err, "空模式名应该不会导致错误")
	assert.Equal(t, expectedResult3, result3, "返回结果应该符合预期")

	// 测试用例4：特殊字符
	modeSlug4 := "special!@#"
	reason4 := "Special chars: !@#$%^&*()"
	expectedResult4 := "Switching to mode 'special!@#' because: Special chars: !@#$%^&*()"

	result4, err := SwitchMode(modeSlug4, reason4)

	// 验证结果
	assert.NoError(t, err, "特殊字符应该不会导致错误")
	assert.Equal(t, expectedResult4, result4, "返回结果应该符合预期")
}
