package mcp

const (
	system_extract_paramater_prompt string = "我会给你几个工具，工具里会包含使用工具需要的参数。如果我提出问题，你要从这些工具中上找到最可能解决我的问题的工具，然后把我说的话转换成工具名称和并且提取出参数，我把使用工具的结果告诉你。如果我告诉你结果，你要把结果给我提取出来。"
	system_extarct_answer_prompt    string = "我告诉你结果，你要把结果提取出来。"
	request_content_type            string = "application/json; charset=utf-8"
)
