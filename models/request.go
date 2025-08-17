package models

type ProductRequest struct {
	DocuTitle   string                 `json:"docu_title"` // 文档标题
	Folder      string                 `json:"folder"`     // 文档所在文件夹
	ChatID      string                 `json:"chat_id"`
	Description string                 `json:"description"`
	Stream      bool                   `json:"stream,omitempty"`    // 是否开启流式响应
	Detail      bool                   `json:"detail,omitempty"`    // 是否返回详细信息
	Variables   map[string]interface{} `json:"variables,omitempty"` // 可选的变量
}

type GptResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
