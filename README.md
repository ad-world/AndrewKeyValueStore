# Andrew Key Value Store
- This README is outdated, as the program now uses websockets instead of RPC. I will update the README soon <3
- A distribued key value store that uses semantics based on the Andrew File System (developed at Carnegie Mellon University).

## Contents
1. [Server Process](#server-process)
2. [Client Process](#client-process)

### Server Process
The server process is where the actual key-value pairs are stored. It is the single source of truth.
The server process exposes two RPC methods.
- `Get(String key)`
- `Put(String key, String value)`

These two RPCs work as expected. `Get` will try and retrieve the value at the key, and return an error code if it cannot. `Put` will create a new entry with the `(key, value)` tuple.

The server keeps track of all clients that have a specific (key, value) pair stored in their cache. It does this so that it can notify those clients when a specific value has changed, so that each of those clients can invalidate their clients.

### Client Process
The client process is a process that can make `Get` and `Put` RPC requests to the server process. The client exposes one RPC method (that is only called by the server).
- `CacheInvalidation(String key)`


On a client `Get` call, first the cache is checked. If a value for the key is found in the cache, the value is returned and the operation is completed. If the value is not completed, the client makes an RPC call to the server.

On a client `Put` call, the client issues a `Put` RPC call to the server. On success, it populates the cache with the key-value pair.

In either case, the tuple is kept in the cache until one of the following happens
- The client crashes, at which point the cache is reset.
- The server sends the client a `CacheInvalidation` request, at which point the specific key is flushed from the cache.


This is how the client maintains consistency of caches.

