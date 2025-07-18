### read_stage_file

描述：读取指定阶段，指定路径下的文件内容。在需要检查非当前阶段的指定文件的内容时（例如，在代码开发阶段读取设计阶段的设计文件，测试验证阶段读取代码开发阶段的指定代码），使用此功能。输出包括每行前的行号（例如 "1 | const x = 1"），以便在创建差异或讨论代码时更容易引用特定行。通过指定 `start_line` 和 `end_line` 参数，可以高效地读取大文件的特定部分，而无需将整个文件加载到内存中。对于 PDF 和 DOCX 文件，会自动提取原始文本内容。不适合其他类型的二进制文件，因为该工具会以字符串形式返回其原始内容。

返回值格式如下：
```
read stage file {successed|failed|warning}, path: {file_path} {描述信息}
---
{文件内容}
```

参数：
- stage: (必填) 指定的阶段名称（例如：`code`, `design`, `test`, `publish`）。
- path: (必选) 要读取的文件路径（相对于当前工作区目录）
- start_line: (可选) 要读取的起始行号（以 1 为基准）。如果未提供，则从文件开头读起。
- end_line: (optional) 要读取的结束行号（以 1 为单位，包含在内）。如果没有提供，则读到文件的末尾。

使用方法：
```xml
<read_stage_file>
<stage>阶段名称</stage>
<path>文件路径</path>
<start_line>起始行号（可选）</start_line>
<end_line>结束行号（可选）</end_line>
</read_stage_file>
```

示例：

1. 读取code阶段整个frontend-config.json文件
```xml
<read_stage_file>
<stage>code</stage>
<path>frontend-config.json</path>
</read_stage_file>
```

返回值:
```
read stage file successed, path: frontend-config.json
---
1| {
2|  "compilerOptions": {
3|    "target": "esnext",
4|    "module": "esnext",
5|    "moduleResolution": "node",
6|    "allowSyntheticDefaultImports": true,
7|    "skipLibCheck": true
8|  },
9|  "include": [
10|    "src/**/*",
11|    ".umirc.ts",
12|    "config/**/*",
13|    ".eslintrc.js",
14|    "typings.d.ts"
15|  ],
16|  "exclude": ["node_modules", "dist", ".umi", ".umi-production"]
17|}
```


2. 读取test阶段大型日志文件logs/application.log的前 1000 行：
```xml
<read_stage_file>
<stage>test</stage>
<path>logs/application.log</path>
<end_line>1000</end_line>
</read_stage_file>
```

返回值：
```
read file failed, path: logs/application.log, error: File not found
```

3. 读取code阶段源文件src/app.ts中的特定函数
```xml
<read_stage_file>
<stage>code</stage>
<path>src/app.ts</path>
<start_line>46</start_line>
<end_line>68</end_line>
</read_stage_file>
```

返回值：
```
read file warning, path: src/app.ts, 当前文件内容为空
```


注意：当同时提供 `start_line` 和 `end_line` 时，该工具仅对请求的行进行高效流式处理，非常适合处理日志、CSV 文件和其他大数据集，无需担心内存问题。
