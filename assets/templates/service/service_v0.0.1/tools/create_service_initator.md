### create_service_initiator

描述: 创建一个新的服务模块 (initiator)
参数:
- **path**: (必填) 服务模块在项目中的相对路径，模块会创建在该目录下。
- **name**: (必填) 模块名。

用法示例:
```xml
<create_service_initiator>
  <path>internal/initiators</path>
  <name>routeRegistry</name>
</create_service_initiator>
```

说明:
在项目根目录 `/home/test/demo` 下的 `internal/initiators` 目录中创建一个名为 `routeRegistry` 的服务模块。
