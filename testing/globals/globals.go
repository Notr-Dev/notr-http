package globals

import "database/sql"

type DataType struct {
	DB *sql.DB
}

var Data = DataType{}
