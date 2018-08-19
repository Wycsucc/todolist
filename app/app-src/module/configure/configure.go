package configure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// TODOList configure struct
type TODOList struct {
	Mode          string `json:"mode"`
	Host          string `json:"host"`
	URL           string `json:"url"`
	Port          int    `json:"port"`
	TimeOutReadS  int    `json:"time_out_read_s"`
	TimeOutWriteS int    `json:"time_out_write_s"`
}

// Log save log configure information
type Log struct {
	File   string `json:"file"`
	Access string `json:"access"`
}

//CollectionName mongodb collection names configure
type CollectionName struct {
	User string `json:"user"`
}

// Mongo mongodb configure struct
type Mongo struct {
	Hosts           string         `json:"hosts"`
	ConnectTimeOutS string         `json:"connect_time_out_s"`
	Username        string         `json:"username"`
	Password        string         `json:"password"`
	DatabaseName    string         `json:"database_name"`
	CollectionNames CollectionName `json:"collection_names"`
}

//configure struct
type configure struct {
	TODOList TODOList `json:"todolist"`
	Log      Log      `json:"log"`
	Mongo    Mongo    `json:"mongo"`
}

var (
	conf     *configure
	confOnce sync.Once
)

// Configure load json configure file
func Configure(ctx context.Context, file string) *configure {
	confOnce.Do(func() {
		conf = &configure{}
		if err := conf.init(file); err != nil {
			log.Fatalln(err)
		}
		conf.readEnv()
	})
	return conf
}

// init load json configure file
func (c *configure) init(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error open file %s fail, %v", file, err)
	}
	defer fd.Close()

	decoder := json.NewDecoder(fd)
	for {
		if err := decoder.Decode(c); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (c *configure) readEnv() {
	mongoHost := os.Getenv("MONGO_DB_HOST")
	if mongoHost != "" {
		c.Mongo.Hosts = mongoHost
	}
}

func (c *configure) String() string {
	js, _ := json.MarshalIndent(c, "", "\t")
	return fmt.Sprintf("%s", js)
}

// GetConfigure get the configure object
func GetConfigure() *configure {
	return conf
}
