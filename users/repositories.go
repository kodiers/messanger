package users

import (
	"database/sql"
	"log"
	"messanger/libs/infrastructure/configuration"
)

type UserRepository struct {
	TableName string
	DB        *sql.DB
}

func InitUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		DB: db,
	}
}

var UserRep = InitUserRepository(configuration.DB)

func (uc *UserRepository) GetUserById(id int) (User, error) {
	row := uc.DB.QueryRow("SELECT * FROM USERS WHERE id=$1", id)
	user := new(User)
	err := row.Scan(&user.ID, &user.Username, &user.LastLogin, &user.PasswordHash, &user.Created, &user.Updated)
	if err != nil {
		log.Println("Could not parse rows from db.", err)
		return *user, err
	}
	return *user, nil
}

func (uc *UserRepository) GetUserByUsername(name string) (User, error) {
	row := uc.DB.QueryRow("SELECT * FROM USERS WHERE username=$1", name)
	user := new(User)
	err := row.Scan(&user.ID, &user.Username, &user.LastLogin, &user.PasswordHash, &user.Created, &user.Updated)
	if err != nil {
		log.Println("Could not parse rows from db.", err)
		return *user, err
	}
	return *user, nil
}

func (uc *UserRepository) InsertUser(user *User) (*User, error) {
	_, err := uc.DB.Exec("INSERT INTO USERS (USERNAME, PASSWORD_HASH) VALUES ($1, $2);",
		user.Username, user.PasswordHash)
	if err != nil {
		log.Println("Could not create user record ", err)
		return nil, err
	}
	createdUser, err := uc.GetUserByUsername(user.Username)
	if err != nil {
		log.Println("Could not get user record ", err)
		return nil, err
	}
	return &createdUser, nil
}
