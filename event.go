package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Event struct {
	ID      *uint64       `json:"id"`
	HostOrg *Organization `json:"host_org"`
	Title   string        `json:"title"`
	Start   time.Time     `json:"start"`
	End     time.Time     `json:"end"`
	Users   []User        `json:"users"`
}

func DbLoadEvent(db *sql.DB, id uint64, rec_depth int) *Event {
	if rec_depth < 0 {
		return nil
	}

	// load event
	var event Event
	var host_org_id *uint64
	err := db.QueryRow("SELECT * FROM events WHERE id =?", id).Scan(
		&event.ID,
		&host_org_id,
		&event.Title,
		&event.Start,
		&event.End,
	)
	if err != nil {
		return nil
	}

	if rec_depth < 1 {
		return &event
	}
	// load organization
	if host_org_id != nil {
		event.HostOrg = DbLoadOrganization(db, *host_org_id, rec_depth-1)
	}

	// load users
	query := `
	SELECT users.id FROM users
	JOIN events_users ON users.id = user_id
	JOIN events ON events.id = event_id
	WHERE events.id = ?
	`
	rows, err := db.Query(query, event.ID)
	if err != nil {
		panic(err)
	}
	var user_ids []uint64
	for rows.Next() {
		var user_id uint64
		err := rows.Scan(&user_id)
		if err != nil {
			panic(err)
		}
		user_ids = append(user_ids, user_id)
	}
	rows.Close()
	for _, user_id := range user_ids {
		event.Users = append(event.Users, *DbLoadUser(db, user_id, rec_depth-1))
	}

	return &event
}

func (e *Event) DbStore(db *sql.DB) *uint64 {
	// update event
	if e.ID != nil {
		_, err := db.Exec("UPDATE events SET title=?, start=?, end=? WHERE id = ?",
			e.Title,
			e.Start,
			e.End,
			e.ID,
		)
		if err != nil {
			panic(err)
		}

		// insert event
	} else {
		res, err := db.Exec("INSERT INTO events (title, start, end) VALUES (?,?,?)",
			e.Title,
			e.Start,
			e.End,
		)
		if err != nil {
			panic(err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		u_id := uint64(id)
		e.ID = &u_id
	}

	// store host organization
	if e.HostOrg != nil {
		e.HostOrg.ID = e.HostOrg.DbStore(db)
		if e.HostOrg.ID != nil {
			_, err := db.Exec("UPDATE events SET host_org_id = ? WHERE id = ?", e.HostOrg.ID, e.ID)
			if err != nil {
				panic(err)
			}
		}
	}

	// store users
	query := `
	SELECT * FROM users
	JOIN events_users ON users.id = user_id
	JOIN events ON events.id = event_id
	WHERE users.id = ? AND events.id = ?
	`
	for _, user := range e.Users {
		user.DbStore(db)
		rows, err := db.Query(query, e.ID, user.ID)
		if err != nil {
			panic(err)
		}
		if rows.Next() == false {
			_, err := db.Exec("INSERT INTO events_users (event_id, user_id) VALUES (?,?)", e.ID, user.ID)
			if err != nil {
				panic(err)
			}
		}
		rows.Close()
	}

	return e.ID
}
