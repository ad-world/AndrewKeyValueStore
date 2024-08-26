package state

import (
	"akv/src/akv"
	"os"
	"testing"
)

func CreateTestDir(dir string) {
	os.Mkdir(dir, 0777)	
}

func DeleteTestDir(dir string) {
	os.RemoveAll(dir)
}

func TestSaveState(t *testing.T) {
	var test_dir string = "test_temp"
	// Setup
	CreateTestDir(test_dir)

	// Create AKV
	store := akv.CreateAndrewKeyValueStore()
	state := &State{ dir: test_dir}

	// Put a key in the store
	putRequest := &akv.PutRequest{Key: "test_key", Value: "test_value"}
	putReply, _ := store.Put(putRequest)

	if putReply != true {
		t.Errorf("Error writing key: %v", putReply)
	}

	// Put another key in the store
	putRequest = &akv.PutRequest{Key: "test_key2", Value: "test_value2"}
	putReply, _ = store.Put(putRequest)
	if putReply != true {
		t.Errorf("Error writing key: %v", putReply)
	}

	// Save the state
	err := SaveState(state, store)
	if err != nil {
		t.Errorf("Error saving state: %v", err)
	}	

	// Cleanup
	DeleteTestDir(test_dir)
}

func TestRestoreState(t *testing.T) {
	var test_dir string = "test_temp"
	// Setup
	CreateTestDir(test_dir)
	
	// Create AKV
	store := akv.CreateAndrewKeyValueStore()
	state := &State{ dir: test_dir}

	// Put a key in the store
	putRequest := &akv.PutRequest{Key: "test_key", Value: "test_value"}
	putReply, _ := store.Put(putRequest)

	if putReply != true {
		t.Errorf("Error writing key: %v", putReply)
	}

	// Put another key in the store
	putRequest = &akv.PutRequest{Key: "test_key2", Value: "test_value2"}
	putReply, _ = store.Put(putRequest)
	if putReply != true{
		t.Errorf("Error writing key: %v", putReply)
	}

	// Save the state
	err := SaveState(state, store)
	if err != nil {
		t.Errorf("Error saving state: %v", err)
	}
	
	new_store := akv.CreateAndrewKeyValueStore()
	new_state := &State{ dir: test_dir}
	// Restore the state
	err = RestoreState(new_state, new_store)
	if err != nil {
		t.Errorf("Error restoring state: %v", err)
	}

	// Check that the keys are the same
	getRequest := &akv.GetRequest{Key: "test_key"}
	getReply, err := new_store.Get(getRequest)

	if err != nil {
		t.Errorf("Error getting key: %v", err)
	}

	if getReply.Value != "test_value" {
		t.Errorf("Expected value: test_value, got: %v", getReply)
	}

	getRequest = &akv.GetRequest{Key: "test_key2"}
	getReply, err = new_store.Get(getRequest)
	
	if err != nil {
		t.Errorf("Error getting key: %v", err)
	}

	if getReply.Value != "test_value2" {
		t.Errorf("Expected value: test_value2, got: %v", getReply)
	}
	
	// Cleanup
	DeleteTestDir(test_dir)
}