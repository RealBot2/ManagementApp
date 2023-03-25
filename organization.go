package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Organization struct {
	ID   *uint64 `json:"id"`
	Name string  `json:"name"`
}

func DbLoadOrganization(db *sql.DB, id uint64, rec_depth int) *Organization {
	if rec_depth < 0 {
		return nil
	}

	var organization Organization
	err := db.QueryRow("SELECT * FROM organizations WHERE id =?", id).Scan(&organization.ID, &organization.Name)
	if err != nil {
		panic(err)
	}
	return &organization
}

func (o *Organization) DbStore(db *sql.DB) *uint64 {
	// update
	if o.ID != nil {
		_, err := db.Exec("UPDATE organizations SET name = ? WHERE id = ?", o.Name, o.ID)
		if err != nil {
			panic(err)
		}
		return o.ID
	}

	// insert
	res, err := db.Exec("INSERT INTO organizations (name) VALUES (?)", o.Name)
	if err != nil {
		panic(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	u_id := uint64(id)
	return &u_id
}
