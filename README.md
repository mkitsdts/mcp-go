# mcp-go

简易实现的 MCP 服务框架，目前是在可以使用的地步。没有使用任何外部库或框架，直接下载代码就可使用

一个 Service 就是与一个模型，注册一个 Service 就会注册一个模型，对话时需要在 Service 上 注册一个交流的 Client ，目前 client 能保存上下文不受限制，需要特别注意内存

添加工具有两种方式，分别是 Service 端提供的 AddGlobalTool() 以及 Service 下属的 Client 提供的 AddTool() 函数。通过函数名很容易判断 AddGlobalTool 注册的是全局工具，也就是所有请求都会带上的工具。 AddTool 注册的工具只会附带在当前 Client 的请求中

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

## 新增更新：

1、实现分步工具调用，增强通用性

2、优化提示词

3、加入工具数量限制，不能超过10个，超过部分无法再添加

## 之前实现：

1、实现与大模型的对话（目前仅支持非流式调用，后期加入参数实现流式与非流式可控转换）

2、实现大模型对工具调用

3、实现上下文的存储，没有上限存储上下文

4、添加了对带有 api key key型的支持（目前已经通过 deepseek api 以及本地 LMStudio 部署的大模型测试）

5、优化了注册工具的方式，可以一个函数注册整个工具

6、更高层次的抽象，一个 service 对应一种大模型，一个 Client 对应一段对话

7、优化项目结构，实现在 Service 端全局添加工具，无需为 Client 重复添加工具

8、为 Client 分配名称，减轻通过下表手动管理 Client 的烦恼

## 计划实现：

1、错误检测过于简陋，后续重点加强错误检测，避免排错消耗大量精力

2、还需重点优化提示词，目前提示词仍处于婴幼儿阶段，稳定性表现不佳

3、加入更多限制，例如上下文长度，对不符合的提示词加以限制，以确保系统的高效运行