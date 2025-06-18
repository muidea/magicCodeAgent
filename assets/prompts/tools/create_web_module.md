### create_web_module

描述: 创建一个新的web模块 (module)  
参数:  
- **path**: (必填) 服务模块在项目中的相对路径，模块会创建在该目录下。  
- **name**: (必填) 模块名。  

用法示例:  
```xml
<create_web_module>
  <path>src/pages</path>
  <name>user</name>
</create_web_module>
```

说明:
在项目根目录 `/home/test/demo` 下的 `src/pages/` 目录中创建一个名为 `user` 的web模块。
