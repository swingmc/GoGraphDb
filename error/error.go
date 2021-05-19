package error

import "errors"

var (
	VertexDataVersionTooLate = errors.New("vertex data version too late, rollback and retry please.")
	EdgeDataVersionTooLate = errors.New("edge data version too late, rollback and retry please.")
)
