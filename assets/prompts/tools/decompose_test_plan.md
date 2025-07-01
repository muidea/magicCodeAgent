### decompose_test_plan

**描述：**  

在测试验证阶段，基于设计阶段输出的设计文档，当前的代码实现和需要解决的错误问题，制定当前阶段测试验证目标。

**注意：**  

- `decompose_test_plan`不允许批量执行，只能单次执行。
- 任何时候开始测试验证，都必须先使用`decompose_test_plan`来制定代码开发目标。

**用法示例：**  

```xml
<decompose_test_plan>
</decompose_test_plan>
```
