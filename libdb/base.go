package libdb

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/topher200/forty-thieves/libunix"
)

type InsertResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (ir *InsertResult) LastInsertId() (int64, error) {
	return ir.lastInsertId, nil
}

func (ir *InsertResult) RowsAffected() (int64, error) {
	return ir.rowsAffected, nil
}

func newDbForTest(t *testing.T) *sqlx.DB {
	var err error

	pgdsn := os.Getenv("DSN")
	if pgdsn == "" {
		pguser, _, pghost, pgport, pgsslmode := os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGSSLMODE")
		if pguser == "" {
			pguser, err = libunix.CurrentUser()
			if err != nil {
				t.Fatalf("Getting current user should never fail. Error: %v", err)
			}
		}

		if pghost == "" {
			pghost = "localhost"
		}

		if pgport == "" {
			pgport = "5432"
		}

		if pgsslmode == "" {
			pgsslmode = "disable"
		}

		pgdsn = fmt.Sprintf("postgres://%v@%v:%v/forty_thieves_test?sslmode=%v", pguser, pghost, pgport, pgsslmode)
	}

	db, err := sqlx.Connect("postgres", pgdsn)
	if err != nil {
		t.Fatalf("Connecting to local postgres should never fail. Error: %v", err)
	}
	return db
}

func newBaseDBForTest(t *testing.T) *Base {
	base := &Base{}
	base.db = newDbForTest(t)

	return base
}

type Base struct {
	db    *sqlx.DB
	table string
	hasID bool
}

func (b *Base) newTransactionIfNeeded(tx *sqlx.Tx) (*sqlx.Tx, bool, error) {
	var err error
	wrapInSingleTransaction := false

	if tx != nil {
		return tx, wrapInSingleTransaction, nil
	}

	tx, err = b.db.Beginx()
	if err == nil {
		wrapInSingleTransaction = true
	}

	if err != nil {
		return nil, wrapInSingleTransaction, err
	}

	return tx, wrapInSingleTransaction, nil
}

func (b *Base) InsertIntoTable(tx *sqlx.Tx, data map[string]interface{}) (sql.Result, error) {
	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}

	keys := make([]string, 0)
	dollarMarks := make([]string, 0)
	values := make([]interface{}, 0)

	loopCounter := 1
	for key, value := range data {
		keys = append(keys, key)
		dollarMarks = append(dollarMarks, fmt.Sprintf("$%v", loopCounter))
		values = append(values, value)

		loopCounter++
	}

	var query string
	if loopCounter > 1 {
		query = fmt.Sprintf(
			"INSERT INTO %v (%v) VALUES (%v)",
			b.table,
			strings.Join(keys, ","),
			strings.Join(dollarMarks, ","))
	} else {
		// we're inserting a row with no data
		query = fmt.Sprintf(
			"INSERT INTO %v VALUES (default)",
			b.table,
		)
	}

	result := &InsertResult{}
	result.rowsAffected = 1

	if b.hasID {
		query = query + " RETURNING id"

		var lastInsertId int64
		err = tx.QueryRow(query, values...).Scan(&lastInsertId)
		if err != nil {
			if wrapInSingleTransaction == true {
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					return nil, rollbackErr
				}
			}
			return nil, err
		}

		result.lastInsertId = lastInsertId
	} else {
		_, err := tx.Exec(query, values...)
		if err != nil {
			if wrapInSingleTransaction == true {
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					return nil, rollbackErr
				}
			}
			return nil, err
		}
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}
	if err != nil {
		return nil, err
	}

	return result, err
}

func (b *Base) UpdateFromTable(tx *sqlx.Tx, data map[string]interface{}, where string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}

	keysWithDollarMarks := make([]string, 0)
	values := make([]interface{}, 0)

	loopCounter := 1
	for key, value := range data {
		keysWithDollarMark := fmt.Sprintf("%v=$%v", key, loopCounter)
		keysWithDollarMarks = append(keysWithDollarMarks, keysWithDollarMark)
		values = append(values, value)

		loopCounter++
	}

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v",
		b.table,
		strings.Join(keysWithDollarMarks, ","),
		where)

	result, err = tx.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}
	if err != nil {
		return nil, err
	}

	return result, err
}

func (b *Base) UpdateById(tx *sqlx.Tx, data map[string]interface{}, id int64) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}

	keysWithDollarMarks := make([]string, 0)
	values := make([]interface{}, 0)

	loopCounter := 1
	for key, value := range data {
		keysWithDollarMark := fmt.Sprintf("%v=$%v", key, loopCounter)
		keysWithDollarMarks = append(keysWithDollarMarks, keysWithDollarMark)
		values = append(values, value)

		loopCounter++
	}

	// Add id as part of values
	values = append(values, id)

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE id=$%v",
		b.table,
		strings.Join(keysWithDollarMarks, ","),
		loopCounter)

	result, err = tx.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}
	if err != nil {
		return nil, err
	}

	return result, err
}

func (b *Base) UpdateByKeyValueString(tx *sqlx.Tx, data map[string]interface{}, key, value string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}

	keysWithDollarMarks := make([]string, 0)
	values := make([]interface{}, 0)

	loopCounter := 1
	for key, value := range data {
		keysWithDollarMark := fmt.Sprintf("%v=$%v", key, loopCounter)
		keysWithDollarMarks = append(keysWithDollarMarks, keysWithDollarMark)
		values = append(values, value)

		loopCounter++
	}

	// Add value as part of values
	values = append(values, value)

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v=$%v",
		b.table,
		strings.Join(keysWithDollarMarks, ","),
		key,
		loopCounter)

	result, err = tx.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}
	if err != nil {
		return nil, err
	}

	return result, err
}

func (b *Base) DeleteFromTable(tx *sqlx.Tx, where string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}

	query := fmt.Sprintf("DELETE FROM %v", b.table)

	if where != "" {
		query = query + " WHERE " + where
	}

	result, err = tx.Exec(query)
	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}
	if err != nil {
		return nil, err
	}

	return result, err
}

func (b *Base) DeleteById(tx *sqlx.Tx, id int64) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}

	query := fmt.Sprintf("DELETE FROM %v WHERE id=$1", b.table)

	result, err = tx.Exec(query, id)
	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}
	if err != nil {
		return nil, err
	}

	return result, err
}
