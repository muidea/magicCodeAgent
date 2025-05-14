## 项目规划目录结构

```
├── .gitignore           # Git忽略规则文件
├── Makefile             # 项目构建Makefile
├── README.md            # 项目说明文档
├── application/         # 应用程序主入口
│   └── {服务名}/         # 单个服务目录
│       └── cmd /        # 服务启动入口
│           └── main.go  # main函数
├── bin/                 # 编译生成的可执行文件目录
├── build.sh             # 项目构建脚本
├── doc/                 # 项目文档
├── docker/              # Docker相关配置
│   ├── bootstrap.sh     # Docker启动脚本
│   └── build.sh         # Docker构建脚本
├── internal/            # 内部模块(主要业务逻辑)
│   ├── config/          # 配置管理
│   ├── initators/       # 各类initator模块,每个模块对应不同子目录
│   │   ├── pprof/       # 性能分析initator模块，默认内置基础模块
│   │   ├── routeregistry/ # 路由管理initator模块，支持HTTP API访问的基础模块(可选)
│   │   └── timer/       # 定时任务initator模块，支持启动定时任务的基础模块(可选)
│   │   └── {自定义模块名}/ # 其他自定义initator模块(可选)
│   └── modules/         # 业务模块，包含各类子功能的业务模块
│       ├── blocks/      # 基础子功能业务块模块，每个子功能模块对应不同子目录
│       │   └── {自定义模块名}/ # 各业务自定义基础业务模块
│       └── kernel/      # 核心业务模块模块
│       │   └── {自定义模块名}/ # 各业务自定义业务模块
│       └── {自定义模块名}/     # 各业务自定义业务模块
├── pkg/                 # 公共库(可被外部引用)
├── test/                # 测试相关代码目录
└── vendor/              # 项目依赖(Go vendor目录)
```

### initator与module的目录结构

├── biz/                # 业务逻辑层，包含业务逻辑代码，包括发起和响应模块间EventHub调用，initator不支持通过EventHub进行模块间调用
├── module.go           # 模块主文件，实现模块接口
├── pkg/                # 公共定义和工具
└── service/            # 服务层，提供HTTP API接口，initator不包含service层

### initator与module的区别说明

1. initator: 系统内部模块，主要用来辅助系统内部的功能，以插件的方式扩展系统基础功能，对外提供业务访问接口，不支持通过EventHub的方式进行功能调用。initator可以通过create_service_initator工具创建，再完善模块功能

2. module: 业务模块，主要用来实现业务功能，对外提供业务访问接口，支持通过EventHub的方式进行功能调用，module可以包含多个子功能模块，如blocks、kernel等，各个module之间都可以对外提供子功能项，可以调用initators中的功能模块，各个module之间通过EventHub的方式进行功能调用。module可以通过create_service_module工具创建，再完善模块功能。各个自定义module的存放目录，根据其业务承载的业务性质进行分组，如blocks、kernel等。

3. 所有的initator和module在定义完成后，都需要在main.go 进行显示引用，以便系统启动时加载

**强调说明：如果当前项目未完成初始化，请先执行`initialization_service_project`工具进行初始化**