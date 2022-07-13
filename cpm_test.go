package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func createTestDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("sql.Open() failed: %s", err)
	}

	return db, nil
}

func openTestDatabase(sqlDb *sql.DB) func() (*CpmDatabase, error) {
	return func() (*CpmDatabase, error) {
		var db CpmDatabase
		db.Database = sqlDb
		return &db, nil
	}
}

func closeTestDatabase(db *CpmDatabase) error {
	return nil
}

func TestInsert(t *testing.T) {
	db, err := createTestDatabase()
	defer db.Close()
	if err != nil {
		t.Errorf("createTestDatabase() err = %q, want nil", err)
	}
	OldOpenDatabase := OpenDatabase
	OpenDatabase = openTestDatabase(db)
	defer func() { OpenDatabase = OldOpenDatabase }()
	OldCloseDatabase := CloseDatabase
	CloseDatabase = closeTestDatabase
	defer func() { CloseDatabase = OldCloseDatabase }()
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	expectedType := "plain"
	os.Args = []string{"", "create", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Errorf("Main() = %q, want %q", actualRet, expectedRet)
	}
	rows, err := db.Query("select machine, service, user, password, type from passwords")
	if err != nil {
		t.Errorf("db.Query() err = %q, want nil", err)
	}
	var actualMachine string
	var actualService string
	var actualUser string
	var actualPassword string
	var actualType string
	expectedNext := true
	actualNext := rows.Next()
	if actualNext != expectedNext {
		t.Errorf("rows.Next() = %v, want %v", actualNext, expectedNext)
	}
	err = rows.Scan(&actualMachine, &actualService, &actualUser, &actualPassword, &actualType)
	if err != nil {
		t.Errorf("rows.Scan() = %q, want nil", err)
	}
	if actualMachine != expectedMachine {
		t.Errorf("actualMachine = %q, want %q", actualMachine, expectedMachine)
	}
	if actualService != expectedService {
		t.Errorf("actualService = %q, want %q", actualService, expectedService)
	}
	if actualUser != expectedUser {
		t.Errorf("actualUser = %q, want %q", actualUser, expectedUser)
	}
	if actualPassword != expectedPassword {
		t.Errorf("actualPassword = %q, want %q", actualPassword, expectedPassword)
	}
	if actualType != expectedType {
		t.Errorf("actualType = %q, want %q", actualType, expectedType)
	}
	expectedNext = false
	actualNext = rows.Next()
	if actualNext != expectedNext {
		t.Errorf("rows.Next() = %v, want %v", actualNext, expectedNext)
	}
}

func TestSelect(t *testing.T) {
	db, err := createTestDatabase()
	defer db.Close()
	if err != nil {
		t.Errorf("createTestDatabase() err = %q, want nil", err)
	}
	OldOpenDatabase := OpenDatabase
	OpenDatabase = openTestDatabase(db)
	defer func() { OpenDatabase = OldOpenDatabase }()
	OldCloseDatabase := CloseDatabase
	CloseDatabase = closeTestDatabase
	defer func() { CloseDatabase = OldCloseDatabase }()
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	expectedType := "plain"
	err = initDatabase(db)
	if err != nil {
		t.Errorf("initDatabase() = %q, want nil", err)
	}
	err = createPassword(db, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
	if err != nil {
		t.Errorf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Errorf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine, service: myservice, user: myuser, password type: plain, password: mypassword\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Errorf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
}

func TestUpdate(t *testing.T) {
	db, err := createTestDatabase()
	defer db.Close()
	if err != nil {
		t.Errorf("createTestDatabase() err = %q, want nil", err)
	}
	OldOpenDatabase := OpenDatabase
	OpenDatabase = openTestDatabase(db)
	defer func() { OpenDatabase = OldOpenDatabase }()
	OldCloseDatabase := CloseDatabase
	CloseDatabase = closeTestDatabase
	defer func() { CloseDatabase = OldCloseDatabase }()
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "newpassword"
	expectedType := "plain"
	err = initDatabase(db)
	if err != nil {
		t.Errorf("initDatabase() = %q, want nil", err)
	}
	err = createPassword(db, expectedMachine, expectedService, expectedUser, "oldpassword", expectedType)
	if err != nil {
		t.Errorf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "update", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Errorf("Main() = %q, want %q", actualRet, expectedRet)
	}
	rows, err := db.Query("select machine, service, user, password, type from passwords")
	if err != nil {
		t.Errorf("db.Query() err = %q, want nil", err)
	}
	var actualMachine string
	var actualService string
	var actualUser string
	var actualPassword string
	var actualType string
	expectedNext := true
	actualNext := rows.Next()
	if actualNext != expectedNext {
		t.Errorf("rows.Next() = %v, want %v", actualNext, expectedNext)
	}
	err = rows.Scan(&actualMachine, &actualService, &actualUser, &actualPassword, &actualType)
	if err != nil {
		t.Errorf("rows.Scan() = %q, want nil", err)
	}
	if actualMachine != expectedMachine {
		t.Errorf("actualMachine = %q, want %q", actualMachine, expectedMachine)
	}
	if actualService != expectedService {
		t.Errorf("actualService = %q, want %q", actualService, expectedService)
	}
	if actualUser != expectedUser {
		t.Errorf("actualUser = %q, want %q", actualUser, expectedUser)
	}
	if actualPassword != expectedPassword {
		t.Errorf("actualPassword = %q, want %q", actualPassword, expectedPassword)
	}
	if actualType != expectedType {
		t.Errorf("actualType = %q, want %q", actualType, expectedType)
	}
	expectedNext = false
	actualNext = rows.Next()
	if actualNext != expectedNext {
		t.Errorf("rows.Next() = %v, want %v", actualNext, expectedNext)
	}
}
