package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID     *uint64       `json:"id"`
	Name   string        `json:"name"`
	Org    *Organization `json:"org"`
	Invite []Invitation  `json:"invitation"`
}

func DbLoadUser(db *sql.DB, id uint64, rec_depth int) *User {
	if rec_depth < 0 {
		return nil
	}

	var user User
	var org_id *uint64
	err := db.QueryRow("SELECT * FROM users WHERE id =?", id).Scan(&user.ID, &user.Name, &org_id)
	if err != nil {
		panic(err)
	}

	if rec_depth < 1 {
		return &user
	}

	// load organization
	if org_id != nil {
		user.Org = DbLoadOrganization(db, *org_id, rec_depth-1)
	}

	// load invitations
	query := `
	SELECT invitations.id FROM invitations
	JOIN users ON invitee_id = users.id
	WHERE users.id =?
	`
	rows, err := db.Query(query, user.ID)
	if err != nil {
		panic(err)
	}
	var invite_ids []uint64
	for rows.Next() {
		var invite_id uint64
		err := rows.Scan(&invite_id)
		if err != nil {
			panic(err)
		}
		invite_ids = append(invite_ids, invite_id)
	}
	rows.Close()
	for _, id := range invite_ids {
		user.Invite = append(user.Invite, *DbLoadInvitation(db, id, rec_depth-1))
	}

	return &user
}

func (u *User) DbStore(db *sql.DB) *uint64 {
	// update
	if u.ID != nil {
		_, err := db.Exec("UPDATE users SET name = ? WHERE id = ?", u.Name, u.ID)
		if err != nil {
			panic(err)
		}

		// insert
	} else {
		res, err := db.Exec("INSERT INTO users (name) VALUES (?)", u.Name)
		if err != nil {
			panic(err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		u_id := uint64(id)
		u.ID = &u_id
	}

	// store org
	if u.Org != nil {
		u.Org.ID = u.Org.DbStore(db)
		if u.Org.ID != nil {
			_, err := db.Exec("UPDATE users SET org_id =? WHERE id =?", u.Org.ID, u.ID)
			if err != nil {
				panic(err)
			}
		}
	}

	// store invitations
	for _, i := range u.Invite {
		i.DbStoreExternID(db, i.Meeting.ID, u.ID)
	}

	return u.ID
}
