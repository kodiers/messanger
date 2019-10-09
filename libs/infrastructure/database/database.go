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
	if migrations != nil && len(migrations) > 0 {
		for _, m := range migrations {
			if !utils.Contains(files, m.Name) {
				notAppliedMigrations = append(notAppliedMigrations, m.Name)
			}
		}
	} else {
		for _, f := range files {
			notAppliedMigrations = append(notAppliedMigrations, f)
		}
	}

	for _, m := range notAppliedMigrations {
		ok := repo.ApplyMigration(m)
		if ok {
			repo.CreateMigrationRecord(m)
		}
	}
}
