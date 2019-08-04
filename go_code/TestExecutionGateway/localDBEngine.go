package TestExecutionGateway

import (
	"errors"
	"github.com/jlambert68/FenixInception/go_code/common_code"
	gRPC "github.com/jlambert68/FenixInception/go_code/common_code/Gateway_gRPC_api"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

// ********************************************************************************************
// Open the database file in same directory as the executable file
//
func initiateDB(localDBFile string) {
	var boltDBName = "bolt.db"
	var boltDBNameUsed string

	if localDBFile == "" {
		// No database name attached then use predined name
		boltDBNameUsed = boltDBName
	} else {
		// Database name atteched
		boltDBNameUsed = localDBFile
	}

	dbRef, err := bolt.Open(boltDBNameUsed, 0644, nil)
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"ID":    "e1dbd9be-e790-41a5-9a54-a5cc1952219f",
			"error": err,
		}).Fatal("Error when trying to open databse: '" + boltDBNameUsed + "'")
	} else {
		Logger.WithFields(logrus.Fields{
			"ID": "c702a65c-5417-471b-a00a-1055e864e8e0",
		}).Debug("Database was initiated with name: '" + boltDBNameUsed + "'")
	}

	// If no errors then save reference to DB in gateway object
	db = dbRef

	// Start Database Engine as a go-routine
	go databaseEngine()
}

// ********************************************************************************************
// Close the database file
//
func closeDB() {

	err := db.Close()
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"ID":  "bfee7b43-d719-45b2-a099-54e09b53b3ab",
			"err": err,
		}).Error("Error when closing local database")
	} else {
		Logger.WithFields(logrus.Fields{
			"ID": "377091eb-f4a0-42a0-9d25-46b50a0862ec",
		}).Info("SUccess in closing local database")
	}
}

