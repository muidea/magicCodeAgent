### read_stage_file

功能描述: 在指定阶段时，请求读取指定路径下的文件内容。当你需要检查一个你不知道内容的现有文件的内容时，例如分析代码、审查文本文件或从配置文件中提取信息时，请使用此功能。输出包括每行前的行号（例如 “1 | const x = 1”），这样在创建差异或讨论代码时就能更容易地引用特定行。通过指定 start_line 和 end_line 参数，你可以高效地读取大文件的特定部分，而无需将整个文件加载到内存中。自动从 PDF 和 DOCX 文件中提取原始文本。可能不适合其他类型的二进制文件，因为它会以字符串形式返回原始内容。

参数说明:
- stage: (必选) 当前指定的阶段名(例如： "code", "design", "test", "publish").
- path: (必选) 要读取的文件路径 (文件路径要求是相对项目根目录的相对路径)
- start_line: (可选) 要读取的起始行号（以 1 为基准）。如果未提供，则从文件开头读起。.
- end_line: (可选) 要读取的结束行号（以 1 为基准，包含在内）。如果没有提供，则读到文件的末尾.
调用方法:
<read_stage_file>
<stage>阶段名</stage>
<path>文件路径</path>
<start_line>起始行号 (可选)</start_line>
<end_line>结束行号 (可选)</end_line>
</read_stage_file>


示例:
1. 读取整个文件:
<read_stage_file>
<stage>code</stage>
<path>frontend-config.json</path>
</read_stage_file>

2. 从log/application.log中读取前1000行记录:
<read_stage_file>
<stage>test</stage>
<path>logs/application.log</path>
<end_line>1000</end_line>
</read_stage_file>

3. 从data/large-dataset.csv文件中读取500-1000行之间的记录:
<read_stage_file>
<stage>test</stage>
<path>data/large-dataset.csv</path>
<start_line>500</start_line>
<end_line>1000</end_line>
</read_stage_file>

4. 从src/app.ts文件中读取46-68行之间的记录，这个对应的是一个函数实现:
<read_stage_file>
<stage>code</stage>
<path>src/app.ts</path>
<start_line>46</start_line>
<end_line>68</end_line>
</read_stage_file>


注意：当同时提供 start_line 和 end_line 时，该工具只对请求的行进行高效流式处理，因此适用于处理日志、CSV 文件等大文件和其他大数据集，而不会产生内存问题。