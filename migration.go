package main

import (
	"database/sql"
	"errors"
)

type Migration struct {
	Id       int
	Name     string
	Batch    int
	UpPath   string
	DownPath string
}

func NewMigration(Id int, Name string, Batch int, FileName string) *Migration {
	return &Migration{Id, Name, Batch, "", ""}
}

// Function to load migrations from a database.
func MigrationsFromDatabase(db *sql.DB) ([]Migration, error) {
	res, err := db.Query("SELECT * FROM migrations")
	if err != nil {
		return nil, errors.New("the migration table does not exist. Run `migrator init` first")
	}
	defer res.Close()
	var migrations []Migration
	for res.Next() {
		var m Migration
		if err := res.Scan(&m.Id, &m.Name, &m.Batch); err != nil {
			return nil, err
		}
		migrations = append(migrations, m)
	}
	return migrations, nil
}

// Function to get the current batch id. Either retrieves the batch id from the latest database entry or 1, as default.
func GetBatchId(migrations *[]Migration) int {
	var batchId int
	if len(*migrations) > 0 {
		batchId = (*migrations)[len(*migrations)-1].Batch + 1
	} else {
		batchId = 1
	}
	return batchId
}
