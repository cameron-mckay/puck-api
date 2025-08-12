package db

import (
	"context"
	"database/sql"

	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB

func Init(connString string) error {
	d, err := sql.Open("sqlserver", connString)

	db = d

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	if err != nil {
		return err
	}

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
	if db == nil {
		panic("db instance not initialized")
	}
	_, err := db.Exec("EXEC BulkAddRPAPuckToNetwork @FulfillmentBinID = @p1, @NetworkID = @p2", sql.Named("p1", binId), sql.Named("p2", networkId))
	return err
}

func DeleteAllSensorsOnNetwork(networkId int) error {
	if db == nil {
		panic("db instance not initialized")
	}
	_, err := db.Exec("EXEC BulkDeleteRPAPuckFromNetwork @NetworkID = @p1", sql.Named("p1", networkId))
	return err
}

type MessageCounts struct {
	SensorID   int  `json:"sensorId"`
	Count      int  `json:"count"`
	MinBattery int  `json:"minBattery"`
	IsDirty    bool `json:"isDirty"`
}

func GetMessageCounts(networkId int) ([]MessageCounts, error) {
	var mcs []MessageCounts

	if db == nil {
		panic("db instance not initialized")
	}
	rows, err := db.Query("EXEC BulkRPAPuckPastDay @NetworkID = @p1", sql.Named("p1", networkId))

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var mc MessageCounts
		err = rows.Scan(&mc.SensorID, &mc.Count, &mc.MinBattery, &mc.IsDirty)
		if err != nil {
			return nil, err
		}
		mcs = append(mcs, mc)
	}

	return mcs, nil
}
