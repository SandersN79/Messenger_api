package configuration

import (
	"encoding/json"
	"os"
)

// Configuration is a struct designed to hold the applications variable configuration settings
type Configuration struct {
	MongoURI string
	Secret   string
	Database string
	MasterAdminUsername string
	MasterAdminInitialPassword string
	Registration  string
	Host string
	Port string
	DBName string
}

// ConfigurationSettings is a function that reads a json configuration file and outputs a Configuration struct
func ConfigurationSettings(env string) *Configuration {
	confFile := "tsconfig.json"
	if env == "test" {
		confFile = "test_conf.json"
	}
	file, _ 				:= os.Open(confFile)
	decoder 				:= json.NewDecoder(file)
	configurationSettings 	:= Configuration{}
	err 					:= decoder.Decode(&configurationSettings)
	if err != nil {
		panic(err)
	}
	return &configurationSettings
}

// InitializeEnvironmentals
func (c *Configuration) InitializeEnvironmentals() {
	os.Setenv("MONGO_URI", c.MongoURI)
	os.Setenv("SECRET", c.Secret)
	os.Setenv("DATABASE", c.Database)
	os.Setenv("MASTER_ADMIN_USERNAME", c.MasterAdminUsername)
	os.Setenv("MASTER_ADMIN_INITIAL_PASSWORD", c.MasterAdminInitialPassword)
	os.Setenv("REGISTRATION", c.Registration)
	os.Setenv("HOST", c.Host)
	os.Setenv("PORT", c.Port)
	os.Setenv("DBNAME", c.DBName)
}
