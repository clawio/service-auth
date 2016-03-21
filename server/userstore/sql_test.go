package userstore

import (
	"database/sql"
	"io/ioutil"
	"path"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestFindByCredentials(t *testing.T) {
	// test for invalid driver
	_, err := NewSQLUserStore("fakedriver", "")
	if err == nil {
		t.Fatal("newSQLUserStore must fail when driver is invalid")
	}

	// test for not being able to open db
	_, err = NewSQLUserStore("sqlite3", "/this/dir/does/not/exist")
	if err == nil {
		t.Fatal("newSQLUserStore must fail when db cannot be open")
	}

	// TODO(labkode) test for automigrate

	tmpDir, err := ioutil.TempDir("", "sql_userstore_test")
	if err != nil {
		t.Fatal(err)
	}
	dbPath := path.Join(tmpDir, "sqlite3_userstore.db")
	store, err := NewSQLUserStore("sqlite3", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if store == nil {
		t.Fatal("store cannot be nil")
	}

	// test for unexistent identity
	_, err = store.FindByCredentials("fake", "fake")
	if err == nil {
		t.Fatal("findByCredentials must fail when identity is not found")
	}

	// test valid identity
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `insert into user_entities values ("test", "test@test.com", "Test", "testpwd")`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		t.Fatal(err)
	}
	identity, err := store.FindByCredentials("test", "testpwd")
	if err != nil {
		t.Fatal(err)
	}
	if identity == nil {
		t.Fatal("identity cannot be nil")
	}
	if identity.Username != "test" ||
		identity.Email != "test@test.com" ||
		identity.DisplayName != "Test" {
		t.Fatal("identity fields are different from the ones in the db")
	}
}
