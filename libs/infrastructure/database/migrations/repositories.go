package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"messanger/libs/infrastructure/database/repositories"
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

func (mr MigrationRepository) GetAppliedMigrationsFromDb(db *sql.DB) ([]Migration, bool) {
	query := fmt.Sprintf("SELECT id, name, created FROM %v;", mr.TableName)
	rows, err := repositories.SelectQuery(db, query)
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
		err := rows.Scan(migration.Id, migration.Name, migration.Created)
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

func (mr MigrationRepository) ApplyMigration(db *sql.DB, filePath string) bool {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading file: ", filePath, " error: ", err)
		return false
	}
	queries := string(data)
	err = repositories.RunCreateInsertQuery(db, queries)
	if err != nil {
		log.Println("Error on running migration")
		return false
	}
	return true
}

func (mr MigrationRepository) CreateMigrationRecord(migrationName string) bool {
	query := fmt.Sprintf("INSERT INTO %v (name) VALUES ('%v');", mr.TableName, migrationName)
	err := repositories.RunCreateInsertQuery(mr.DB, query)
	if err != nil {
		log.Println("Error on create migration log. ", err)
		return false
	}
	return true
}
