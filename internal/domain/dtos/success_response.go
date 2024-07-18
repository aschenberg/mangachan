package dtos

type SuccessResponse struct {
	Message    string      `json:"msg"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}
