package network

type Interface struct {
	uuid       string
	Name       string
	MacAddress string
	IPAddress  string
	Netmask    string
	Gateway    string
}

func NewInterface(name, mac, ip, netmask, gateway string) *Interface {
	return &Interface{
		Name:       name,
		MacAddress: mac,
		IPAddress:  ip,
		Netmask:    netmask,
		Gateway:    gateway,
	}

}
func (ni *Interface) GetUuid() string   { return ni.uuid }
func (ni *Interface) SetUuid(id string) { ni.uuid = id }
