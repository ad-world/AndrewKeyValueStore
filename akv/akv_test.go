package akv

import (
	"akv/fs_ops"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func CreateTestDir() {
	// Create the test directory for this test
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fullPath := filepath.Join(cwd, "../test_store")
	err = os.Mkdir(fullPath, fs.ModePerm)
	if err != nil {
		panic(err)
	}
}

func DeleteTestDir() {
	// Delete the test directory for this test
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fullPath := filepath.Join(cwd, "../test_store")
	err = os.RemoveAll(fullPath)
	if err != nil {
		panic(err)
	}
}

func TestAKV(t *testing.T) {
	// Test setup
	CreateTestDir();
	store := &AndrewKeyValueStore{fs_operator: fs_ops.CreateFsOps("test_store")}
	
	// Test reading a key that doesn't exist
	getRequest := &GetRequest{Key: "nonexistent_key"}
	var reply string
	err := store.Get(getRequest, &reply)

	// Should receive an error here
	if err == nil {
		t.Errorf("Expected error while reading key: 'nonexistent_key', got nil")
	}
	
	// Reply should be empty string
	if reply != "" {
		t.Errorf("Expected empty string, got %v", reply)
	}
	
	// Test reading a key that does exist
	putRequest := &PutRequest{Key: "test_key", Value: "test_value"}
	var putReply bool
	err = store.Put(putRequest, &putReply)
	
	// err should be nil
	if err != nil {
		t.Errorf("Error writing key: %v", err)
	}
	
	// putReply should be true
	if(!putReply) {
		t.Errorf("Expected true, got false")
	}
	getRequest = &GetRequest{Key: "test_key"}
	err = store.Get(getRequest, &reply)

	// err should be nil
	if err != nil {
		t.Errorf("Error reading key: %v", err)
	}
	
	// reply should be "test_value"
	if reply != "test_value" {
		t.Errorf("Expected test_value, got %v", reply)
	}
	DeleteTestDir();
}