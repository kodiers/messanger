package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	migrations2 "messanger/libs/infrastructure/database/migrations"
	"messanger/libs/utils"
)

func ConnectToDb(conStr string) *sql.DB {
	if conStr == "" {
		log.Fatalln("Could not get DB connections string.")
	}
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Fatalln("Could not connect to database.", err)
	}
	return db
}

func RunMigrations(migrationsFolder string, db *sql.DB) {
	repo := migrations2.InitMigrationRepository(migrationsFolder, db)
	migrations, ok := repo.GetAppliedMigrationsFromDb()
	files, _ := repo.GetMigrationsFiles()
	if !ok {
		log.Println("Cannot get info from migration table.")
	}
	var notAppliedMigrations []string
	var appliedMigrations []string
	for _, m := range migrations {
		appliedMigrations = append(appliedMigrations, m.Name)
	}
	for _, f := range files {
		if len(appliedMigrations) > 0 && utils.Contains(appliedMigrations, f) {
			continue
		}
		notAppliedMigrations = append(notAppliedMigrations, f)
	}

	for _, m := range notAppliedMigrations {
		log.Println("Applying migration: ", m)
		ok := repo.ApplyMigration(m)
		if ok {
			repo.CreateMigrationRecord(m)
		}
	}
}
