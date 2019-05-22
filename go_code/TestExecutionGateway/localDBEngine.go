package TestExecutionGateway

import (
	"errors"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

// ********************************************************************************************
// Open the database file in same directory as the executable file
//
func (gatewayObject *GatewayTowardsPluginObject_struct) initiateDB() {
	var boldDBName string = "bolt.db"

	db, err := bolt.Open(boldDBName, 0644, nil)
	if err != nil {
		gatewayObject.logger.WithFields(logrus.Fields{
			"ID":    "e1dbd9be-e790-41a5-9a54-a5cc1952219f",
			"error": err,
		}).Fatal("Error when trying to open databse: '" + boldDBName + "'")
	}

	// If no errors then save reference to DB in gateway object
	gatewayObject.db = db

	// Start Database Engine as a go-routine
	go gatewayObject.databaseEngine()
}

// ********************************************************************************************
// Do Reads and Writes to local database
//
func (gatewayObject *GatewayTowardsPluginObject_struct) databaseEngine() {

	var err error

	for {

		// Wait for data comes from channel to dtabase engine
		messageToDbEngine := <-gatewayObject.dbMessageQueue
		gatewayObject.logger.WithFields(logrus.Fields{
			"ID":                "5bdb83d8-e913-4933-969b-5035f41e4a70",
			"messageToDbEngine": messageToDbEngine,
		}).Info("Received a new message to Database engine")

		// Decide if it's a Read- or Write-instruction
		switch messageToDbEngine.messageType {
		case DB_READ:
			// Read data from Database and send back using incoming return-channel
			err = gatewayObject.db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte(messageToDbEngine.bucket))
				if bucket == nil {
					gatewayObject.logger.WithFields(logrus.Fields{
						"ID":     "f4093818-80ee-48a2-aa47-3fb7a0792045",
						"err":    err,
						"Bucket": bucket,
					}).Warning("Bucket not found")

				} else {
					gatewayObject.logger.WithFields(logrus.Fields{
						"ID":     "8ccfac34-90a0-4b50-8e37-bc4f10d76f62",
						"Bucket": bucket,
					}).Info("Success in finding Bucket")
				}

				// Retrieve value from key
				value := bucket.Get([]byte(messageToDbEngine.key))
				valueString := string(value)
				gatewayObject.logger.WithFields(logrus.Fields{
					"ID":     "4fa038d7-c135-41bb-804a-7a4d249e1bc9",
					"Bucket": bucket,
					"Key":    messageToDbEngine.key,
					"Value":  valueString,
				}).Info("Success in reading Key")

				// Send back value using attached channel
				readResultMessage := dbResultMessage_struct{
					err,
					messageToDbEngine.key,
					value}
				messageToDbEngine.resultsQueue <- readResultMessage

				return nil
			})

		case DB_WRITE:
			// Store incoming data in defined bucket
			err = gatewayObject.db.Update(func(tx *bolt.Tx) error {
				// Create Bucket if it not exist
				bucket, err := tx.CreateBucketIfNotExists([]byte(messageToDbEngine.bucket))
				if err != nil {
					gatewayObject.logger.WithFields(logrus.Fields{
						"ID":     "044b668e-3762-4b4c-98c3-0e7dae7e7fda",
						"err":    err,
						"Bucket": bucket,
					}).Error("Error when creating bucket")
				} else {
					gatewayObject.logger.WithFields(logrus.Fields{
						"ID":     "e0359bee-de08-420f-b417-9635fc7b1e9b",
						"Bucket": bucket,
					}).Info("Success in creating Bucket")
				}

				// Save Value using its key
				err = bucket.Put(
					[]byte(messageToDbEngine.key),
					[]byte(messageToDbEngine.value))
				if err != nil {
					if err != nil {
						gatewayObject.logger.WithFields(logrus.Fields{
							"ID":     "48ca41de-cb14-44dd-902d-fb146bd0a9fa",
							"err":    err,
							"Bucket": bucket,
							"Key":    messageToDbEngine.key,
							"Value":  messageToDbEngine.value,
						}).Error("Error when saving Key-Value in bucket")
					} else {
						gatewayObject.logger.WithFields(logrus.Fields{
							"ID":     "b5ac3b38-9aa6-4b42-b289-94c722a2dbbe",
							"Bucket": bucket,
							"Key":    messageToDbEngine.key,
							"Value":  messageToDbEngine.value,
						}).Info("Success in saving Key-Value in bucket")
					}
				}

				// Send back value using attached channel
				dbWritedResultMessage := dbResultMessage_struct{
					err,
					messageToDbEngine.key,
					nil}
				messageToDbEngine.resultsQueue <- dbWritedResultMessage

				return nil
			})

			// No need to take care of error from return due to it is always nil

		default:
			gatewayObject.logger.WithFields(logrus.Fields{
				"ID":                            "16c7f34d-507b-4d04-b765-9334648320cd",
				"messageToDbEngine.messageType": messageToDbEngine.messageType,
			}).Warning("No known messageType sent to Database Enging")

			// Send back value using attached channel
			var errorMessage = errors.New("messageToDbEngine.messageType is not a known type")
			dbUnknownResultMessage := dbResultMessage_struct{
				errorMessage,
				messageToDbEngine.key,
				nil}

			messageToDbEngine.resultsQueue <- dbUnknownResultMessage

		}
	}
}
