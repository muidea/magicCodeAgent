### attempt_finish

**描述：**

每次使用工具后，用户会回复反馈信息，收到反馈信息后，根据返回信息决定后续需要继续执行的动作。

**注意：**  

在用户确认之前的工具使用成功之前，不得使用该工具。否则将导致代码损坏和系统故障。在使用此工具之前，您必须在 <thinking></thinking> 标记中询问自己是否已从用户处确认之前的任何工具使用都是成功的。如果没有，请勿使用此工具。


**使用方法：**  
<attempt_finish>
<result>
Your final result description here
</result>
</attempt_finish>

**参数：**  

- **result** (必填)： 任务结果。以最终结果的方式表述，不需要用户进一步输入。不要以提问或提供进一步帮助的方式结束结果。

**用法示例：**  

<attempt_finish>
<result>
已经完成产品信息更新
</result>
</attempt_finish>
