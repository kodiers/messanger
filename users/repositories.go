package users

import (
	"database/sql"
	"log"
)

type UserRepository struct {
	TableName string
	DB        *sql.DB
}

func MakeUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		TableName: "USERS",
		DB:        db,
	}
}

func (uc *UserRepository) getUserById(id int) (User, error) {
	row := uc.DB.QueryRow("SELECT * FROM USERS WHERE id=$1", id)
	user := new(User)
	err := row.Scan(&user.ID, &user.Username, &user.LastLogin, &user.PasswordHash, &user.Created, &user.Updated)
	if err != nil {
		log.Println("Could not parse rows from db.", err)
		return *user, err
	}
	return *user, nil
}
