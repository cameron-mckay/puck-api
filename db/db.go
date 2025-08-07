package db

import (
	"context"
	"database/sql"

	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB

func Init(connString string) error {
	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		return err
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func Close() {
	db.Close()
}

func AddBinToNetwork(binId int, networkId int) error {
	_, err := db.Exec("EXEC BulkAddRPAPuckToNetwork @FullfillmentBinID = ?, @NetworkID = ?", binId, networkId)
	return err
}

func DeleteAllSensorsOnNetwork(networkId int) error {
	_, err := db.Exec("EXEC BulkDeleteRPAPuckFromNetwork @NetworkID = ?", networkId)
	return err
}
