package db

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"sec-ctl/cloud/config"

	"github.com/jmoiron/sqlx"
	// load postgres driver
	_ "github.com/lib/pq"
)

const bcryptSaltSize = 8

// UUID represents a PostgreSQL uuid
type UUID string

func (id UUID) String() string {
	return strings.Replace(string(id), "-", "", -1)
}

type DB struct {
	conn *sqlx.DB
}

type User struct {
	ID    UUID
	Email string
}

type Site struct {
	ID          UUID
	OwnerID     UUID `db:"owner_id"`
	StateShadow string
}

type Event struct {
	Level string
	Time  time.Time
	Data  string
}

func OpenDB(cfg config.Config) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort,
		cfg.DBUsername, cfg.DBPassword,
		cfg.DBName,
	)
	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to db: %v", err)
	}

	if _, err = conn.Exec("SELECT 1;"); err != nil {
		return nil, fmt.Errorf("unable to communicate with db: %v", err)
	}

	db := &DB{conn: conn}

	return db, nil
}

func (db *DB) AuthUser(email string, password string) (User, error) {

	var u User
	err := db.conn.Select(&u, "SELECT * FROM users where email = $1 AND password = crypt($2, password)", email, password)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (db *DB) AuthUserByToken(token string) (User, error) {
	var u User
	err := db.conn.Get(&u, `
		SELECT users.*
			FROM users
				JOIN auth_tokens ON auth_tokens.rec_id = users.id
			WHERE token = $1
	`)

	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (db *DB) CreateUser(email string, password string) (User, string, error) {

	u := User{
		Email: email,
	}

	tx, err := db.conn.Beginx()
	if err != nil {
		return User{}, "", err
	}

	r := tx.QueryRow(`
		INSERT
			INTO users(id, email, password)
			VALUES (gen_random_uuid(), $2, crypt($3, gen_salt('fo', $4))
			RETURNING id
	`, u.ID, email, password, bcryptSaltSize)
	err = r.Scan(&u.ID)
	if err != nil {
		return User{}, "", err
	}

	tok, err := createAuthToken(tx, u.ID, time.Time{})
	if err != nil {
		return User{}, "", err
	}

	if err = tx.Commit(); err != nil {
		return User{}, "", err
	}

	return u, tok, nil
}

func (db *DB) AuthSiteByToken(token string) (Site, error) {

	var s Site
	err := db.conn.Select(&s, "SELECT * FROM users where token = $1 AND owner_id IS NOT NULL", token)
	if err != nil {
		return Site{}, err
	}

	return s, nil
}

func (db *DB) ClaimSite(user User, siteID UUID, claimToken string) error {
	tx, err := db.conn.Beginx()
	if err != nil {
		return err
	}

	r, err := tx.Exec(`
		UPDATE sites
			SET owner_id = $1
			FROM auth_tokens
			WHERE users.id = $2
				AND owner_id IS NULL
				AND auth_tokens.rec_id = sites.id
				AND auth_tokens.token = $3
	`, user.ID, siteID, claimToken)
	if err != nil {
		return err
	}

	n, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("Invalid Site id, invalid claim token, or already claimed")
	}

	_, err = tx.Exec(`DELETE FROM auth_tokens WHERE token = $1`, claimToken)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (db *DB) FetchSiteByID(id UUID) (Site, error) {
	s := Site{}
	err := db.conn.Get(&s, `SELECT * FROM sites WHERE id = $1`, id)
	return s, err
}

func (db *DB) CreateSite() (Site, string, string, error) {

	tx, err := db.conn.Beginx()
	if err != nil {
		return Site{}, "", "", err
	}

	s := Site{}

	err = tx.QueryRow(`
		INSERT
			INTO sites(id)
			VALUES (gen_random_uuid())
			RETURNING normalize_uuid(id)
	`).Scan(&s.ID)

	if err != nil {
		return Site{}, "", "", err
	}

	tok, err := createAuthToken(tx, s.ID, time.Time{})
	if err != nil {
		return Site{}, "", "", err
	}

	tmpTok, err := createAuthToken(tx, s.ID, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return Site{}, "", "", err
	}

	if err = tx.Commit(); err != nil {
		return Site{}, "", "", err
	}

	return s, tok, tmpTok, nil
}

func (db *DB) SaveEvent(level string, tstamp time.Time, siteID UUID, evt interface{}) error {

	data, err := json.Marshal(evt)
	if err != nil {
		panic(fmt.Errorf("failed to jsonify event %v: %v", evt, err))
	}

	_, err = db.conn.Exec(`
		INSERT INTO events(level, time, site_id, data)
		VALUES ($1, $2, $3, $4)
	`, level, tstamp, siteID, data)
	return err
}

func (db *DB) GetLatestEvents(siteID UUID, max uint) ([]Event, error) {

	var evts []Event
	err := db.conn.Select(&evts, "SELECT * FROM events WHERE site_id = $1 ORDER BY time DESC LIMIT $2", siteID, max)
	if err != nil {
		return nil, err
	}

	return evts, nil

}

func createAuthToken(tx *sqlx.Tx, recID UUID, expiresAt time.Time) (string, error) {

	var expiresAtOrNull interface{}
	if expiresAt.IsZero() {
		expiresAtOrNull = nil
	} else {
		expiresAtOrNull = expiresAt
	}

	var tok string
	r := tx.QueryRow(`
		INSERT INTO
			auth_tokens(rec_id, token, expires_at)
			VALUES ($1, gen_random_uuid(), $2)
			RETURNING normalize_uuid(token)
	`, recID, expiresAtOrNull)

	err := r.Scan(&tok)
	if err != nil {
		return "", err
	}

	return tok, nil
}
