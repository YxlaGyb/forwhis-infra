package types

type ToggleLikeRequest struct {
	TargetType string `json:"target_type" binding:"required,oneof=content comment"`
	TargetID   int    `json:"target_id" binding:"required,min=1"`
}
