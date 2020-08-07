package BlockIo

func (blockIo *Client) Withdraw(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._withdraw("POST", "withdraw", args); }
func (blockIo *Client) WithdrawFromUser(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._withdraw("POST", "withdraw_from_user", args); }
func (blockIo *Client) WithdrawFromLabel(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._withdraw("POST", "withdraw_from_label", args); }
func (blockIo *Client) WithdrawFromAddress(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._withdraw("POST", "withdraw_from_address", args); }
func (blockIo *Client) WithdrawFromLabels(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._withdraw("POST", "withdraw_from_labels", args); }
func (blockIo *Client) WithdrawFromAddresses(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._withdraw("POST", "withdraw_from_addresses", args); }
func (blockIo *Client) WithdrawFromUsers(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._withdraw("POST", "withdraw_from_users", args); }
func (blockIo *Client) WithdrawFromDtrustAddress(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._withdraw("POST", "withdraw_from_dtrust_address", args); }
func (blockIo *Client) WithdrawFromDtrustAddresses(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._withdraw("POST", "withdraw_from_dtrust_addresses", args); }
func (blockIo *Client) WithdrawFromDtrustLabels(args map[string]interface{}) (map[string]interface{}, error) { return blockIo._withdraw("POST", "withdraw_from_dtrust_labels", args); }