package BlockIo

func (blockIo *Client) Withdraw(args string) map[string]interface{} { return blockIo._withdraw("POST", "withdraw", args); }
func (blockIo *Client) WithdrawFromUser(args string) map[string]interface{} { return blockIo._withdraw("POST", "withdraw_from_user", args); }
func (blockIo *Client) WithdrawFromLabel(args string) map[string]interface{} { return blockIo._withdraw("POST", "withdraw_from_label", args); }
func (blockIo *Client) WithdrawFromAddress(args string) map[string]interface{} { return blockIo._withdraw("POST", "withdraw_from_address", args); }
func (blockIo *Client) WithdrawFromLabels(args string) map[string]interface{} { return blockIo._withdraw("POST", "withdraw_from_labels", args); }
func (blockIo *Client) WithdrawFromAddresses(args string) map[string]interface{} { return blockIo._withdraw("POST", "withdraw_from_addresses", args); }
func (blockIo *Client) WithdrawFromUsers(args string) map[string]interface{} { return blockIo._withdraw("POST", "withdraw_from_users", args); }
func (blockIo *Client) WithdrawFromDtrustAddress(args string) map[string]interface{} { return blockIo._withdraw("POST", "withdraw_from_dtrust_address", args); }
func (blockIo *Client) WithdrawFromDtrustAddresses(args string) map[string]interface{} { return blockIo._withdraw("POST", "withdraw_from_dtrust_addresses", args); }
func (blockIo *Client) WithdrawFromDtrustLabels(args string) map[string]interface{} { return blockIo._withdraw("POST", "withdraw_from_dtrust_labels", args); }