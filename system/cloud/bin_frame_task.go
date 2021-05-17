package cloud

type BinFrameTask struct {
	SessionId string
	Client    *Connection
	Frame     *BinFrame
}
