package BlockIo

import "encoding/json"

func (blockIo *Client) GetNewAddress(args map[string]interface{}) (map[string]interface{}, error) { 
	argsObj, _ := json.Marshal(args)
	return blockIo._request("POST", "get_new_address", string(argsObj))
 }
func (blockIo *Client) GetBalance(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("GET", "get_balance", string(argsObj))
 }
func (blockIo *Client) GetMyAddresses(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_my_addresses", string(argsObj))
 }
func (blockIo *Client) GetAddressReceived(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_address_received", string(argsObj))
 }
func (blockIo *Client) GetAddressByLabel(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_address_by_label", string(argsObj))
 }
func (blockIo *Client) GetAddressBalance(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_address_balance", string(argsObj))
 }
func (blockIo *Client) CreateUser(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "create_user", string(argsObj))
 }
func (blockIo *Client) GetUsers(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_users", string(argsObj))
 }
func (blockIo *Client) GetUserBalance(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_user_balance", string(argsObj))
 }
func (blockIo *Client) GetUserAddress(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_user_address", string(argsObj))
 }
func (blockIo *Client) GetUserReceived(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_user_received", string(argsObj))
 }
func (blockIo *Client) GetTransactions(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_transactions", string(argsObj))
 }
func (blockIo *Client) SignAndFinalizeWithdrawal(args string) (map[string]interface{}, error) {
    return blockIo._request("POST", "sign_and_finalize_withdrawal", args)
 }
func (blockIo *Client) GetNewDtrustAddress(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_new_dtrust_address", string(argsObj))
 }
func (blockIo *Client) GetMyDtrustAddresses(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_my_dtrust_addresses", string(argsObj))
 }
func (blockIo *Client) GetDtrustAddressByLabel(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_dtrust_address_by_label", string(argsObj))
 }
func (blockIo *Client) GetDtrustTransactions(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_dtrust_transactions", string(argsObj))
 }
func (blockIo *Client) GetDtrustAddressBalance(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_dtrust_address_balance", string(argsObj))
 }
func (blockIo *Client) GetNetworkFeeEstimate(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_network_fee_estimate", string(argsObj))
 }
func (blockIo *Client) ArchiveAddress(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "archive_address", string(argsObj))
 }
func (blockIo *Client) UnarchiveAddress(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "unarchive_address", string(argsObj))
 }
func (blockIo *Client) GetMyArchivedAddresses(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_my_archived_addresses", string(argsObj))
 }
func (blockIo *Client) ArchiveDtrustAddress(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "archive_dtrust_address", string(argsObj))
 }
func (blockIo *Client) UnarchiveDtrustAddress(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "unarchive_dtrust_address", string(argsObj))
 }
func (blockIo *Client) GetMyArchivedDtrustAddresses(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_my_archived_dtrust_addresses", string(argsObj))
 }
func (blockIo *Client) GetDtrustNetworkFeeEstimate(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_dtrust_network_fee_estimate", string(argsObj))
 }
func (blockIo *Client) CreateNotification(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "create_notification", string(argsObj))
 }
func (blockIo *Client) DisableNotification(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "disable_notification", string(argsObj))
 }
func (blockIo *Client) EnableNotification(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "enable_notification", string(argsObj))
 }
func (blockIo *Client) GetNotifications(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_notifications", string(argsObj))
 }
func (blockIo *Client) GetRecentNotificationEvents(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_recent_notification_events", string(argsObj))
 }
func (blockIo *Client) DeleteNotification(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "delete_notification", string(argsObj))
 }
func (blockIo *Client) ValidateApiKey(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "validate_api_key", string(argsObj))
 }
func (blockIo *Client) SignTransation(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "sign_transaction", string(argsObj))
 }
func (blockIo *Client) FinalizeTransaction(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "finalize_transaction", string(argsObj))
 }
func (blockIo *Client) GetMyAddressesWithoutBalances(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_my_addresses_without_balances", string(argsObj))
 }
func (blockIo *Client) GetRawTransaction(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_raw_transaction", string(argsObj))
 }
func (blockIo *Client) GetDtrustBalance(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_dtrust_balance", string(argsObj))
 }
func (blockIo *Client) ArchiveAddresses(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "archive_addresses", string(argsObj))
 }
func (blockIo *Client) UnarchiveAddresses(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "unarchive_addresses", string(argsObj))
 }
func (blockIo *Client) ArchiveDtrustAddresses(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "archive_dtrust_addresses", string(argsObj))
 }
func (blockIo *Client) UnarchiveDtrustAddresses(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "unarchive_dtrust_addresses", string(argsObj))
 }
func (blockIo *Client) IsValidAddress(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "is_valid_address", string(argsObj))
 }
func (blockIo *Client) GetCurrentPrice(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_current_price", string(argsObj))
 }
func (blockIo *Client) GetAccountInfo(args map[string]interface{}) (map[string]interface{}, error) {
	argsObj, _ := json.Marshal(args)
    return blockIo._request("POST", "get_account_info", string(argsObj))
 }