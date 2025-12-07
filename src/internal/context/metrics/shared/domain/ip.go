package domain

import (
	"net/netip"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type Ip netip.Addr

func NewIp(ip string) (Ip, error) {
	netipAddrPort, err := netip.ParseAddrPort(ip)
	if err == nil {
		return Ip(netipAddrPort.Addr()), nil
	}

	netipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return Ip{}, shared_domain_context.NewValidationError("Ip", "is not a valid IP address")
	}

	return Ip(netipAddr), nil
}

func (ip Ip) String() string {
	return netip.Addr(ip).String()
}
