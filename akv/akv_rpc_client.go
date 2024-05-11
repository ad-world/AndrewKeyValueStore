package akv

import "net/rpc"

type AndrewKeyValueClient struct {
	*rpc.Client
}

func (client *AndrewKeyValueClient) Get(key string) (string, error) {
	var reply string
	err := client.Call("AndrewKeyValueStore.Get", &GetRequest{Key: key}, &reply)
	return reply, err
}

func (client *AndrewKeyValueClient) Put(key string, value string) (bool, error) {
	var reply bool
	err := client.Call("AndrewKeyValueStore.Put", &PutRequest{Key: key, Value: value}, &reply)
	return reply, err
}

func (client *AndrewKeyValueClient) Delete(key string) (bool, error) {
	var reply bool
	err := client.Call("AndrewKeyValueStore.Delete", &DeleteRequest{Key: key}, &reply)
	return reply, err
}
