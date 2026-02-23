package types

type ChatRequest struct {
	ConversationID string `json:"conversationId"`
	Message        string `json:"message"`
}

type DeleteConversationRequest struct {
	ConversationID string `json:"conversationId"`
}
