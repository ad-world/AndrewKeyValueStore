package state

import (
	"akv/akv"
	"encoding/json"
	"os"
)

func SaveState (state *State, akv *akv.AndrewKeyValueStore) error {
	store := akv.Store
	jsonString, err := json.Marshal(store)

	if err != nil {
		return err
	}	

	f, err := os.Create("state.json")
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(jsonString)

	return nil
}
