### create_service_initator

Description: 创建一个新的服务模块(initator)
Parameters:
- path: (required) 服务模块在项目中的相对路径，模块会创建在该目录下
- name: (required) 模块名
Usage:
<create_service_initator>
<path>initatorPath</path>
<name>initatorName</name>
</create_service_initator>

Example: 在/home/test/demo/internal/initators目录下创建routeRegistry模块, /home/test/demo是项目根目录

<create_service_initator>
<path>internal/initators</path>
<name>routeRegistry</name>
</create_service_initator>