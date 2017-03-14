package xorm

import (
	"errors"
	"flag"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	testEngine *Engine
	dbType     string
	connStr    string
)

func prepareSqlite3Engine() error {
	os.Remove("./test.db")
	var err error
	testEngine, err = NewEngine("sqlite3", "./test.db")
	if err != nil {
		return err
	}
	testEngine.ShowSQL(*showSQL)
	return nil
}

func prepareMysqlEngine() error {
	if testEngine == nil {
		var err error
		testEngine, err = NewEngine("mysql", connStr)
		if err != nil {
			return err
		}
		testEngine.ShowSQL(*showSQL)
		_, err = testEngine.Exec("DROP DATABASE")
		if err != nil {
			return err
		}
	}
	testEngine.ShowSQL(*showSQL)
	return nil
}

func prepareEngine() error {
	if dbType == "sqlite" {
		return prepareSqlite3Engine()
	} else if dbType == "mysql" {
		return prepareMysqlEngine()
	}
	return errors.New("Unknown test database driver")
}

var (
	db      = flag.String("db", "sqlite", "the tested database")
	showSQL = flag.Bool("show_sql", true, "show generated SQLs")
)

func TestMain(m *testing.M) {
	flag.Parse()

	if db != nil {
		dbType = *db
	}

	if err := prepareEngine(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	assert.NoError(t, prepareEngine())

	if err := testEngine.Ping(); err != nil {
		t.Fatal(err)
	}
}
