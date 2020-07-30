package BlockIo

func (blockIo *Client) SweepFromAddress(args string) map[string]interface{} { return blockIo._sweep("POST", "sweep_from_address", args) }