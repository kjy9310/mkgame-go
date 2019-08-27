package mysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)
type DbHandler struct {
	db *sql.DB
}
var Con DbHandler
func init() {
	user:="root"
	password:=""
	address:="localhost"
	port:="3306"
	schema:="game"
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+address+":"+port+")/"+schema)
	if err != nil {
		log.Fatal(err)
	}
	Con = (DbHandler{db: db})
}
func (h *DbHandler) SelectList(sql string, args []interface{}) (bool, []map[string]string) {
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
func (h *DbHandler) SelectRow(sql string, args []interface{}, count int) (bool, []string) {
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
func (h *DbHandler) CloseDB() {
	defer h.db.Close()
}

