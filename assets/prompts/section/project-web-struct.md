## 项目规划目录结构

```
.
├── config // 与 [`.umirc.ts`]相同，区别是你可以单独在一个 `config` 文件夹下集中管理所有的配置，保持项目根目录整洁。
│   └── config.ts
├── dist  // 执行 `max build` 后产物的默认输出文件夹。
├── mock
│   └── app.ts｜tsx
├── src
│   ├── .umi
│   ├── .umi-production
│   ├── models
│   │   ├── global.ts
│   │   └── index.ts
│   ├── pages
│   │   ├── index.less
│   │   └── index.tsx
│   ├── utils // 推荐目录
│   │   └── index.ts
│   ├── services // 推荐目录
│   │   └── api.ts
│   ├── app.(ts|tsx)  // 运行时配置文件，可以在这里扩展运行时的能力，比如修改路由、修改 render 方法等。
│   ├── global.ts   //全局前置脚本文件。
│   ├── global.(css|less|sass|scss)   //全局样式文件。
│   ├── overrides.(css|less|sass|scss)  //高优先级全局样式文件。
│   ├── favicon.(ico|gif|png|jpg|jpeg|svg|avif|webp)
│   └── loading.(tsx|jsx)   // 全局加载组件。
├── node_modules
├── .umirc.ts // 非运行时配置 与 config/config 文件 2 选一
├── package.json
├── tsconfig.json
└── typings.d.ts
```

### 路由定义说明

#### pages 目录

约定式路由默认以 `pages/*` 文件夹的文件层级结构来生成路由表。

在配置式路由中，`component` 若写为相对路径，将从该文件夹为起点开始寻找文件：

```ts
  routes: [
    // `./index` === `@/pages/index`
    { path: '/', component: './index' }
  ]
```

##### 基础路由

假设 `pages` 目录结构如下：

```
+ pages/
  + users/
    - index.tsx
  - index.tsx
```

那么，会自动生成路由配置如下：

```ts
[
  { path: '/', component: '@/pages/index.tsx' },
  { path: '/users/', component: '@/pages/users/index.tsx' },
]
```

##### 动态路由

约定带 `$` 前缀的目录或文件为动态路由。若 `$` 后不指定参数名，则代表 `*` 通配，比如以下目录结构：

```
+ pages/
  + foo/
    - $slug.tsx
  + $bar/
    - $.tsx
  - index.tsx
```

会生成路由配置如下：

```ts
[
  { path: '/', component: '@/pages/index.tsx' },
  { path: '/foo/:slug', component: '@/pages/foo/$slug.tsx' },
  { path: '/:bar/*', component: '@/pages/$bar/$.tsx' },
]
```

##### pages/404.tsx

在使用约定式路由时，该文件会自动被注册为全局 404 的 fallback 页面。若你使用配置式路由，需要自行配置兜底路由到路由表最后一个：

```ts
  routes: [
    // other routes ...
    { path: '/*', component: '@/pages/404.tsx' }
  ]
```

### Model的使用说明

#### 创建 Model

数据流管理插件采用约定式目录结构，我们约定可以在 `src/models`, `src/pages/xxxx/models/`目录中，和 `src/pages/xxxx/model.{ts,tsx}` 文件引入 Model 文件。
Model 文件允许使用 `.(tsx|ts)` 四种后缀格式，**命名空间（namespace）** 生成规则如下。

| 路径 | 命名空间 | 说明 |
| :--- |:--- | :--- |
| `src/models/count.ts` | `count` | `src/models` 目录下不支持目录嵌套定义 model |
| `src/pages/pageA/model.ts` | `pageA.model` |  |
| `src/pages/pageB/models/product.ts` | `pageB.product` |  |
| `src/pages/pageB/models/fruit/apple.ts` | `pageB.fruit.apple` |  `pages/xxx/models` 下 model 定义支持嵌套定义 |

当我们需要获取 Model 中的全局数据时，调用该命名空间即可。例如，对于 Model 文件 `userModel.ts`，它的命名空间为 `userModel`。

#### 使用 Model

现在，您想要在某个组件中使用全局的 Model。以用户信息为例，只需要调用 `useModel` 这一钩子函数：

```tsx
// src/components/Username/index.tsx
import { useModel } from '@umijs/max';

export default function Page() {
  const { user, loading } = useModel('userModel');

  return (
    {loading ? <></>: <div>{user.username}</div>}
  );
}
```

其中，`useModel()` 方法传入的参数为 Model 的**命名空间**


**强调说明：如果当前项目未完成初始化，请先执行`initialization_web_project`工具进行初始化**