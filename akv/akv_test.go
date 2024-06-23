package akv

import (
	"testing"
)


func TestAKV(t *testing.T) {
	// Test setup
	store := CreateAndrewKeyValueStore();
	
	// Test reading a key that doesn't exist
	getRequest := &GetRequest{Key: "nonexistent_key"}
	var reply Value
	err := store.Get(getRequest, &reply)

	// Should receive an error here
	if err == nil {
		t.Errorf("Expected error while reading key: 'nonexistent_key', got nil")
	}
	
	// Reply should be empty string
	if reply.Value != "" {
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
	if reply.Value != "test_value" {
		t.Errorf("Expected test_value, got %v", reply)
	}
}