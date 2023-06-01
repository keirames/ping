package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	Conn *sqlx.DB
}

func NewUserRepository(conn *sqlx.DB) *userRepository {
	return &userRepository{
		Conn: conn,
	}
}

func (u *userRepository) IsExist(ids []int64) ([]int64, error) {
	type customRes struct {
		ID int64 `db:"id"`
	}

	var res []customRes

	query, args, err := sqlx.In("SELECT id FROM users WHERE users.id IN(?);", ids)
	if err != nil {
		fmt.Println("sql query error", err)
		return nil, err
	}
	fmt.Println(query)

	query = u.Conn.Rebind(query)
	err = u.Conn.Select(&res, query, args...)
	if err != nil {
		fmt.Println("sql select query error", err)
		return nil, err
	}

	var userIDs []int64
	for _, item := range res {
		userIDs = append(userIDs, item.ID)
	}

	return userIDs, nil
}
