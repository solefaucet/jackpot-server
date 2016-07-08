package core

// GetDestAddress gets destination address
func (w Wallet) GetDestAddress() (string, error) {
	return w.client.GetAccountAddress("")
}
