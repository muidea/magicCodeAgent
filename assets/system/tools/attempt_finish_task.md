### attempt_finish_task

**描述：**

在完成设计目标，开发目标或验证目标对应的任务，并核对工作输出满足要求后，使用此工具进行任务完成确认。

**注意：**  

1. 在确认当前执行任务的目标完成，并符合要求前，不得使用该工具。在使用此工具前，您必须在 <thinking></thinking> 标记中明确已经完成的工作项，并明确给出目标和要求的核对结果。
2. 任务完成确认必须与`decompose_design_task`、`decompose_code_task`、`decompose_debug_task`工具定义的任务一对一对应。


**参数：**  

- **name** (必填)：正在开展工作对应的名称(项目名，业务模块名，功能模块名)
- **level** (必填)：正在进行的工作对应的层级(project, module, block)
- **path** (必填)：正在开展工作对应的文档路径
- **needContinue** (必填)：是否继续进行下一个任务。true表示继续进行下一个任务，false表示所有任务完成。
- **result** (必填)： 任务结果。以最终结果的方式表述，不需要用户进一步输入。不要以提问或提供进一步帮助的方式结束结果。

>上述三个参数必须保持与`decompose_design_task`,`decompose_code_task`,`decompose_debug_task`中传递的参数对应

**参数说明：**

- `level` 为 `project` 时，`name` 当前项目名
- `level` 为 `module` 时，`name` 当前业务模块名
- `level` 为 `block` 时，`name` 当前功能模块名

- `level` 为 `project` 时，`path`值为`./{项目名}-项目功能.md`
- `level` 为 `module` 时，`path`值为`./{项目名}/{业务模块名}-业务模块.md`
- `level` 为 `block` 时，`path`值为`./{项目名}/{业务模块名}/{功能模块名}-功能模块.md`

- `needContinue` 如果当前level存在下级层级，并且下层级存在未完成项，则返回true，否则返回false

**用法示例：**  

1. 完成分类信息查询功能
```xml
<attempt_finish_task>
<name>CategoryQuery</name>
<level>block</level>
<path>./demoProduct/Category/CategoryQuery-功能模块.md</path>
<needContinue>false</needContinue>
<result>
已完成CategoryQuery功能模块开发，包含以下内容：
1. 实现了分类信息查询功能，支持模糊查询和分页
2. 代码实现与设计文档严格一致
3. 使用ProTable组件实现了列表展示和查询表单
4. 已集成到Category业务模块的路由配置中
5. 代码通过静态检查，无错误
</result>
</attempt_finish_task>
```

