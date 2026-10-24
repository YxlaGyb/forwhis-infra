package request

type CreateMetaRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Parent      int    `json:"parent"`
}
