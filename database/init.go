package database

func CreateDatabaseConnections() {
	MongoConn = getMongoConnection()
}