// ********************************************************************************************
// Do Reads and Writes to local database
//
func databaseEngine() {

	var err error

	Logger.WithFields(logrus.Fields{
		"ID":    "d6b7454e-eb99-4c7d-9ec9-84249a7ee848",
		"error": err,
	}).Debug("Database engine started")

	// Used for Debugging
	defer Logger.WithFields(logrus.Fields{
		"ID":    "ea3e5c3f-0c27-4303-80c5-7d76b875d03b",
		"error": err,
	}).Debug("Exiting database engine with 'defer'")

	for {
		// Wait for data comes from channel to dtabase engine
		messageToDbEngine := <-DbMessageQueue
		Logger.WithFields(logrus.Fields{
			"ID":                             "5bdb83d8-e913-4933-969b-5035f41e4a70",
			"messageToDbEngine.messageType,": messageToDbEngine.MessageType,
			"messageToDbEngine":              messageToDbEngine,
		}).Debug("Received a new message to Database engine")

		// Decide if it's a Read- or Write-instruction
		switch messageToDbEngine.MessageType {
		case common_code.DbRead:
			// Infor entering this part in Debug-mode
			Logger.WithFields(logrus.Fields{
				"ID": "993dc086-84ec-4c45-a9b1-e7c73ea50b50",
			}).Debug("Entering Read-Database")

			// Read data from Database and send back using incoming return-channel
			err = db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte(messageToDbEngine.Bucket))
				if bucket == nil {
					Logger.WithFields(logrus.Fields{
						"ID":     "f4093818-80ee-48a2-aa47-3fb7a0792045",
						"err":    err,
						"Bucket": messageToDbEngine.Bucket,
					}).Warning("Bucket not found")

					err = errors.New("bucket not found")
					// Send back err and empty value using attached channel
					readResultMessage := common_code.DbResultMessageStruct{
						err,
						messageToDbEngine.Bucket,
						[]byte("")}

					messageToDbEngine.ResultsQueue <- readResultMessage

					return nil

				} else {
					Logger.WithFields(logrus.Fields{
						"ID":     "8ccfac34-90a0-4b50-8e37-bc4f10d76f62",
						"Bucket": bucket,
					}).Debug("Success in finding Bucket")

				}
				// Retrieve value from key
				value := bucket.Get([]byte(messageToDbEngine.Key))
				valueString := string(value)
				Logger.WithFields(logrus.Fields{
					"ID":     "4fa038d7-c135-41bb-804a-7a4d249e1bc9",
					"Bucket": bucket,
					"Key":    messageToDbEngine.Key,
					"Value":  valueString,
				}).Debug("Success in reading Key")

				// Send back value using attached channel
				readResultMessage := common_code.DbResultMessageStruct{
					err,
					messageToDbEngine.Key,
					value}
				messageToDbEngine.ResultsQueue <- readResultMessage

				return nil
			})

		case common_code.DbWrite:
			// Infor entering this part in Debug-mode
			Logger.WithFields(logrus.Fields{
				"ID": "5173bc5a-53f2-4a9a-ade6-97e0e875478a",
			}).Debug("Entering Write-Database")

			// Store incoming data in defined bucket
			err = db.Update(func(tx *bolt.Tx) error {
				// Create Bucket if it not exist
				bucket, err := tx.CreateBucketIfNotExists([]byte(messageToDbEngine.Bucket))
				if err != nil {
					Logger.WithFields(logrus.Fields{
						"ID":     "044b668e-3762-4b4c-98c3-0e7dae7e7fda",
						"err":    err,
						"Bucket": bucket,
					}).Error("Error when creating bucket")

					// Send back err and empty value using attached channel
					readResultMessage := common_code.DbResultMessageStruct{
						err,
						messageToDbEngine.Bucket,
						[]byte("")}

					messageToDbEngine.ResultsQueue <- readResultMessage

					return nil

				} else {
					Logger.WithFields(logrus.Fields{
						"ID":     "e0359bee-de08-420f-b417-9635fc7b1e9b",
						"Bucket": bucket,
					}).Debug("Success in creating Bucket")
				}

				// Save Value using its key
				err = bucket.Put(
					[]byte(messageToDbEngine.Key),
					[]byte(messageToDbEngine.Value))
				if err != nil {
					if err != nil {
						Logger.WithFields(logrus.Fields{
							"ID":     "48ca41de-cb14-44dd-902d-fb146bd0a9fa",
							"err":    err,
							"Bucket": bucket,
							"Key":    messageToDbEngine.Key,
							"Value":  messageToDbEngine.Value,
						}).Error("Error when saving Key-Value in bucket")
					} else {
						Logger.WithFields(logrus.Fields{
							"ID":     "b5ac3b38-9aa6-4b42-b289-94c722a2dbbe",
							"Bucket": bucket,
							"Key":    messageToDbEngine.Key,
							"Value":  messageToDbEngine.Value,
						}).Debug("Success in saving Key-Value in bucket")
					}
				}

				// Send back value using attached channel
				dbWritedResultMessage := common_code.DbResultMessageStruct{
					err,
					messageToDbEngine.Key,
					nil}
				messageToDbEngine.ResultsQueue <- dbWritedResultMessage

				return nil
			})

		// No need to take care of error from return due to it is always nil

		case common_code.DBGetFirstObjectFromBucket:
			// Info entering this part in Debug-mode
			Logger.WithFields(logrus.Fields{
				"ID": "becfe123-4ff4-41f5-bcf0-26e6d70fe176",
			}).Debug("Entering Get-First-Object-In-Bucket-Database")

			// Read data from Database and send back using incoming return-channel
			err = db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte(messageToDbEngine.Bucket))
				if bucket == nil {
					Logger.WithFields(logrus.Fields{
						"ID":     "7871629e-d622-463f-a526-1ee82a80ed97",
						"err":    err,
						"Bucket": messageToDbEngine.Bucket,
					}).Warning("Bucket not found")

					err = errors.New("bucket not found")

					// Send back err and empty value using attached channel
					readResultMessage := common_code.DbResultMessageStruct{
						err,
						messageToDbEngine.Bucket,
						[]byte("")}

					messageToDbEngine.ResultsQueue <- readResultMessage

					return nil

				} else {
					Logger.WithFields(logrus.Fields{
						"ID":     "6d53f3f1-1d01-4d92-a012-e8667eae8aae",
						"Bucket": bucket,
					}).Debug("Success in finding Bucket")

					// Retrieve key of first object
					// Create a cursor for iteration.
					cursor := bucket.Cursor()

					// Get key for first object in Bucket
					key, _ := cursor.First()

					// Send back value using attached channel
					dbGetFirstObjectResultMessage := common_code.DbResultMessageStruct{
						err,
						string(key),
						nil}

					messageToDbEngine.ResultsQueue <- dbGetFirstObjectResultMessage

					return nil
				}
			})

		case common_code.DBDelete:
			// Info entering this part in Debug-mode
			Logger.WithFields(logrus.Fields{
				"ID": "45b6fbe7-0c73-48ba-9267-7a17ed677867",
			}).Debug("Entering Delete Object in Bucket-Database")

			// Delete the key in a write transaction.
			err := db.Update(func(tx *bolt.Tx) error {
				return tx.Bucket([]byte(messageToDbEngine.Bucket)).Delete([]byte(messageToDbEngine.Key))
			})

			if err != nil {
				LogErrorAndSendInfoToFenix(
					"ab67dec4-9040-4f39-88ed-8c7dbb56698a",
					gRPC.InformationMessage_FATAL,
					"FATAL error when deleting object from bucket",
					"FATAL error when deleting object from bucket",
					"FATAL error when deleting key: '"+messageToDbEngine.Key+"' in bucket: '"+messageToDbEngine.Bucket+"'",
					"FATAL error when deleting key: '"+messageToDbEngine.Key+"' in bucket: '"+messageToDbEngine.Bucket+"'",
				)
			}

			// Send back value using attached channel
			dbDeleteResultMessage := common_code.DbResultMessageStruct{
				err,
				messageToDbEngine.Key,
				nil}
			messageToDbEngine.ResultsQueue <- dbDeleteResultMessage

		default:
			Logger.WithFields(logrus.Fields{
				"ID":                            "16c7f34d-507b-4d04-b765-9334648320cd",
				"messageToDbEngine.messageType": messageToDbEngine.MessageType,
			}).Warning("No known messageType sent to Database Enging")

			// Send back value using attached channel
			var errorMessage = errors.New("messageToDbEngine.messageType is not a known type")
			dbUnknownResultMessage := common_code.DbResultMessageStruct{
				errorMessage,
				messageToDbEngine.Key,
				nil}

			messageToDbEngine.ResultsQueue <- dbUnknownResultMessage

		}
	}
}
