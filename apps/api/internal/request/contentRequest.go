package request

type CreateContentRequest struct {
	Title       string `json:"title"`
	Text        string `json:"text"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Password    string `json:"password"`
	CategoryMid int    `json:"categoryMid"`
	TagIds      []int  `json:"tagIds"`
}

type ContentPasswordRequest struct {
	Password string `json:"password"`
}
