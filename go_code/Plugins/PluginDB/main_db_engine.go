package PluginKeyValueDBStore

import (
	//"database/sql"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/golang/protobuf/ptypes"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/pluginDBgRPCApi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"jlambert/FenixInception2/go_code/common_code"
	"time"
)

var mainDB *sqlx.DB

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "pluginkeyvaluestore"
)

// **********************************************************************************************************
// Initiate connection to Main Database for KeyValueStore
//
func initiateMainDB() {

	logger.WithFields(logrus.Fields{
		"ID": "8be4a508-9584-4ff3-8d97-028591ece5d4",
	}).Debug("Initiate connection to Main Database for Fenix Inception")

	var err error

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	//mainDB, err = sql.Open("postgres", dbinfo)
	mainDB, err = sqlx.Connect("postgres", dbinfo)

	// If some error occured then close the application
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":  "8be4a508-9584-4ff3-8d97-028591ece5d4",
			"err": err,
		}).Fatal("Error when initiateting connection to Main Database for Fenix Inception, closing Fenix Inception")
	}
}

// **********************************************************************************************************
// Close connection connection to Main Database for KeyValueStore
//
func closeMainDB() {

	logger.WithFields(logrus.Fields{
		"ID": "8be4a508-9584-4ff3-8d97-028591ece5d4",
	}).Debug("Close connection to Main Database for KeyValueStore")

	var err error

	err = mainDB.Close()

	// some error occured when closing connection to Database
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ID":  "faf440bc-d94c-4521-8cd1-193e15874bdf",
			"err": err,
		}).Error("Error when closing connection to Main Database for KeyValueStore")
	}
}

// **********************************************************************************************************
// Read ValueResponseMessage from Plugin-Key-Value-Store-DB
//
func readKeyValuetMessageFromDB(readKeyRequestMessage gRPC.ReadKeyRequestMessage) (valueResponseMessage *gRPC.ValueResponseMessage, err error) {

	// Prepare SQL
	var sqlToBeExecuted = "SELECT Key, Bucket, ValueSaveTypeId, ValueSaveTypeName, Value, ValueString, updatedDateTime "
	sqlToBeExecuted = sqlToBeExecuted + "FROM keyvaluestore.keyvaluestore "
	sqlToBeExecuted = sqlToBeExecuted + "WHERE Key = $1 "
	sqlToBeExecuted = sqlToBeExecuted + "ORDER BY updatedDateTime ASC "

	//sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)

	if err != nil {
		// Error while preparing SQL
		logger.WithFields(logrus.Fields{
			"ID":                    "6e7e206b-7966-47d7-9efa-ba5674abb832",
			"err":                   err,
			"readKeyRequestMessage": readKeyRequestMessage,
			"sqlToBeExecuted":       sqlToBeExecuted,
		}).Error("Error when Praparing SQL for getting Value for Key in Plugin KeyValue-storeDB")

	} else {
		// Preparing SQL was OK
		// Now Execute SQL in DB
		sqlResult, err := mainDB.Queryx(sqlToBeExecuted, readKeyRequestMessage.Key)
		if err != nil {
			// Error while executing SQL
			logger.WithFields(logrus.Fields{
				"ID":        "",
				"err":       err,
				"sqlResult": sqlResult,
			}).Error("Error when executing SQL for getting Value for Key in Plugin KeyValue-storeDB")
		} else {
			// Executing SQL was OK
			logger.WithFields(logrus.Fields{
				"ID":        "8102b4d5-c2ed-493e-918a-3365669d9cba",
				"sqlResult": sqlResult,
			}).Debug("sqlResult for getting Value for Key in Plugin KeyValue-storeDB was OK")

			// Convert SQL-result into correct return message for function
			//var testInstructionId string
			var numberOfRowsInResult int
			type savedKeyValueStruct struct {
				Key               string    `db:"key"`
				Bucket            string    `db:"bucket"`
				ValueSaveTypeId   int32     `db:"valuesavetypeid"`
				ValueSaveTypeName string    `db:"valuesavetypename"`
				Value             []byte    `db:"value"`
				ValueString       string    `db:"valuestring"`
				UpdatedDateTime   time.Time `db:"updateddatetime"`
			}

			var savedKeyValue savedKeyValueStruct

			numberOfRowsInResult = 0
			for sqlResult.Next() {
				// Only get the latest object, if there are more rows
				if numberOfRowsInResult == 0 {
					err = sqlResult.StructScan(&savedKeyValue)
					if err != nil {
						// Error while looping through result
						logger.WithFields(logrus.Fields{
							"ID":        "2006d78e-44c4-417c-92a6-9b749e48b946",
							"err":       err,
							"sqlResult": sqlResult,
						}).Error("Error when processing SQL-result for for TestInstruction-Payload for getting TestInstructions-Payload to be sent for execution")
					}
				} else {
					numberOfRowsInResult++
				}
			}

			// If the SQL gave a result that differs from 1 row then something is wrong
			if numberOfRowsInResult != 0 {
				// Eimport "github.com/golang/protobuf/ptypes"rror while looping through result
				logger.WithFields(logrus.Fields{
					"ID":                   "bf187ebf-ddbd-49f5-ac07-0be1e3251985",
					"numberOfRowsInResult": numberOfRowsInResult,
				}).Error("Error, expected exactly one row in result set")

				// Create an error
				err = errors.New("Error, expected exactly one row in result set but got " + string(numberOfRowsInResult))
			} else {
				// Transfer result into correct return-structure
				protoTimeStamp, err := ptypes.TimestampProto(savedKeyValue.UpdatedDateTime)

				if err != nil {
					// Could not convert from Proto Time stamp
					logger.WithFields(logrus.Fields{
						"ID":            "97a45e34-ed25-4907-a2ab-e6e6e752e527",
						"err":           err,
						"savedKeyValue": savedKeyValue,
					}).Error("Could not convert from Proto Time stamp after reading KeyValue Store")
					valueResponseMessage = &gRPC.ValueResponseMessage{
						Key:             savedKeyValue.Key,
						Bucket:          savedKeyValue.Key,
						ValueSaveType:   gRPC.ValueSaveTypeEnum(savedKeyValue.ValueSaveTypeId), //(gRPC.CurrentVersionEnum).gR [savedKeyValue.ValueSaveTypeId]
						Value:           savedKeyValue.Value,
						ValueString:     savedKeyValue.ValueString,
						UpdatedDateTime: nil,
						Acknack:         false,
						Comments:        "Could not convert from Proto Time stamp after reading KeyValue Store",
					}
				} else {

					/*
						//TODO Remove the code below
						tempTime, _ := ptypes.Timestamp(protoTimeStamp)
						fmt.Println(tempTime)
						fmt.Println(savedKeyValue.UpdatedDateTime.In(time.Local))
						fmt.Println(tempTime.String())
					*/

					// OK in convert from Proto Time stamp
					valueResponseMessage = &gRPC.ValueResponseMessage{
						Key:             savedKeyValue.Key,
						Bucket:          savedKeyValue.Bucket,
						ValueSaveType:   gRPC.ValueSaveTypeEnum(savedKeyValue.ValueSaveTypeId), //(gRPC.CurrentVersionEnum).gR [savedKeyValue.ValueSaveTypeId]
						Value:           savedKeyValue.Value,
						ValueString:     savedKeyValue.ValueString,
						UpdatedDateTime: protoTimeStamp,
						Acknack:         true,
						Comments:        "Readable  UTC timestamp: " + ptypes.TimestampString(protoTimeStamp),
					}
				}
			}
		}
	}

	// Return the saved value for the key
	return valueResponseMessage, err
}

