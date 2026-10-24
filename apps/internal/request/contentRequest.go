package request

type CreateContentRequest struct {
	Title       string `json:"title"`
	Text        string `json:"text"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Password    string `json:"password"`
	Template    string `json:"template"`
	CategoryMid int    `json:"categoryMid"`
}

type ContentPasswordRequest struct {
	Password string `json:"password"`
}
