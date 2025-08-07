### decompose_code_task

**描述：**  

在代码开发阶段，基于设计阶段输出的设计文档，制定当前阶段代码开发任务。

**注意：**  

- `decompose_code_task`不允许批量执行，只能单次执行。
- 任何时候开始代码开发，都必须先使用`decompose_code_task`来制定代码开发任务。
- 在制定新开发任务前，务必先使用`attempt_finish_task`提交前一个任务执行完成确认。

**功能层级 (level) 定义：**  

- **project：** 项目级，根据设计阶段输出的业务模块清单完成业务模块代码初始化。
- **module：** 业务模块级，根据设计阶段输出功能模块设计文档，完成当前业务模块代码开发。
- **block：** 功能模块级，根据设计阶段输出功能模块设计文档，核对当前代码实现是否符合设计要求。

**参数：**  

- **name** (必填)：待开发功能名
- **level** (必填)：待开发的功能层级
- **path** (必填)：待开发功能对应的设计文档路径(相对对项目根目录的路径)

**参数说明：**

- `level` 为 `project` 时，`name` 待开发项目名, `description` 为项目功能描述。
- `level` 为 `module` 时，`name` 待开发业务模块名, `description` 为业务模块功能描述。
- `level` 为 `block` 时，`name` 待开发功能模块名, `description` 为功能模块功能描述。

- `level` 为 `project` 时，`path`值为`./{项目名}-项目功能.md`
- `level` 为 `module` 时，`path`值为`./{项目名}/{业务模块名}-业务模块.md`
- `level` 为 `block` 时，`path`值为`./{项目名}/{业务模块名}/{功能模块名}-功能模块.md`

**用法示例：**  

1. 制定产品管理对应的功能模块级开发任务

```xml
<decompose_code_task>
    <name>ProductEdit</name>
    <level>block</level>
    <path>./demo/Product/ProductEdit-功能模块.md</path>
    <description>提供产品信息修改功能，包括基本信息、分类和图片的更新</description>
</decompose_code_task>
```

2. 制定产品展示系统对应的项目功能级开发任务

```xml
<decompose_code_task>
    <name>demoProject</name>
    <level>project</level>
    <path>./demoProject-项目功能.md</path>
    <description>产品展示系统，主要实现产品信息、分类信息和图片的管理功能，支持用户对产品的登记和发布操作</description>
</decompose_code_task>
```
