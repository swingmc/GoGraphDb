package transaction

import (
	"sync"
)

var (
	TransactionCounter = make(chan int64, 10000)
	lock = sync.RWMutex{}
)


var TransactionGetter map[int64]*Transaction

func init() {
	close(StopTheWorld)
	TransactionGetter = map[int64]*Transaction{}
}

func GetTransaction(versionId int64) *Transaction{
	lock.RLock()
	t, ok := TransactionGetter[versionId]
	if !ok {
		lock.RUnlock()
		return nil
	}
	lock.RUnlock()
	return t
}

func addTransaction(versionId int64, transaction *Transaction) {
	lock.Lock()
	TransactionGetter[versionId] = transaction
	lock.Unlock()
}

func CleanTransaction() {
	lock.Lock()
	TransactionGetter = map[int64]*Transaction{}
	lock.Unlock()
}

func LockTransactionMap() {
	lock.Lock()
}

func UnlockTransactionMap(){
	lock.Unlock()
}

func RLockTransactionMap() {
	lock.RLock()
}

func RUnlockTransactionMap(){
	lock.RUnlock()
}