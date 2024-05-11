package akv

type AndrewKeyValueAPI interface {
	Get(args *GetRequest, reply *string) error
	Put(args *PutRequest, reply *bool) error
	Delete(args *DeleteRequest, reply *bool) error
}

type DeleteRequest struct {
	Key string
}

type GetRequest struct {
	Key string
}

type PutRequest struct {
	Key string
	Value string
}