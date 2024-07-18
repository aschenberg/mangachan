package dtos

type Pagination struct {
	Page  int64 `bson:"page"`
	Limit int64 `bson:"limit"`
	Total int64 `bson:"total"`
}
