package config

import "os"

func Config() map[string]string {
	dbName, exists := os.LookupEnv("DB")
	connString, connExists := os.LookupEnv("CS")

	if exists && connExists {
		return map[string]string{
			"db": dbName,
			"cs": connString,
		}
	}
	dbName = os.Getenv("DB")
	connString = os.Getenv("CS")
	return map[string]string{
		"db": dbName,
		"cs": connString,
	}
}
