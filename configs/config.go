package configs

import (
	"os"
)

var CERT_PATH = getStringParameter("CERT_PATH", "./")
var KEY_PRIVATE = CERT_PATH + "problserv.key"
var CERTIFICATE = CERT_PATH + getStringParameter("PROBLEMS_CERT_NAME", "probllocal.crt")

var PG_HOST = getStringParameter("PG_HOST", "localhost")
var PG_PORT = getStringParameter("PG_PORT", "5444")
var POSTGRES_DB = getStringParameter("POSTGRES_DB", "scooterdb")
var POSTGRES_USER = getStringParameter("POSTGRES_USER", "scooteradmin")
var POSTGRES_PASSWORD = getStringParameter("POSTGRES_PASSWORD", "Megascooter!")
var GRPC_PORT = getStringParameter("PROBLEMS_GRPC_PORT", "4444")

var USERS_GRPC_PORT = getStringParameter("USERS_GRPC_PORT", "5555")
var USER_SERVICE = getStringParameter("USER_SERVICE", "localhost")

func getStringParameter(paramName, defaultValue string) string {
	result, ok := os.LookupEnv(paramName)
	if !ok {
		result = defaultValue
	}
	return result
}
