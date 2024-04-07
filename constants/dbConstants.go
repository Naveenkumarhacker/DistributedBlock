package constants

import "gitlab.com/avarf/getenvs"

var SqliteDbFileLocation = getenvs.GetEnvString("INTERNAL_DB_FILE_LOCATION", "sqlite.db")

var (
	DbUser     = getenvs.GetEnvString("DB_USERNAME", "user")
	DbPassword = getenvs.GetEnvString("DB_PASSWORD", "password")
	DbName     = getenvs.GetEnvString("DB_NAME", "blocks")
	DbHost     = getenvs.GetEnvString("DB_HOST", "127.0.0.1")
	DbPort     = getenvs.GetEnvString("DB_PORT", "3306")
)
