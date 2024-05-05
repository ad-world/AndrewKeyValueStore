package akv

type AndrewKeyValueAPI interface {
	Get(args *GetRequest, reply *string) error
	Put(args *PutRequest, value string, reply *int) error
}

type GetRequest struct {
	Key string
}

type PutRequest struct {
	Key string
	Value string
}