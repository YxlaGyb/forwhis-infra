package types

type CreateCommentRequest struct {
	Cid    int    `json:"cid"`
	Text   string `json:"text"`
	Url    string `json:"url"`
	Type   string `json:"type"`
	Parent int    `json:"parent"`
}
