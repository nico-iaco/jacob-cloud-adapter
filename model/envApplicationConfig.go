package model

type EnvApplicationConfig struct {
	Postgres DbConfig
	Db2      DbConfig
	Oracle   DbConfig
}
