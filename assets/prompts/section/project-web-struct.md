## 项目规划目录结构

```plaintext
.
├── config              // 配置文件目录，集中管理项目配置，与 .umirc.ts 文件功能相同，二者二选一。
│   └── config.ts
├── dist                // 默认输出文件夹，执行 `max build` 后的构建产物。
├── mock                // 模拟数据文件夹，用于开发时的 Mock 数据。
│   └── app.ts｜tsx
├── src                 // 源代码目录
│   ├── .umi            // umi 的临时文件目录（开发环境）
│   ├── .umi-production // umi 的临时文件目录（生产环境）
│   ├── models          // 数据流管理模型文件夹
│   │   ├── global.ts   // 全局 Model 文件
│   │   └── index.ts    // Model 索引文件
│   ├── pages           // 页面目录（约定式路由）
│   │   ├── index.less  // 首页样式文件
│   │   └── index.tsx   // 首页组件
│   ├── utils           // 工具函数目录（推荐使用）
│   │   └── index.ts
│   ├── services        // 服务请求目录（推荐使用）
│   │   └── api.ts
│   ├── app.(ts|tsx)    // 运行时配置文件，可扩展运行时能力，如修改路由、render 方法等。
│   ├── global.ts       // 全局前置脚本文件
│   ├── global.(css|less|sass|scss) // 全局样式文件
│   ├── overrides.(css|less|sass|scss) // 高优先级全局样式文件
│   ├── favicon.(ico|gif|png|jpg|jpeg|svg|avif|webp) // 网站图标
│   └── loading.(tsx|jsx) // 全局加载组件
├── node_modules        // 项目依赖文件夹
├── .umirc.ts           // 非运行时配置文件，与 config/config 文件二选一
├── package.json        // 项目配置文件
├── tsconfig.json       // TypeScript 配置文件
└── typings.d.ts        // TypeScript 类型定义文件
````

---

### 路由定义说明

#### Pages 目录

* **约定式路由**：`pages/*` 文件夹的文件层级结构将自动生成路由表。
* **配置式路由**：若 `component` 写为相对路径，文件会从 `pages` 文件夹开始寻找。

示例配置：

```ts
routes: [
  // `./index` === `@/pages/index`
  { path: '/', component: './index' },
];
```

##### 基础路由

目录结构：

```plaintext
+ pages/
  + users/
    - index.tsx
  - index.tsx
```

生成的路由配置：

```ts
[
  { path: '/', component: '@/pages/index.tsx' },
  { path: '/users/', component: '@/pages/users/index.tsx' },
];
```

##### 动态路由

带 `$` 前缀的目录或文件为动态路由，例如：

```plaintext
+ pages/
  + foo/
    - $slug.tsx
  + $bar/
    - $.tsx
  - index.tsx
```

生成的路由配置：

```ts
[
  { path: '/', component: '@/pages/index.tsx' },
  { path: '/foo/:slug', component: '@/pages/foo/$slug.tsx' },
  { path: '/:bar/*', component: '@/pages/$bar/$.tsx' },
];
```

##### 404 页面

在约定式路由中，`pages/404.tsx` 会自动注册为全局 404 页面。若使用配置式路由，需要手动在最后一项添加兜底路由：

```ts
routes: [
  // other routes...
  { path: '/*', component: '@/pages/404.tsx' },
];
```

---

### Model 的使用说明

#### 创建 Model

数据流管理采用约定式目录结构，可在以下位置定义 Model 文件：

* `src/models`
* `src/pages/xxxx/models/`
* `src/pages/xxxx/model.{ts,tsx}`

**命名空间规则**：

| 路径                                      | 命名空间                | 说明                            |
| :-------------------------------------- | :------------------ | :---------------------------- |
| `src/models/count.ts`                   | `count`             | `src/models` 下不支持嵌套目录定义 Model |
| `src/pages/pageA/model.ts`              | `pageA.model`       |                               |
| `src/pages/pageB/models/product.ts`     | `pageB.product`     |                               |
| `src/pages/pageB/models/fruit/apple.ts` | `pageB.fruit.apple` | 支持嵌套目录定义 Model                |

#### 使用 Model

示例代码：

```tsx
// src/components/Username/index.tsx
import { useModel } from '@umijs/max';

export default function Page() {
  const { user, loading } = useModel('userModel');

  return loading ? <></> : <div>{user.username}</div>;
}
```

`useModel()` 的参数为 Model 的命名空间。

---

### 项目辅助工具和命令

1. **代码格式化**：

   ```bash
   yarn run format --quiet --write --ignore-path .prettierignore
   ```

2. **静态检查**：

   ```bash
   yarn lint --quiet --max-warnings 0
   ```

3. **安装依赖**：

   ```bash
   yarn install --silent --no-progress
   ```

4. **清理临时文件**：

   ```bash
   rm -rf ./dist src/.umi src/.umi-production
   ```

---

### 特别提示

1. **项目初始化**：请先使用 `initialization_web_project` 工具完成项目初始化，确保目录结构和配置文件正确生成。
2. **包管理器**：本项目使用 `yarn` 作为包管理工具。
