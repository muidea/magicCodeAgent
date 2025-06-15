### decompose_code_plan

**描述：**  

在代码开发阶段，基于设计阶段输出的设计文档，制定当前阶段代码开发目标。

**注意：**  

- `decompose_code_plan`不允许批量执行，只能单次执行。
- 任何时候开始代码开发，都必须先使用`decompose_code_plan`来制定代码开发目标。

**功能层级 (level) 定义：**  

- **project：** 项目级，根据项目功能设计文档，完成项目代码目录创建和初始化。
- **module：** 业务模块级，根据业务模块设计文档，完成业务模块代码代码开发和。
- **block：** 功能模块级，根据功能模块设计文档，完成功能模块代码代码开发和。

**参数：**  

- **name** (必填)：待开发功能名
- **level** (必填)：待开发的功能层级
- **path** (必填)：待开发功能对应的设计文档路径(相对对项目根目录的路径)

**参数说明：**

- `level` 为 `project` 时，`name` 当前项目名, `description` 为项目功能描述。
- `level` 为 `module` 时，`name` 当前业务模块名, `description` 为业务模块功能描述。
- `level` 为 `block` 时，`name` 当前功能模块名, `description` 为功能模块功能描述。

- `level` 为 `project` 时，`path`值为`./{项目名}-项目功能.md`
- `level` 为 `module` 时，`path`值为`./{项目名}/{业务模块名}-业务模块.md`
- `level` 为 `block` 时，`path`值为`./{项目名}/{业务模块名}/{功能模块}-功能模块.md`

**用法示例：**  

```xml
<decompose_code_plan>
    <name>产品新增</name>
    <level>block</level>
    <path>./产品管理/产品新增-业务模块.md</path>
    <description>提供产品信息新增功能</description>
</decompose_code_plan>
```
