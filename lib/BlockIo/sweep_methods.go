package BlockIo

func (blockIo *Client) SweepFromAddress(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._sweep("POST", "sweep_from_address", args) }