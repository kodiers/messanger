package users

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"messanger/libs/infrastructure/configuration"
	"time"
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

func (uc *UserRepository) loadLocation() time.Location {
	location, _ := time.LoadLocation("Europe/Moscow")
	return *location
}

func (uc *UserRepository) parseUserFromRow(row sql.Row) (User, error) {
	var lastLogin, created, updated string
	user := new(User)
	err := row.Scan(&user.ID, &user.Username, &lastLogin, &user.PasswordHash, &created, &updated)
	if err != nil {
		log.Println("Could not parse rows from db.", err)
		return *user, err
	}
	location := uc.loadLocation()
	user.LastLogin, _ = time.ParseInLocation(time.RFC3339, lastLogin, &location)
	user.Created, _ = time.ParseInLocation(time.RFC3339, created, &location)
	user.Updated, _ = time.ParseInLocation(time.RFC3339, updated, &location)
	return *user, err
}

func (uc *UserRepository) GetUserById(id int) (User, error) {
	row := uc.DB.QueryRow("SELECT * FROM USERS WHERE id=$1", id)
	user, _ := uc.parseUserFromRow(*row)
	return user, nil
}

func (uc *UserRepository) GetUserByUsername(name string) (User, error) {
	row := uc.DB.QueryRow("SELECT * FROM USERS WHERE username=$1", name)
	user, _ := uc.parseUserFromRow(*row)
	return user, nil
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

func (uc *UserRepository) GetUsers() ([]User, error) {
	rows, err := uc.DB.Query("SELECT * FROM USERS;")
	if err != nil {
		log.Println("Could not get users from db", err)
		return nil, err
	}
	users := make([]User, 0)
	for rows.Next() {
		user := new(User)
		var lastLogin, created, updated string
		err := rows.Scan(&user.ID, &user.Username, &lastLogin, &user.PasswordHash, &created, &updated)
		if err != nil {
			log.Println("Could not read rows data", err)
		}
		location := uc.loadLocation()
		user.LastLogin, _ = time.ParseInLocation(time.RFC3339, lastLogin, &location)
		user.Created, _ = time.ParseInLocation(time.RFC3339, created, &location)
		user.Updated, _ = time.ParseInLocation(time.RFC3339, updated, &location)
		users = append(users, *user)
	}
	_ = rows.Close()
	return users, nil
}

func (uc *UserRepository) UpdateUser(user User) error {
	_, err := uc.DB.Exec("UPDATE users SET LAST_LOGIN=$1, UPDATED=$2", pq.FormatTimestamp(user.LastLogin), pq.FormatTimestamp(time.Now()))
	if err != nil {
		log.Println("Could not update user", err)
		return err
	}
	return nil
}

func (uc *UserRepository) UpdateUserLastLogin(user User) {
	UserUpdateMutex.Lock()
	user.LastLogin = time.Now()
	err := UserRep.UpdateUser(user)
	defer UserUpdateMutex.Unlock()
	if err != nil {
		log.Println("Cannot update user LastLogin field")
	}
}
