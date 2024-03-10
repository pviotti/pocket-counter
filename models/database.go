package models

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const DATABASE_PATH = "./data/pocket.db"

func GetUnreadCountsForPastYearFromDB() ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	var unreadCounts []map[string]interface{}
	rows, err := db.Query("select date, count from unread_count where date >= date('now', '-1 year') order by date")
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		unreadCounts = append(unreadCounts, map[string]interface{}{
			"date":         date,
			"unread_count": count,
		})
	}

	return unreadCounts, nil
}

func SaveToDatabase(unreadCount int) error {
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`create table if not exists unread_count (
        date date primary key,
        count integer
    )`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	date := time.Now().Format("2006-01-02")
	_, err = db.Exec("insert or replace into unread_count (date, count) values (?, ?)", date, unreadCount)
	if err != nil {
		return fmt.Errorf("error inserting into table: %v", err)
	}

	return nil
}
