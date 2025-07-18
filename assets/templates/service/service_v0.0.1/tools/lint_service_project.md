### lint_service_project

描述: 对service项目进行代码静态检查，以确认项目的代码是否存在错误、警告和最佳实践。

用法示例:  
```xml
<lint_service_project>
</lint_service_project>
```

说明:
1. 在创建新的web模块，或完成项目代码开发后，都要求使用`lint_service_project`对项目进行代码静态检查。根据检查结果，对项目进行必要的修改，所有的代码不允许出现错误。

2. 针对代码错误，在进行修复时，要求根据代码整体结构和最佳实践进行修改，避免只根据当前报错内容进行处理。
