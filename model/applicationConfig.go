package model

type ApplicationConfig struct {
	Base BaseApplicationConfig
	Prod EnvApplicationConfig
	Coll EnvApplicationConfig
}
