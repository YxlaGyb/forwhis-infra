package request

type SetOptionRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
