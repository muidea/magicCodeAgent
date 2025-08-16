### decompose_debug_task

**描述：**  

在测试验证阶段，制定测试验证目标。

**注意：**  

- `decompose_debug_task`不允许批量执行，只能单次执行。
- 任何时候开始测试验证，都必须先使用`decompose_debug_task`来制定测试验证目标。
- 在制定新测试验证任务前，务必先使用`attempt_finish_task`提交前一个任务确认。

**参数：**  

- **name** (必填)：待测试验证功能名
- **level** (必填)：固定值`project`
- **path** (必填)：待测试验证功能对应路径(基于根目录的相对路径)
- **description** (必填)：待测试验证功能描述

**用法示例：**  

1. 修复Product项目代码错误

```xml
<decompose_debug_task>
    <name>Product</name>
    <level>project</level>
    <path>.</path>
    <description>修复产品管理项目代码错误</description>
</decompose_debug_task>
```
