package akv

import (
	"akv/fs_ops"
)

type AndrewKeyValueStore struct {
	store map[string]string
	fs_operator fs_ops.FsOpsInterface
}

func CreateAndrewKeyValueStore() *AndrewKeyValueStore {
	return &AndrewKeyValueStore{
		store: make(map[string]string),
		fs_operator: fs_ops.CreateFsOps("store"),
	}
}

func (store *AndrewKeyValueStore) Get(args *GetRequest, reply *string) error {
	data, err := store.fs_operator.ReadKey(args.Key)
	if err != nil {
		*reply = ""
		return err
	}

	*reply = string(data)

	return nil
}

func (store *AndrewKeyValueStore) Put(args *PutRequest, reply *bool) error {
	err := store.fs_operator.WriteKey(args.Key, []byte(args.Value), 0644);

	if err != nil {
		*reply = false
		return err
	} else {
		*reply = true
	}

	return nil
}

func (store *AndrewKeyValueStore) Delete(args *DeleteRequest, reply *bool) error {
	err := store.fs_operator.DeleteKey(args.Key)

	if err != nil {
		*reply = false
		return err
	} else {
		*reply = true
	}

	return nil
}