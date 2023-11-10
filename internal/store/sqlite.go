package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func SetDBFilepath(filepath string) {
	var err error
	db, err = sql.Open("sqlite3", filepath)
	if err != nil {
		fmt.Println("Failed to open database:", filepath, err)
		return
	}

	if _, err = db.Exec(`create table if not exists price_history 
(coin string, price string, timestamp datetime default current_timestamp)`); err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}
}

func RecordPrice(coin string, price string) error {
	_, err := db.Exec(fmt.Sprintf(`insert into price_history (coin, price) 
select '%s', '%s' where not exists 
(select timestamp from price_history where coin = '%s' and timestamp > 
(select datetime('now', '-1 hour')))`, coin, price, coin))
	return err
}

func GetPriceHistory(coin string, asc bool) ([][]string, error) {
	var allValues = [][]string{}
	limit := 10

	order := "asc"
	if !asc {
		order = "desc"
	}

	rows, err := db.Query(fmt.Sprintf(`select coin, price, datetime(timestamp, "localtime")
from price_history where coin = '%s' order by timestamp %s limit %d`, coin, order, limit))
	if err != nil {
		return allValues, err
	}
	defer rows.Close()

	for rows.Next() {
		var values = []string{"", "", ""}
		if err := rows.Scan(&values[0], &values[1], &values[2]); err != nil {
			return allValues, err
		}
		allValues = append(allValues, values)
	}
	err = rows.Err()

	return allValues, err
}

func isThereOldValues(coin string) (bool, error) {
	// TODO: hardcoded
	limit := 10
	query := fmt.Sprintf(`select case when 
(select count(*) from price_history where coin="%s") > %d 
then "yes" else "no"
end;`, coin, limit)

	rows, err := db.Query(query)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	rows.Next()

	var result string
	if err := rows.Scan(&result); err != nil {
		return false, err
	}

	return result == "yes", nil
}

func DeleteOlder(coin string) error {
	oldValues, err := isThereOldValues(coin)
	if err != nil {
		return err
	}
	if !oldValues {
		return nil
	}

	// TODO only 1, because it runs everytime
	query := fmt.Sprintf(`delete from price_history 
where timestamp in 
(select timestamp from price_history 
where coin="%s" order by timestamp asc limit 1)`, coin)

	_, err = db.Exec(query)
	return err
}

func Close() {
	db.Close()
}
