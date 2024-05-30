package mysql

import "testing"

func TestDumper(t *testing.T) {

	username := "root"
	password := "root"
	host := "localhost"
	port := "3306"
	inputDir := "/Users/apple/migrate/gep"
	err := restoreFromDirectory(username, password, host, port, inputDir)
	if err != nil {
		t.Log(err.Error())
	}
}
