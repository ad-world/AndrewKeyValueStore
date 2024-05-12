package state

import (
	"akv/akv"
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveState (state *State, akv *akv.AndrewKeyValueStore) error {
	store := akv.Store
	jsonString, err := json.Marshal(store)

	if err != nil {
		return err
	}	

	f, err := os.Create(filepath.Join(state.dir, "state.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(jsonString)

	return nil
}
