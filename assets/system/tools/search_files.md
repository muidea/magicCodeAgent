### search_files

描述：请求对指定目录下的文件执行 regex 搜索，提供上下文丰富的结果。该工具可在多个文件中搜索模式或特定内容，显示每个匹配的封装上下文。

返回值格式如下：
```
search files {successed|failed}, path: {file_path} {状态信息}
---
[
  "file1.txt",
  "file2.ts",
  "file3.py"
]
```

参数：
- path: (必选) 要搜索的目录路径（相对于当前工作区目录）。该目录将被递归搜索。
- regex: (必选) 要搜索的正则表达式模式。使用 Rust regex 语法。
- file_pattern: (可选) 用于过滤文件的全局模式（例如，`*.ts`表示 TypeScript 文件）。如果不提供，将搜索所有文件 (*)。

使用方法:
```xml
<search_files>
<path>目录路径</path>
<regex>正则表达式</regex>
<file_pattern>过滤规则(可选)</file_pattern>
</search_files>
```

使用举例: 请求搜索src目录下的所有.ts文件
```xml
<search_files>
<path>src</path>
<regex>.*</regex>
<file_pattern>*.ts</file_pattern>
</search_files>
```

返回值
```
search files successed, path: src
---
[
    "src/app.ts",
    "src/index.ts"
]
```
