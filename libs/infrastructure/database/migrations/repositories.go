package migrations

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type MigrationRepository struct {
	FolderPath string
	TableName  string
	DB         *sql.DB
}

func InitMigrationRepository(folderPath string, db *sql.DB) MigrationRepository {
	return MigrationRepository{FolderPath: folderPath, TableName: "migrations", DB: db}
}

func (mr MigrationRepository) GetAppliedMigrationsFromDb() ([]Migration, bool) {
	quoted := pq.QuoteIdentifier(mr.TableName)
	rows, err := mr.DB.Query(fmt.Sprintf("SELECT id, name, created FROM %s;", quoted))
	if err != nil {
		if rows != nil {
			_ = rows.Close()
		}
		log.Println("Table 'migrations' does not exists. Will create it.")
		return nil, false
	}
	migrations := make([]Migration, 0)
	for rows.Next() {
		migration := new(Migration)
		err := rows.Scan(&migration.Id, &migration.Name, &migration.Created)
		if err != nil {
			log.Println("Could not read rows data", err)
		}
		migrations = append(migrations, *migration)
	}
	_ = rows.Close()
	return migrations, true
}

func (mr MigrationRepository) GetMigrationsFiles() ([]string, bool) {
	var files []string
	err := filepath.Walk(mr.FolderPath, func(path string, info os.FileInfo, err error) error {
		fInfo, _ := os.Stat(path)
		if !fInfo.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Println("Could not get files from migration folder.")
		return nil, false
	}
	return files, true
}

func (mr MigrationRepository) ApplyMigration(filePath string) bool {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading file: ", filePath, " error: ", err)
		return false
	}
	queries := string(data)
	_, err = mr.DB.Exec(queries)
	if err != nil {
		log.Println("Error on running migration", err)
		return false
	}
	return true
}

func (mr MigrationRepository) CreateMigrationRecord(migrationName string) bool {
	quoted := pq.QuoteIdentifier(mr.TableName)
	_, err := mr.DB.Exec(fmt.Sprintf("INSERT INTO %s (name) VALUES ($1);", quoted), migrationName)
	if err != nil {
		log.Println("Error on create migration log. ", err)
		return false
	}
	return true
}
