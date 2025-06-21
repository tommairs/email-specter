package database

func init() {
	MongoConn = getMongoConnection()
}
