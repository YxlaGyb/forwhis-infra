package types

type RejectReviewRequest struct {
	Reason string `json:"reason"`
}

type ApproveReviewRequest struct {
	NewValue string `json:"newValue"`
}
