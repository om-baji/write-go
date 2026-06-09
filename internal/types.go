package internal

type LedgerEntry struct {
	Magic     int32  `json:"magic"`
	Seq       int64  `json:"seq"`
	Crc       uint64 `json:"crc"`
	Timestamp string `json:"timestamp"`
	Body      string `json:"body"`
}

type Segment struct {
	Id   int    `json:"id"`
	Size int    `json:"size"`
	Path string `json:"path"`
}
