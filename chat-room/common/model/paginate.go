package commonmodel

type PaginatedRes[T any] struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Data  T   `json:"data"`
}
