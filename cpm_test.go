package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func openMockDatabase() (*CpmDatabase, error) {
	var db CpmDatabase
	var err error
	db.Database, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("sql.Open() failed: %s", err)
	}

	return &db, nil
}

func closeMockDatabase(db *CpmDatabase) error {
	err := db.Database.Close()
	if err != nil {
		return fmt.Errorf("db.Database.Close() failed: %s", err)
	}

	return nil
}

func TestInsert(t *testing.T) {
	OldOpenDatabase := OpenDatabase
	OpenDatabase = openMockDatabase
	defer func() { OpenDatabase = OldOpenDatabase }()
	OldCloseDatabase := CloseDatabase
	CloseDatabase = closeMockDatabase
	defer func() { CloseDatabase = OldCloseDatabase }()

	os.Args = []string{"", "create", "-m", "mymachine", "-s", "myservice", "-u", "myuser", "-p", "mypassword"}
	buf := new(bytes.Buffer)

	ret := Main(buf)

	wantedRet := 0
	if ret != wantedRet {
		t.Errorf("Main() return = %q, want %q", ret, wantedRet)
	}

	// TODO also assert that the password is indeed stored
}
