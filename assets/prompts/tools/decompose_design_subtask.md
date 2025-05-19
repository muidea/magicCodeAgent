### decompose_design_subtask

Description: 在功能设计阶段，需要根据当前的任务设计层级和背景描述，逐步分解出对应的子任务，在完成前一个分解任务后，可以使用该 tool 来获取子任务分解计划的相关要求。

根据当前任务的背景描述，该 tool 会提供子任务分解计划要求，用于指导继续分解对应的子任务还是处理下一个同层级的兄弟任务。

只有在所有的子任务和兄弟任务都分解完成后，才结束所有的任务，并使用attempt_completion工具结束设计工作。


任务层级(level)定义如下：
- project: 项目级，当前正在对项目目标进行分析和分解，需要根据项目需求和背景分解出对应各个子业务模块。
- module: 模块级，当前正在对指定的业务模块进行分析和分解，需要根据模块功能结合项目需求，分解出该业务模块包含的功能块。
- block: 功能块级，当前正在对指定的功能块进行分析和分解，需要根据功能块结合业务模块功能、项目需求，分解出该功能块需要包含的功能点和实现该功能点对应函数，并根据项目目录规划，定义函数实现代码文件和及其文件目录。
- function: 函数级，根据函数功能规划，实现代码功能。

**补充说明：函数级任务对应的是代码开发阶段任务，其他层级对应设计阶段任务**

Parameters:
- name: (required) 当前任务名
- level: (required) 当前任务层级
- description: (required) 当前任务的描述


Usage:
<decompose_design_subtask>
<name>当前任务名</name>
<level>taskLevel</level>
<description>基于当前任务背景描述</description>
</decompose_design_subtask>

Example: 

<decompose_design_subtask>
<name>人员信息管理项目功能设计</name>
<level>project</level>
<description>当前正在进行人员信息管理项目功能设计</description>
</decompose_design_subtask>