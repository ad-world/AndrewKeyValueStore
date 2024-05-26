package akv

import "net/rpc"

type AndrewKeyValueClient struct {
	*rpc.Client
}

// Get calls the Get method of the AndrewKeyValueStore on the server and returns the value of the key if it exists, otherwise it will return an error.
func (client *AndrewKeyValueClient) Get(key string) (string, error) {
	var reply string
	err := client.Call("AndrewKeyValueStore.Get", &GetRequest{Key: key}, &reply)
	return reply, err
}

// Put calls the Put method of the AndrewKeyValueStore on the server and returns true if the operation is successful, otherwise it will return an error.
func (client *AndrewKeyValueClient) Put(key string, value string) (bool, error) {
	var reply bool
	err := client.Call("AndrewKeyValueStore.Put", &PutRequest{Key: key, Value: value}, &reply)
	return reply, err
}

// Delete calls the Delete method of the AndrewKeyValueStore on the server and returns true if the operation is successful, otherwise it will return an error.
func (client *AndrewKeyValueClient) Delete(key string) (bool, error) {
	var reply bool
	err := client.Call("AndrewKeyValueStore.Delete", &DeleteRequest{Key: key}, &reply)
	return reply, err
}

func (client *AndrewKeyValueClient) GetLastUpdated(key string) (string, error) {
	var reply string
	err := client.Call("AndrewKeyValueStore.GetLastUpdated", &GetLastUpdatedRequest{Key: key}, &reply)
	return reply, err
}
