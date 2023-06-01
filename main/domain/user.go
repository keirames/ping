package domain

type User struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type UserRepository interface {
	IsExist(ids []int64) ([]int64, error)
}
