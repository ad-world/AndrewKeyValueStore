package state

import (
	"akv/akv"
	"testing"
)

func TestSaveState(t *testing.T) {
	// Create AKV
	store := akv.CreateAndrewKeyValueStore()
	state := &State{}

	// Put a key in the store
	putRequest := &akv.PutRequest{Key: "test_key", Value: "test_value"}
	reply := false
	putReply := store.Put(putRequest, &reply)

	if putReply != nil {
		t.Errorf("Error writing key: %v", putReply)
	}

	// Put another key in the store
	putRequest = &akv.PutRequest{Key: "test_key2", Value: "test_value2"}
	putReply = store.Put(putRequest, &reply)
	if putReply != nil {
		t.Errorf("Error writing key: %v", putReply)
	}

	// Save the state
	err := SaveState(state, store)
	if err != nil {
		t.Errorf("Error saving state: %v", err)
	}	
}