package mocks

type MockIPAddress struct {
	Name    string
	Address string
}

func (addr MockIPAddress) Network() string {
	return addr.Name
}

func (addr MockIPAddress) String() string {
	return addr.Address
}

func NewMockIPAddress(name, address string) MockIPAddress {
	return MockIPAddress{
		Name:    name,
		Address: address,
	}
}
