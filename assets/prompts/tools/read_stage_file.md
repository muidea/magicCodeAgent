### read_stage_file

描述：在指定阶段时，请求读取指定路径下的文件内容。当需要检查未知内容的现有文件时（例如，分析代码、审查文本文件或从配置文件中提取信息），使用此功能。输出包括每行前的行号（例如 "1 | const x = 1"），以便在创建差异或讨论代码时更容易引用特定行。通过指定 `start_line` 和 `end_line` 参数，可以高效地读取大文件的特定部分，而无需将整个文件加载到内存中。对于 PDF 和 DOCX 文件，会自动提取原始文本内容。不适合其他类型的二进制文件，因为该工具会以字符串形式返回其原始内容。

参数：
- **stage**: (必填) 当前指定的阶段名称（例如：`"code"`, `"design"`, `"test"`, `"publish"`）。
- **path**: (必填) 要读取的文件路径（路径应为相对项目根目录的相对路径）。
- **start_line**: (可选) 要读取的起始行号（以 1 为基准）。如果未提供，则从文件开头读取。
- **end_line**: (可选) 要读取的结束行号（以 1 为基准，包含在内）。如果未提供，则读取到文件末尾。

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

1. **读取整个文件**：

```xml
<read_stage_file>
<stage>code</stage>
<path>frontend-config.json</path>
</read_stage_file>
```

2. **从 `logs/application.log` 文件中读取前 1000 行记录**：

```xml
<read_stage_file>
<stage>test</stage>
<path>logs/application.log</path>
<end_line>1000</end_line>
</read_stage_file>
```

3. **从 `data/large-dataset.csv` 文件中读取 500-1000 行之间的记录**：

```xml
<read_stage_file>
<stage>test</stage>
<path>data/large-dataset.csv</path>
<start_line>500</start_line>
<end_line>1000</end_line>
</read_stage_file>
```

4. **从 `src/app.ts` 文件中读取 46-68 行之间的记录（对应一个函数实现）**：

```xml
<read_stage_file>
<stage>code</stage>
<path>src/app.ts</path>
<start_line>46</start_line>
<end_line>68</end_line>
</read_stage_file>
```

注意：当同时提供 `start_line` 和 `end_line` 时，该工具仅对请求的行进行高效流式处理，非常适合处理日志、CSV 文件和其他大数据集，无需担心内存问题。
