package request

type ChatRequest struct {
	ConversationID string `json:"conversationId"`
	Message        string `json:"message"`
}

type DeleteConversationRequest struct {
	ConversationID string `json:"conversationId"`
}
