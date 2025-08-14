package database_type

type DXDatabaseType int64

const (
	UnknownDatabaseType DXDatabaseType = iota
	PostgreSQL
	MySQL
	Oracle
	SQLServer
	MariaDb
)

func (t DXDatabaseType) String() string {
	switch t {
	case PostgreSQL:
		return "postgres"
	case MySQL:
		return "mysql"
	case Oracle:
		return "oracle"
	case SQLServer:
		return "sqlserver"
	case MariaDb:
		return "mariadb"
	default:

		return "unknown"
	}
}

func (t DXDatabaseType) Driver() string {
	switch t {
	case PostgreSQL:
		return "postgres"
	case MySQL:
		return "mysql"
	case Oracle:
		return "oracle"
	case SQLServer:
		return "sqlserver"
	case MariaDb:
		return "mariadb"
	default:

		return "unknown"
	}
}

func StringToDXDatabaseType(v string) DXDatabaseType {
	switch v {
	case "postgres", "postgresql":
		return PostgreSQL
	case "mysql":
		return MySQL
	case "mariadb":
		return MariaDb
	case "oracle":
		return Oracle
	case "sqlserver":
		return SQLServer
	default:

		return UnknownDatabaseType
	}
}
