## 代码规划目录结构

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
│   │   ├── index       // 首页目录
│   │   │   ├── index.tsx     // 首页组件
│   │   │   ├── index.less    // 首页样式文件（可选）
│   │   │   ├── model.ts      // 模型文件
│   │   │   ├── service.ts    // 服务文件
│   │   │   └── types         // 类型定义
│   │   │   │   ├── index.ts  // 辅助函数定义&类型导出
│   │   │   │   └── types.ts  // 类型定义
│   │   └── xxx             // 其他页面
│   │       ├── xxx         // 其他页面
│   │       │   ├── index.tsx   // 页面文件
│   │       │   ├── index.less  // 样式文件（可选）
│   │       │   ├── model.ts    // 模型文件
│   │       │   ├── service.ts  // 服务文件
│   │       │   └── types       // 类型定义
│   │       └── xxx         // 其他页面
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
├── .prettierignore     // prettier  配置文件
├── eslint.config.mjs   // eslint 配置文件
├── package.json        // 项目配置文件
├── tsconfig.json       // TypeScript 配置文件
└── typings.d.ts        // TypeScript 类型定义文件
```

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

### 服务定义说明

服务定义文件 `services.ts` 是一个 TypeScript 模块，用于定义项目的服务。服务是一个对象，对象中的每个属性表示一个服务，属性的名称表示服务名称，属性的值表示服务定义。

`service.ts` 文件内容格式如下：

```typescript
import { genApi } from "@/utils";
import { RequestConfig } from "@umijs/max";

const index: { [index: string]: string } = {
  // servcie接口定义清单
};

type APIFunctionType = {
  [key: string]: (
    val: unknown,
    headers?: Record<string, string>,
    config?: RequestConfig
  ) => Promise<unknown>;
};

const APIFunction: APIFunctionType = {};
Object.keys(index).forEach((key) => {
  APIFunction[key] = genApi(index[key]);
});

export default APIFunction;
```

以product模块为例，product模块的接口定义清单如下：

```typescript
import { genApi } from "@/utils";
import { RequestConfig } from "@umijs/max";

const index: { [index: string]: string } = {
  filterProductList: "GET /product/filter/",
  queryProduct: "GET /product/query/:id",
  insertProduct: "POST /product/insert/",
  updateProduct: "PUT /product/update/:id",
  deleteProduct: "DELETE /product/delete/:id",
  countProduct: "GET /product/count/",
};

type APIFunctionType = {
  [key: string]: (
    val: unknown,
    headers?: Record<string, string>,
    config?: RequestConfig
  ) => Promise<unknown>;
};

const APIFunction: APIFunctionType = {};
Object.keys(index).forEach((key) => {
  APIFunction[key] = genApi(index[key]);
});

export default APIFunction;

```

各个服务接口依据设计文档定义出对应的service模块

### 模型定义说明

模型定义文件 `model.ts` 是一个 TypeScript 模块，定义了模型接口和数据结构。

`model.ts` 文件内容格式如下, 以Product模型为例:

```typescript
import { ApiFilter, handleApiCall } from "@/utils/api";
import { useCallback, useState } from "react";
import api from "./service";
import { Product } from "./types";

const useProductModel = () => {
  const [loading, setLoading] = useState<boolean>(false);
  const [currentErrorMessage, setCurrentErrorMessage] = useState<string | null>(
    null
  );
  // 数据定义区
  const [currentProduct, setCurrentProduct] = useState<Product | null>(null);
  const [productList, setProductList] = useState<Product[]>([]);

  // 接口定义区
  const filterProduct = useCallback(
    async (filter: ApiFilter) => {
      handleApiCall(
        api.filterProductList(filter),
        (values) => {
          setProductList(values);
        },
        setCurrentErrorMessage,
        setLoading
      );
    },
    [loading, currentErrorMessage, productList]
  );

  const insertProduct = useCallback(
    async (val: Product): Promise<Product> => {
      return await handleApiCall(
        api.insertProduct(val),
        (value) => {
          setCurrentProduct(value);
          setProductList(productList ? [...productList, value] : [value]);
        },
        setCurrentErrorMessage,
        setLoading
      );
    },
    [loading, currentErrorMessage, currentProduct, productList]
  );
  // 其他接口定义

  // 数据和接口导出
  return {
    loading,
    currentErrorMessage,
    currentProduct,
    // 其他数据
    productList,
    filterProduct,
    insertProduct,
    // 其他接口
  };
};

export default useProductModel;
```

### 类型定义说明

定义当前模块使用的数据类型声明和各种辅助函数，并导出以便当前模块使用。

数据类型需要严格参照设计文档结构构进行定义，并使用类型检查工具进行类型检查。