package db

import (
	"strings"
	"time"

	"github.com/globalsign/mgo"
)

// Session struct
type Session struct {
	S          *mgo.Session
	Database   *mgo.Database
	Collection *mgo.Collection
}

// MongoClient 定义Mongo的操作实例
type MongoClient struct {
	Hosts                string
	Database             string
	Collection           string
	ConnectTimeOutSecond int
	Session
}

func (m *MongoClient) Connect() error {
	dialInfo := &mgo.DialInfo{
		Addrs:    strings.Split(m.Hosts, ","),
		Timeout:  time.Duration(m.ConnectTimeOutSecond) * time.Second,
		Database: m.Database,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}
	m.S = session
	// Optional. Switch the session to a monotonic behavior.
	m.S.SetMode(mgo.Monotonic, true)
	return nil
}

// InitDatabase func
func (m *MongoClient) InitDatabase() {
	db := m.Session.S.DB(m.Database)
	m.Session.Database = db
}

// InitCollection func
func (m *MongoClient) InitCollection() {
	c := m.Session.Database.C(m.Collection)
	m.Session.Collection = c
}

// Close func
func (m *MongoClient) Close() {
	if m.S != nil {
		m.S.LogoutAll()
		m.S.Close()
	}
}

// GetCollection func
func (m *MongoClient) GetCollection() *mgo.Collection {
	return m.Session.Collection
}
