package powershell

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//
// the goal around context cache is to get around golang's problem of not being able to
// marshal to cgo a golang pointer to a golang pointer.
//
var contextCache sync.Map
var contextLookupKey uint64

func storeRunspaceContext(context runspaceContext) uint64 {
	contextLookup := atomic.AddUint64(&contextLookupKey, 1)
	contextCache.Store(contextLookup, context)
	return contextLookup
}
func getRunspaceContext(key uint64) runspaceContext {

	contextInterface, ok := contextCache.Load(key)
	if !ok {
		panic(fmt.Sprint("failed to load context key:", key))
	}
	return contextInterface.(runspaceContext)
}
func deleteRunspaceContextLookup(key uint64) {
	contextCache.Close(key)
}
