package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Meeting struct {
	ID          *uint64      `json:"id"`
	HostUser    *User        `json:"host_user"`
	Event       *Event       `json:"event"`
	Title       string       `json:"title"`
	Start       time.Time    `json:"start"`
	End         time.Time    `json:"end"`
	Invitations []Invitation `json:"invitations"`
	Accepted    bool         `json:"accepted"`
	Denied      bool         `json:"denied"`
}

func DbLoadMeeting(db *sql.DB, id uint64, rec_depth int) *Meeting {
	if rec_depth < 0 {
		return nil
	}

	// load meeting
	var m Meeting
	var user_id *uint64
	var event_id *uint64
	err := db.QueryRow("SELECT * FROM meetings WHERE id =?", id).Scan(
		&m.ID,
		&user_id,
		&event_id,
		&m.Title,
		&m.Start,
		&m.End,
		&m.Accepted,
		&m.Denied,
	)
	if err != nil {
		panic(err)
	}

	if rec_depth < 1 {
		return &m
	}

	// load host user
	if user_id != nil {
		m.HostUser = DbLoadUser(db, *user_id, rec_depth-1)
	}

	// load event
	if event_id != nil {
		m.Event = DbLoadEvent(db, *event_id, rec_depth-1)
	}

	// load invitations
	query := `
	SELECT invitations.id FROM invitations
	JOIN meetings ON meetings.id = meeting_id
	WHERE meetings.id = ?
	`
	rows, err := db.Query(query, m.ID)
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
		m.Invitations = append(m.Invitations, *DbLoadInvitation(db, id, rec_depth-1))
	}

	return &m
}

func (m *Meeting) DbStore(db *sql.DB) *uint64 {
	// update meeting
	query := `
	UPDATE meetings
	SET (title, start, end, accepted, denied)
	VALUES (?,?,?,?,?)
	WHERE id =?
	`
	if m.ID != nil {
		_, err := db.Exec(query,
			m.Title,
			m.Start,
			m.End,
			m.Accepted,
			m.Denied,
			m.ID,
		)
		if err != nil {
			panic(err)
		}

		// insert meeting
	} else {
		query = `
		INSERT INTO meetings
		(title, start, end, accepted, denied)
		VALUES (?,?,?,?,?)
		`
		res, err := db.Exec(query,
			m.Title,
			m.Start,
			m.End,
			m.Accepted,
			m.Denied,
		)
		if err != nil {
			panic(err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		u_id := uint64(id)
		m.ID = &u_id
	}

	// store host user
	if m.HostUser != nil {
		m.HostUser.ID = m.HostUser.DbStore(db)
		_, err := db.Exec("UPDATE meetings SET host_user_id =? WHERE id =?", m.HostUser.ID, m.ID)
		if err != nil {
			panic(err)
		}
	}

	// store event
	if m.Event != nil {
		m.Event.ID = m.Event.DbStore(db)
		_, err := db.Exec("UPDATE meetings SET event_id =? WHERE id =?", m.Event.ID, m.ID)
		if err != nil {
			panic(err)
		}
	}

	// store invitations
	for _, i := range m.Invitations {
		i.DbStoreExternID(db, m.ID, i.Invitee.ID)
	}

	return m.ID
}
