# mcp-go

简易实现的 MCP 服务框架，目前是在可以使用的地步。没有使用任何外部库或框架，直接下载代码就可使用

一个 Service 就是与一个模型，注册一个 Service 就会注册一个模型，对话时需要在 Service 上 注册一个交流的 Client ，目前 client 能保存上下文不受限制，需要特别注意内存

## 具体使用方法：

### 一、注册 service

通过提供的 NewMCPService 接口实现，具体操作如下：

service := mcp.NewMCPService("gemma-3-12b-it", "http://localhost:1234/v1/chat/completions", "")

其中第一个参数是大模型名称，例如 deepseek-chat 、 gemma-3-12b-it 等

第二个参数是远程服务地址，特别注意的是，地址需要包含端口号以及具体的路由，例如 chat/completions 、 v1/chat/completions

第三个参数则为 api-key key本地部署没有 api-key 的就填空key

### 二、开启新对话

通过调用 service 的 NewClient 接口开启一个新对话，代码如下：

dialog := s.NewClient(tag)

tag 是 Client 的标签，后续通过标签直接找到对应的 Client ，防止下标管理让人头疼

新对话可以保存上下文，但目前还没有任何的安全检查，后续加强

###  三、开始对话

通过 Client 的 Chat 接口与大模型进行对话，需要在对话前注册好工具，以便大模型调用。

### 四、清空对话历史

通过 Client 提供的 ClearHistory 接口清空对话，但仍然保留工具

### 五、工具注册

工具分为全局工具和局部工具，全局工具会附带在该 Service 的全部请求中，局部工具仅局限于当前 Client

全局工具通过 Service 提供的 AddGolbalTool 接口添加，局部工具通过 Client 提供的 AddTool 接口添加

具体操作如下 

```
dialog.AddTool("weather_query","查询指定位置的天气情况",
		mcp.Paramaters{
			Type: "object",
			Properties: map[string]any{
				"latitude": map[string]any{
					"type":        "number",
					"description": "latitude",
				},
				"longitude": map[string]any{
					"type":        "number",
					"description": "longitude",
				},
			},
			Required: []string{"latitude", "longitude"},
		},
		getWeather,
	)
```

AddGlobalTool 和 AddTool 操作是一样的，除了添加的位置不同之外其他都一样

第一个参数是工具名称，第二个参数是工具描述，第三个参数是描述工具的json信息，Type字段的值统一的 object ， Properties 的类型是 map[string]any ，用来填写工具参数， Properties 下面的 key 是参数名称，value 是描述参数的信息，比如是字符串就在 type 里填写string， description是对工具参数的描述，接下来的 required 字样存储的是 []string 的值，这里填写的是必须要的参数。第四个参数就是工具函数的 Handler ，参数类型统一为 map[string]any ，取参数的时候通过参数名称作为键读取参数值，注意错误检查，返回类型统一为 （string,err）也就是必须要把处理结果转换成字符串，并且抛出错误

### 五、文件注册

与工具注册的设计类似，分为全局文件添加和局部文件添加。参数非常简单，只有一个参数 path ， path 的值就是相对 main.go 文件的路径

## 新增内容：

1、实现文件的添加

2、优化提示词

3、加入工具数量限制，不能超过10个，超过部分无法再添加

## 已经实现：

1、实现与大模型的对话（目前仅支持非流式调用，后期加入参数实现流式与非流式可控转换）

2、实现大模型对工具调用

3、实现上下文的存储，没有上限存储上下文

4、添加了对带有 api key key型的支持（目前已经通过 deepseek api 以及本地 LMStudio 部署的大模型测试）

5、优化了注册工具的方式，可以一个函数注册整个工具

6、更高层次的抽象，一个 service 对应一种大模型，一个 Client 对应一段对话

7、优化项目结构，实现在 Service 端全局添加工具，无需为 Client 重复添加工具

8、为 Client 分配名称，减轻通过下表手动管理 Client 的烦恼

9、实现分步工具调用，增强通用性

10、优化提示词

11、加入工具数量限制，不能超过10个，超过部分无法再添加

## 计划实现：

1、错误检测过于简陋，后续重点加强错误检测，避免排错消耗大量精力

2、还需重点优化提示词，目前提示词仍处于婴幼儿阶段，稳定性表现不佳

3、加入更多限制，例如上下文长度，对不符合的提示词加以限制，以确保系统的高效运行