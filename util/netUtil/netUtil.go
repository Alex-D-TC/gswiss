package netUtil

type NetAddr struct {
	Address  string
	ConnType string
}

func (addr *NetAddr) Network() string {
	return addr.ConnType
}

func (addr *NetAddr) String() string {
	return addr.Address
}
