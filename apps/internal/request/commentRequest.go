package request

type CreateCommentRequest struct {
	Cid    int    `json:"cid"`
	Text   string `json:"text"`
	Url    string `json:"url"`
	Parent int    `json:"parent"`
}
