### list_stage_files

描述：在指定阶段时，请求列出指定目录内的文件和目录。如果递归为 `true`，将递归列出所有文件和目录。如果递归为 `false` 或未提供，则只会列出顶层内容。请勿使用此工具来确认您可能已创建的文件是否存在，因为用户会通知您文件是否已成功创建。

返回值格式如下：
```
list stage files {successed|failed}, path: {file_path} {状态信息}
---
[
  "file1.txt",
  "file2.ts",
  "file3.py"
]
```

参数
- stage: (必填) 指定的阶段名称（例如：`code`, `design`, `test`, `publish`）。
- path: (必选) 要列出内容的目录路径（相对于当前工作区目录）
- recursive: (可选) 是否递归列出文件。使用 true 表示递归列出，使用 false 或省略表示只列出顶级文件。


使用方法：
```xml
<list_stage_files>
<stage>阶段名称</stage>
<path>目录路径</path>
<recursive>true 或 false（可选）</recursive>
</list_stage_files>
```

使用示例：请求列出 `code` 阶段当前目录下的所有文件，不要求递归

```xml
<list_stage_files>
<stage>code</stage>
<path>src</path>
<recursive>false</recursive>
</list_stage_files>
```

返回值
```
list stage files successed, path: src
---
[
    "src/app.ts",
    "src/index.ts"
]
```
