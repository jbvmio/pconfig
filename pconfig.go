package pconfig

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// PConfig is the config used for processors.
type PConfig struct {
	Environment         string
	Datacenter          string
	LogLevel            string
	LogDir              string
	GroupName           string
	MemberName          string
	DeleteGroupOnStart  bool
	InitialOffsetOldest bool
	KafkaBrokers        []string
	InputTopics         []string
	OutputTopics        []string
	InOutMap            map[string]string
	HTTP                HConfig
	DB                  DBConfig
	SleepTime           int // Used by testing-processor only
}

// HConfig contains any HTTP options.
type HConfig struct {
	Name     string
	LogLevel string
	Address  string
}

// DBConfig contains any Database options.
type DBConfig struct {
	DataDir string
}

// GetPConfig reads in a pconfig file and returns a PConfig.
func GetPConfig(filePath string) *PConfig {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to Read Config: %v\n", err)
	}
	var P PConfig
	viper.SetDefault(`sleeptime`, 30)
	viper.SetDefault(`loglevel`, `info`)
	P.Environment = viper.GetString(`environment`)
	P.Datacenter = viper.GetString(`datacenter`)
	P.LogLevel = viper.GetString(`loglevel`)
	P.LogDir = viper.GetString(`logfile`)
	P.GroupName = viper.GetString(`groupname`)
	P.MemberName = viper.GetString(`membername`)
	P.DeleteGroupOnStart = viper.GetBool(`delete-group-on-start`)
	P.InitialOffsetOldest = viper.GetBool(`initial-offset-oldest`)
	P.KafkaBrokers = viper.GetStringSlice(`kafka-brokers`)
	P.InputTopics = viper.GetStringSlice(`input-topics`)
	P.OutputTopics = viper.GetStringSlice(`output-topics`)
	P.InOutMap = viper.GetStringMapString(`in-out-map`)

	viper.SetDefault(`http.name`, `defaultName`)
	viper.SetDefault(`http.loglevel`, `info`)
	viper.SetDefault(`http.address`, `:8181`)
	P.HTTP.Name = viper.GetString(`http.name`)
	P.HTTP.LogLevel = viper.GetString(`http.loglevel`)
	P.HTTP.Address = viper.GetString(`http.address`)

	viper.SetDefault(`database.datadir`, `./datadir`)
	P.DB.DataDir = viper.GetString(`database.datadir`)

	P.SleepTime = viper.GetInt(`sleeptime`)
	return &P
}

// VerifyPConfig verifies a PConfig.
func VerifyPConfig(P *PConfig) error {
	var err error
	switch {
	case P.GroupName == "":
		err = fmt.Errorf("missing groupname")
	case len(P.KafkaBrokers) < 1:
		err = fmt.Errorf("missing kafka brokers")
	case len(P.InOutMap) == 0 && len(P.InputTopics) == 0 && len(P.OutputTopics) == 0:
		err = fmt.Errorf("missing topics")
	}
	return err
}
