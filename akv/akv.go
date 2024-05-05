package akv

import (
	"fmt"
	"akv/fs_ops"
)

type AndrewKeyValueStore struct {
	store map[string]string
	fs_operator fs_ops.FsOpsInterface
}

func CreateAndrewKeyValueStore() *AndrewKeyValueStore {
	return &AndrewKeyValueStore{
		store: make(map[string]string),
		fs_operator: fs_ops.CreateFsOps(),
	}
}

func (store *AndrewKeyValueStore) Get(args *GetRequest, reply *string) error {
	value, ok := store.store[args.Key]
	if !ok {
		return fmt.Errorf("key %s not found", args.Key)
	}
	*reply = value
	return nil
}

func (store *AndrewKeyValueStore) Put(args *PutRequest, value string, reply *int) error {
	store.store[args.Key] = value
	*reply = len(value)
	return nil
}