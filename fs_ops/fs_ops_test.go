package fs_ops

import (
	"io/fs"
	"os"
	"testing"
	"path/filepath"
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

func TestFsOps_ReadKey(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current working directory: %v", err)
	}

	// Setup the test
	CreateTestDir()
	fullPath := filepath.Join(cwd, "../test_store/test")

	temp, err := os.Create(fullPath)

	if err != nil {
		t.Errorf("Error creating temp file: %v", err)
	}

	testData := []byte("test data")
	_, err = temp.Write(testData)

	if err != nil {
		t.Errorf("Error writing to temp file: %v", err)
	}

	// Actual Test
	fsOps := CreateFsOps("test_store")
	data, err := fsOps.ReadKey("test")

	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	if string(data) != string(testData) {
		t.Errorf("Expected %s, got %s", testData, data)
	}

	// Cleanup
	DeleteTestDir()
}

func TestFsOps_WriteKey(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current working directory: %v", err)
	}

	// Setup the test
	CreateTestDir()
	fullPath := filepath.Join(cwd, "../test_store/test_file")

	temp, err := os.Create(fullPath)

	if err != nil {
		t.Errorf("Error creating temp file: %v", err)
	}

	defer os.Remove(temp.Name())

	// Actual test
	fsOps := CreateFsOps("test_store")
	testData := []byte("test data")
	
	err = fsOps.WriteKey("test_file", testData, fs.ModePerm)
	if err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	data, err := fsOps.ReadKey("test_file")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	if string(data) != string(testData) {
		t.Errorf("Expected %s, got %s", testData, data)
	}

	// Cleanup
	DeleteTestDir()
}