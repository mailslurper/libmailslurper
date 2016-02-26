package storage

/*
ConnectionInformation contains data necessary to establish a connection
to a database server.
*/
type ConnectionInformation struct {
	Address  string
	Port     int
	Database string
	UserName string
	Password string
	Filename string
}

/*
NewConnectionInformation returns a new ConnectionInformation structure with
the address and port filled in.
*/
func NewConnectionInformation(address string, port int) *ConnectionInformation {
	return &ConnectionInformation{
		Address: address,
		Port:    port,
	}
}

/*
SetDatabaseInformation fills in the name of a database to connect to, and the user
credentials necessary to do so
*/
func (information *ConnectionInformation) SetDatabaseInformation(database, userName, password string) {
	information.Database = database
	information.UserName = userName
	information.Password = password
}

/*
SetDatabaseFile sets the name of a file-base database. This is used for SQLite
*/
func (information *ConnectionInformation) SetDatabaseFile(filename string) {
	information.Filename = filename
}
