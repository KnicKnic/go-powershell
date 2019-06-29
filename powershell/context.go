package powershell

import (
	"sync"
	"sync/atomic"
)

//
// the goal around context cache is to get around golang's problem of not being able to
// marshal to cgo a golang pointer to a golang pointer.
//
var contextCache sync.Map
var contextLookupKey uint64 = 0

func StoreRunspaceContext(context Context) uint64 {
	contextLookup := atomic.AddUint64(&contextLookupKey, 1)
	contextCache.Store(contextLookup, context)
	return contextLookup
}
func GetRunspaceContext(key uint64) (Context, bool) {

	contextInterface, ok := contextCache.Load(key)
	if ok {
		return contextInterface.(Context), true
	} else {
		return Context{}, false
	}
}
func DeleteRunspaceContextLookup(key uint64) {
	contextCache.Delete(key)
}
