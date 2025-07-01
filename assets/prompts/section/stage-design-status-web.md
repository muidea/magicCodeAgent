## 项目状态

1. 当前处于功能设计阶段(design)，要求按照: 项目(project) -> 业务模块(module) -> 功能模块(block) 的层级关系，逐级进行功能分解。

2. 当前分解层级：{{.Level}}

3. 针对CURD类型的接口，要求设计的接口url符合`^\/(?:[^\/]+\/)*([^\/]+)\/(filter|query\/[^\/]+|insert|update\/[^\/]+|delete\/[^\/]+|count)\/?$`规则。
