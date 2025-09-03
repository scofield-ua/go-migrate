package config

type Config struct {
	DB db
}

type InitConfigParams struct {
	driver     DbDriver
	dbHost     string
	dbUsername string
	dbPassword string
	dbDatabase string
}

func (c *Config) Init(params InitConfigParams) {
	c.DB.SetHost(params.dbHost)
	c.DB.SetUsername(params.dbUsername)
	c.DB.SetPassword(params.dbPassword)
	c.DB.SetDbName(params.dbDatabase)
	c.DB.SetDriver(params.driver)
}
