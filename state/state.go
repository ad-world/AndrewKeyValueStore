package state

import (
	"akv/akv"
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveState (state *State, akv *akv.AndrewKeyValueStore) error {
	store := akv.Store
	// Create a JSON string from the store
	jsonString, err := json.Marshal(store)

	if err != nil {
		return err
	}	

	// Create a file for the string
	f, err := os.Create(filepath.Join(state.dir, "state.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the JSON string to the file
	f.Write(jsonString)

	return nil
}

func RestoreState (state *State, akv *akv.AndrewKeyValueStore) error {
	// Open the state file
	f, err := os.Open(filepath.Join(state.dir, "state.json"))
	if err != nil {
		return err
	}

	store := akv.Store
	// Decode the JSON string into the store
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&store)

	if err != nil {
		return err
	}

	return nil
}