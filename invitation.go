package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Invitation struct {
	ID       *uint64  `json:"id"`
	Meeting  *Meeting `json:"meeting"`
	Invitee  *User    `json:"invitee"`
	Accepted bool     `json:"Accepted"`
	Denied   bool     `json:"Denied"`
}

func DbLoadInvitation(db *sql.DB, id uint64, rec_depth int) *Invitation {
	if rec_depth < 0 {
		return nil
	}

	var i Invitation
	var meeting_id *uint64
	var invitee_id *uint64
	err := db.QueryRow("SELECT * FROM Invitation WHERE id=?", id).Scan(
		&i.ID,
		&meeting_id,
		&invitee_id,
		&i.Accepted,
		&i.Denied,
	)
	if err != nil {
		panic(err)
	}

	if meeting_id != nil {
		i.Meeting = DbLoadMeeting(db, *meeting_id, rec_depth-1)
	}
	if invitee_id != nil {
		i.Invitee = DbLoadUser(db, *invitee_id, rec_depth-1)
	}

	return &i
}

func (i *Invitation) DbStore(db *sql.DB) *uint64 {
	if i == nil {
		return nil
	}
	if i.Meeting == nil || i.Invitee == nil {
		return nil
	}

	return i.DbStoreExternID(db, i.Meeting.ID, i.Invitee.ID)
}

func (i *Invitation) DbStoreExternID(db *sql.DB, meeting_id *uint64, invitee_id *uint64) *uint64 {
	// update invitation
	_, err := db.Exec("UPDATE Invitation SET (meeting_id, invitee_id, accepted, denied) VALUES (?,?,?,?) WHERE id=?",
		meeting_id,
		invitee_id,
		i.Accepted,
		i.Denied,
		i.ID,
	)
	if err != nil {
		panic(err)
	}

	// insert invitation
	res, err := db.Exec("INSERT INTO invitations (meeting_id, invitee_id, accepted, denied) VALUES (?,?,?,?)",
		meeting_id,
		invitee_id,
		i.Accepted,
		i.Denied,
	)
	if err != nil {
		panic(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	u_id := uint64(id)
	i.ID = &u_id

	return i.ID
}
