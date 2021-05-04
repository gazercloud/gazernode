package repeater_bin_client

type BinFrameTask struct {
	SessionId string
	Client    *RepeaterBinClient
	Frame     *BinFrame
}
