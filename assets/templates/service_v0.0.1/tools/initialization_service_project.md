### initialization_service_project

描述：初始化后端服务项目，使用后端服务模板快速创建一个后端服务项目。

参数：
- **name**: (必填) 项目名称，表示初始化的后端服务项目的名称。
- **description**: (可选) 项目描述，对该项目的功能或用途进行简要说明。

使用方法：
```xml
<initialization_service_project>
<name>项目名称</name>
<description>项目描述（可选）</description>
</initialization_service_project>
```

示例：初始化一个名为 "Demo" 的后端服务项目

```xml
<initialization_service_project>
<name>Demo</name>
<description>这是一个示例项目</description>
</initialization_service_project>
```