// **********************************************************************************************************
// Save incoming 'WriteKeyValueMessage' to Plugin KeyValue-DB-store
//
func saveKeyValuetMessageInDB(writeKeyValueMessage *gRPC.WriteKeyValueMessage) (messageSavedInDB bool) {

	messageSavedInDB = true

	// Prepare SQL
	var sqlToBeExecuted = "INSERT INTO keyvaluestore.keyvaluestore "
	sqlToBeExecuted = sqlToBeExecuted + "(Key, Bucket, ValueSaveTypeId, ValueSaveTypeName, Value, ValueString, updatedDateTime) "
	sqlToBeExecuted = sqlToBeExecuted + "VALUES($1,$2,$3,$4,$5,$6,$7) "

	sqlStatement, err := mainDB.Prepare(sqlToBeExecuted)
	if err != nil {
		// Execute SQL in DB
		logger.WithFields(logrus.Fields{
			"ID":                   "38fffbaf-ea07-4c77-aa6f-41b827c9767d",
			"err":                  err,
			"writeKeyValueMessage": writeKeyValueMessage,
		}).Error("Error when Praparing SQL for updating Plugin Key-Value storeDB with data from 'writeKeyValueMessage'")

		messageSavedInDB = false
	} else {
		//SQL prepared OK
		/*
			create table plugins.KeyValueStore
			(
			    Key               uuid                     not null, -- The Key for the saved data -- 1
			    Bucket            varchar default null,              -- A type definition to be able to group content -- 2
			    ValueSaveTypeId   int                      not null, -- Defines how the value is saved, As ByteArray, String or as both of them -- 3
			    ValueSaveTypeName varchar                  not null, -- Defines how the value is saved, As ByteArray, String or as both of them -- 3
			    Value             bytea   default null,              -- The value saved as bytearray; -- 4
			    ValueString       varchar default null,              -- The value as plain text; -- 5
			    updatedDateTime   timestamp with time zone not null  -- The Datetime when the row was created/updates''; -- 6
			);

		*/
		// Values to insert into database
		sqlResult, err := sqlStatement.Exec(
			writeKeyValueMessage.Key,
			writeKeyValueMessage.Bucket,
			writeKeyValueMessage.ValueSaveType,
			gRPC.CurrentVersionEnum_name[int32(writeKeyValueMessage.ValueSaveType)],
			writeKeyValueMessage.Value,
			writeKeyValueMessage.ValueString,
			common_code.GeneraTimeStampUTC())

		if err != nil {
			// Error while executing
			logger.WithFields(logrus.Fields{
				"ID":                   "984e10e0-de20-4440-a63a-8bf444ccf831",
				"err":                  err,
				"sqlResult":            sqlResult,
				"writeKeyValueMessage": writeKeyValueMessage,
			}).Error("Error when updating Plugin Key-Value storeDB with data from 'writeKeyValueMessage'")

			messageSavedInDB = false
		} else {
			//SQL executed OK
			logger.WithFields(logrus.Fields{
				"ID":                   "7c54131a-eeff-41c4-951c-71370e49588b",
				"err":                  err,
				"sqlResult":            sqlResult,
				"writeKeyValueMessage": writeKeyValueMessage,
			}).Debug("Plugin Key-Value storeDB was updated with data from 'writeKeyValueMessage'")
		}
	}

	// Return if the messages was saved or not in database
	return messageSavedInDB
}
