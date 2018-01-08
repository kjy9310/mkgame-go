package main

import (
	"database/sql"
	"log"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
)
type HelloHandler struct {
	db *sql.DB
}
func initDB() HelloHandler {
// sql.DB 객체 생성
	db, err := sql.Open("mysql", "root:12726@tcp(127.0.0.1:3306)/game_card")
	if err != nil {
		log.Fatal(err)
	}
	return (HelloHandler{db: db})
}
func (h *HelloHandler) selectList(sql string, args []interface{}) (bool, []map[string]string) {
	var resultMap []map[string]string
	rows, _ := h.db.Query(sql,args)
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	rowcount:=0;
	for rows.Next() {
		for i, _ := range columns {
            valuePtrs[i] = &values[i]
        }
		err := rows.Scan(valuePtrs...)
		if err != nil{
			return false, nil
		}
		rowMap := make(map[string]string)
        for i, col := range columns {
			//var v interface{}
			b, ok := values[i].([]byte)
			if (ok) {
                //v = string(b)
				rowMap[col] = string(b)
            } else {
                //v = values[i]
            }
        }
		resultMap = append(resultMap, rowMap)
		rowcount+=1
    }
	return true, resultMap
}
func (h *HelloHandler) selectRow(sql string, args []interface{}, count int) (bool, []string) {
	stringValues := make([]string,count)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for i := range values {
        valuePtrs[i] = &values[i]
    }
	row := h.db.QueryRow(sql,args...)
	err := row.Scan(valuePtrs...)
	if err != nil{
		return false, []string{err.Error()}
	}
	for i := range values {
		//var v interface{}
		b, ok := values[i].([]byte)
		if (ok) {
            //v = string(b)
			stringValues[i] = string(b)
        } else {
            //v = values[i]
        }
    }
	return true, stringValues
}
func (h *HelloHandler) closeDB() {
	defer h.db.Close()
}

