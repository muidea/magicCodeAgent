### decompose_design_task

**描述：**  

在功能设计阶段，基于根据当前的设计要求，制定当前阶段代码设计任务。

**注意：**  

- `decompose_design_task`不允许批量执行，只能单次执行。
- 任何时候开始设计，都必须先使用`decompose_design_task`来制定功能设计任务。
- 在制定新设计任务前，务必先使用`attempt_finish_task`提交前一个任务确认。

**功能层级 (level) 定义：**  

- **project：** 项目级，当前正在根据项目需求设计项目功能点，并输出项目功能设计文档
- **module：** 业务模块级，当前正在根据项目功能设计文档，设计说包含的业务模块，并输出各个业务模块设计文档。
- **block：** 功能模块级，当前正在根据业务模块设计文档，设计包含的功能模块，并输出各个功能模块设计文档。

**参数：**  

- **name** (必填)：待设计功能英文名
- **level** (必填)：当前设计的功能层级
- **path** (可选)：待设计功能所归属的父级功能的设计文档路径(相对对项目根目录的路径)
- **description** (必填)：待设计功能描述

**参数说明：**

- `level` 为 `project` 时，`name` 待设计项目英文名。
- `level` 为 `module` 时，`name` 待设计业务模块英文名。
- `level` 为 `block` 时，`name` 待设计功能模块英文名。

- `level` 为 `project` 时，`description` 表示当前项目的需求描述。  
- `level` 为 `module` 时，`description` 表示业务模块的功能描述。  
- `level` 为 `block` 时，`description` 表示功能模块的功能描述。  

- `level` 为 `project` 时，`path`值为空
- `level` 为 `module` 时，`path`值为`./{项目名}-项目功能.md`
- `level` 为 `block` 时，`path`值为`./{项目名}/{业务模块名}-业务模块.md`

>强调：name必须使用待设计功能的英文名**

**用法示例：**  

1. 制定产品信息新增设计任务

```xml
<decompose_design_task>
    <name>ProductCreate</name>
    <level>block</level>
    <path>./demo/Product-业务模块.md</path>
    <description>提供产品信息新增功能</description>
</decompose_design_task>
```
