### list_code_definition_names

描述: 请求列出源代码中的定义名称（类、函数、方法等）。该工具既可分析单个文件，也可分析指定目录顶层的所有文件。它能深入分析代码库结构和重要构造，封装对理解整体架构至关重要的高层概念和关系。

返回值格式如下：
```
list code definition names {successed|failed}, path: {file_path} {状态信息}
---
{定义信息}
```


参数:
- path: (必选) 要分析的文件或目录（相对于当前工作目录）的路径。如果给定一个目录，它会列出所有顶级源文件中的定义。

使用说明:
```xml
<list_code_definition_names>
<path>Directory path here</path>
</list_code_definition_names>
```

使用举例:

1. 列出特定文件(src/main.ts)中的定义：
```xml
<list_code_definition_names>
<path>src/main.ts</path>
</list_code_definition_names>
```

返回值：
```
list code definition names successed, path: src/main.ts
---
[
  {
    "file": "main.ts",
    "definitions": [
      {
        "type": "function",
        "name": "greet",
        "line": 2
      },
      {
        "type": "interface",
        "name": "IUser",
        "line": 6
      },
      {
        "type": "class",
        "name": "User",
        "line": 11,
        "implements": ["IUser"],
        "methods": [
          {
            "name": "constructor",
            "line": 16
          },
          {
            "name": "getDisplayName",
            "line": 21
          }
        ]
      },
      {
        "type": "namespace",
        "name": "Utils",
        "line": 30,
        "functions": [
          {
            "name": "isEmail",
            "line": 31
          }
        ]
      }
    ]
  }
]
```


2. 列出目录中所有文件的定义：
```xml
<list_code_definition_names>
<path>src/</path>
</list_code_definition_names>
```

返回值
```
list code definition names successed, path: src/
---json
[
  {
    "file": "src/index.ts",
    "defs": [
      ["function", "startApp"]
    ]
  },
  {
    "file": "src/user.ts",
    "defs": [
      ["interface", "IUser"],
      ["class", "User"],
      ["method", "constructor"],
      ["method", "greet"]
    ]
  },
  {
    "file": "src/utils/string.ts",
    "defs": [
      ["function", "capitalize"]
    ]
  }
]

```
