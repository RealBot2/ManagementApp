
CREATE TABLE organizations (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
	org_id INTEGER,
	FOREIGN KEY (org_id) REFERENCES organizations(id)
);

CREATE TABLE events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	host_org_id INTEGER,
    title TEXT NOT NULL,
	start DATETIME NOT NULL,
	end DATETIME NOT NULL,
	FOREIGN KEY (host_org_id) REFERENCES organizations(id)
);

CREATE TABLE events_users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	event_id INTEGER,
	user_id INTEGER,
	FOREIGN KEY (event_id) REFERENCES events(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE meetings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	host_user_id INTEGER,
	event_id INTEGER,
    title TEXT NOT NULL,
    start DATETIME NOT NULL,
    end DATETIME NOT NULL,
	accepted INTEGER NOT NULL,
	denied INTEGER NOT NULL,
	FOREIGN KEY (host_user_id) REFERENCES users(id),
	FOREIGN KEY (event_id) REFERENCES events(id)
);

CREATE TABLE invitations (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	meeting_id INTEGER,
    invitee_id INTEGER,
	accepted INTEGER NOT NULL,
	denied INTEGER NOT NULL,
	FOREIGN KEY (meeting_id) REFERENCES meetings(id),
	FOREIGN KEY (invitee_id) REFERENCES users(id)
);
