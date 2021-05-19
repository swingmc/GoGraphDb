package transaction

var TransactionGetter map[int64]*Transaction

func init() {
	close(StopTheWorld)
	TransactionGetter = map[int64]*Transaction{}
}

func GetTransaction(versionId int64) *Transaction{
	t, ok := TransactionGetter[versionId]
	if !ok {
		return nil
	}
	return t
}

func CleanTransaction() {
	TransactionGetter = map[int64]*Transaction{}
}