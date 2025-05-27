### decompose_design_plan

**描述：**  

在功能设计阶段，用于制定功能分解目标。根据当前功能名、当前功能设计文档、当前功能层级和功能描述，制定下一层级功能的分解目标和要求。


**注意：**  
    
- `decompose_design_plan`不允许批量执行，只能单次执行。
- 任何时候开始设计，都必须先使用`decompose_design_plan`来制定功能分解目标。

**功能层级 (level) 定义：**  
- **project：** 项目级，当前正在对项目目标进行功能分解，根据项目需求和背景分解出业务模块。
- **module：** 模块级，当前正在对指定的业务模块进行功能分解，根据业务模块功能分解出功能块。
- **block：** 功能块级，当前正在对指定的功能块进行功能分解，根据功能块分解出功能点。
- **point：** 功能点级，当前正在对指定的功能点进行功能分解，根据功能点分解出子功能点。

**参数：**  
- **name** (必填)：当前功能名(如果是level为project，则name为项目名)
- **path** (必填)：当前功能对应设计文档的路径(相对对项目根目录的路径)
- **level** (必填)：当前功能层级
- **description** (必填)：当前功能描述

**参数说明：**  
- 当 `level` 为 `project` 时，`name` 表示项目名，`description` 表示当前项目的需求描述。  
- 当 `level` 为 `module` 时，`name` 表示业务模块名，`description` 表示当前业务模块功能描述。  
- 当 `level` 为 `block` 时，`name` 表示功能块名，`description` 表示当前功能块功能描述。  
- 当 `level` 为 `point` 时，`name` 表示功能点名，`description` 表示当前功能点功能描述。  

名称根据待分析功能对应的 MD 文件名称生成，例如：  
- 文件路径为 `{业务模块名}/{功能块名}/{功能点名}-功能点设计.md`，提取的 `level` 为 `point`，`name` 为功能点名，`description` 为功能点功能描述。  
- 文件路径为 `{业务模块名}/{功能模块名}-功能模块设计.md`，提取的 `level` 为 `block`，`name` 为功能模块名，`description` 为功能模块功能描述。

**用法示例：**  

```xml
<decompose_design_plan>
    <name>产品详情展示</name>
    <path>./产品管理/产品详情展示</path>
    <level>block</level>
    <description>展示产品详细信息，包括基础信息、规格参数、图片展示等</description>
</decompose_design_plan>
