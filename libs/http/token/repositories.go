package token

import (
	"database/sql"
	"log"
)

type TokenRepository struct {
	DB                    *sql.DB
	TokenExpirationInSecs int64
}

func InitTokenRepository(db *sql.DB, tokenExpiration int64) TokenRepository {
	return TokenRepository{
		DB:                    db,
		TokenExpirationInSecs: tokenExpiration,
	}
}

func (sr TokenRepository) CreateTokenRecord(token Token) bool {
	_, err := sr.DB.Exec("INSERT INTO tokens (TOKEN, USER_ID, EXPIRED) VALUES ($1, $2, $3);",
		token.Token, token.UserID, token.Expired)
	if err != nil {
		log.Println("Could not create session record ", err)
		return false
	}
	return true
}

func (sr TokenRepository) GetTokenFromDb(token string) (Token, error) {
	row := sr.DB.QueryRow("SELECT USER_ID, TOKEN, EXPIRED FROM tokens WHERE TOKEN=$1", token)
	tokObj := Token{}
	err := row.Scan(&tokObj.UserID, &tokObj.Token, &tokObj.Expired)
	if err != nil {
		return tokObj, err
	}
	return tokObj, nil
}
