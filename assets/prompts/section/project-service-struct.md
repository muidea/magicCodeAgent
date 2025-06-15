## 项目规划目录结构

```plaintext
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
│   ├── initiators/      # 各类initiator模块，每个模块对应不同子目录
│   │   ├── pprof/       # 性能分析initiator模块，默认内置基础模块
│   │   ├── routeregistry/ # 路由管理initiator模块，支持HTTP API访问的基础模块(可选)
│   │   ├── timer/       # 定时任务initiator模块，支持启动定时任务的基础模块(可选)
│   │   └── {自定义模块名}/ # 其他自定义initiator模块(可选)
│   └── modules/         # 业务模块，包含各类子功能的业务模块
│       ├── blocks/      # 基础子功能业务块模块，每个子功能模块对应不同子目录
│       │   └── {自定义模块名}/ # 各业务自定义基础业务模块
│       ├── kernel/      # 核心业务模块模块
│       │   └── {自定义模块名}/ # 各业务自定义核心业务模块
│       └── {自定义模块名}/     # 其他自定义业务模块
├── pkg/                 # 公共库(可被外部引用)
├── test/                # 测试相关代码目录
└── vendor/              # 项目依赖(Go vendor目录)
```

---

### Initiator与Module的目录结构

```plaintext
├── biz/                # 业务逻辑层，包含业务逻辑代码，包括发起和响应模块间EventHub调用，initiator不支持通过EventHub进行模块间调用
├── module.go           # 模块主文件，实现模块接口
├── pkg/                # 公共定义和工具
└── service/            # 服务层，提供HTTP API接口，initiator不包含service层
```

---

### Initiator与Module的区别说明

1. **Initiator**：系统内部模块，主要用于扩展系统内部功能，以插件的方式实现基础功能支持，不支持通过EventHub进行模块间调用。适用于提供通用功能（如性能监控、路由管理等）。可以通过 `create_service_initiator` 工具创建并完善模块功能。

2. **Module**：业务模块，专注于实现具体的业务功能。支持通过EventHub进行模块间调用，可包含多个子功能模块（如 `blocks`、`kernel`）。Module之间可以互相调用initiators中的功能模块。可以通过 `create_service_module` 工具创建并完善功能，并根据业务性质分组存放（如 `blocks`、`kernel`）。

3. **引用说明**：所有定义完成的initiator和module，必须在 `main.go` 中显式引用，以确保系统启动时正确加载。

---

### 项目辅助工具和命令

1. **代码格式化**：使用 `execute_command` 工具执行以下命令格式化代码：

   ```bash
   go fmt ./...
   ```

2. **静态检查**：使用 `execute_command` 工具执行以下命令对项目进行静态检查：

   ```bash
   go vet ./...
   ```

3. **依赖整理**：使用 `execute_command` 工具执行以下命令整理和优化 `go.mod` 文件：

   ```bash
   go mod tidy
   ```

4. **运行测试**：使用 `execute_command` 工具执行以下命令运行所有测试用例：

   ```bash
   go test ./...
   ```
