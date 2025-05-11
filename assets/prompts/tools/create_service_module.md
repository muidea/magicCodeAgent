### create_service_module

Description: 创建一个新的服务模块(module)
Parameters:
- path: (required) 服务模块在项目中的相对路径，模块会创建在该目录下
- name: (required) 模块名
Usage:
<create_service_module>
<path>modulePath</path>
<name>moduleName</name>
</create_service_module>

Example: 在/home/test/demo/internal/modules/kernel目录下创建mcp模块, /home/test/demo是项目根目录

<create_service_module>
<path>internal/modules/kernel</path>
<name>mcp</name>
</create_service_module>