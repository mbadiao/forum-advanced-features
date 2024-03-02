package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	_"github.com/mattn/go-sqlite3"
)

func CreateTable() *sql.DB {
	_, errNofile := os.Stat("./internals/database/database.db")

	db, err := sql.Open("sqlite3", "./internals/database/database.db")
	if err != nil {
		log.Println(err.Error())
	}
	if errNofile != nil {
		sqlcode, err := os.ReadFile("./internals/database/table.sql")
		if err != nil {
			log.Println(err.Error())
		}
		_, err = db.Exec(string(sqlcode))

		if err != nil {
			log.Println(err.Error())
		}
	}
	return db
}

func GeneratePrepare(text string) string {
	nb := len(strings.Split(text, ","))
	a := strings.Repeat("?,", nb)
	return "(" + a[:len(a)-1] + ")"
}

func Insert(db *sql.DB, table, values string, data ...interface{}) {
	prepare := GeneratePrepare(values)
	Query := fmt.Sprintf("INSERT INTO %v %v values %v", table, values, prepare)
	insert, err := db.Prepare(Query)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = insert.Exec(data...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}



func Scan(db *sql.DB, request string, data Table) ([]Table, error) {
	stmt, err := db.Prepare(request)
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // Ensure the statement is closed

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure the rows are closed

	tableData := []Table{}

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("data must be a pointer to a struct that implements the Table interface")
	}

	for rows.Next() {
		dynamicType := reflect.New(reflect.TypeOf(data).Elem()).Interface().(Table)
		if err := dynamicType.ScanRows(rows); err != nil {
			return nil, err
		}
		tableData = append(tableData, dynamicType)
	}
	return tableData, nil
}
