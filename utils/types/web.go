package types

// Vo 业务返回 VO
type Vo struct {
	Code     Code        `json:"code"`
	Page     int         `json:"page,omitempty"`
	PageSize int         `json:"page_size,omitempty"`
	Total    int         `json:"total,omitempty"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// ReplyMessage 对话回复消息结构
type ReplyMessage struct {
	Type    WsMsgType   `json:"type"` // 消息类别，start, end, img
	Content interface{} `json:"content"`
}

type WsMsgType string

const (
	WsContent = WsMsgType("content") // 输出内容
	WsEnd     = WsMsgType("end")
	WsErr     = WsMsgType("error")
)

// InputMessage 对话输入消息结构
type InputMessage struct {
	Content string `json:"content"`
	Tools   []int  `json:"tools"`  // 允许调用工具列表
	Stream  bool   `json:"stream"` // 是否采用流式输出
}

type Code int

const (
	Success       = Code(0)
	Failed        = Code(1)
	NotAuthorized = Code(401) // 未授权

	OkMsg       = "Success"
	ErrorMsg    = "系统开小差了"
	InvalidArgs = "非法参数或参数解析失败"
)
