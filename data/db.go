package data

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const dataSource = "root:root@tcp(localhost)/shortlink"

var db *sql.DB
/*

CREATE TABLE `links` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `token` varchar(32) NOT NULL,
  `rawLink` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`)
)
 */

func init() {
	var err error
	db, err = sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
}

type NoTokenError struct {
	token string
}

func (c NoTokenError) Error() string {
	return "Try to read a nonexistent Token " + c.token
}

type DbError struct {
	err error
}

func (c DbError) Error() string {
	return "DB Error:" + c.err.Error()
}

func NoTokenErr(err error) bool {
	_, ok := err.(NoTokenError)
	return ok
}

func GetRawLink(token string) (string, error) {
	var rawLink string
	err := db.QueryRow("SELECT rawLink FROM links WHERE token = ?", token).Scan(&rawLink)
	switch {
	case err == sql.ErrNoRows:
		return "", NoTokenError{token}
	case err != nil:
		return "", DbError{err}
	default:
		return rawLink, nil
	}
}

func InsertLink(token string, rawLink string) (int64, error) {

	var id int64
	err := db.QueryRow("SELECT id FROM links WHERE token = ?", token).Scan(&id)

	switch err {
	case sql.ErrNoRows:
		ret, err := db.Exec("INSERT links (token, rawLink) VALUES (?, ?)", token, rawLink)
		if err != nil {
			return 0, err
		}
		id, err = ret.LastInsertId()
		if err != nil {
			return 0, err
		}
	case nil:
		return id, nil
	default:
		return 0, err
	}
	return id, nil
}
