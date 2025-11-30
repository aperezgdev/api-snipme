package domain

import "net/netip"

type LinkVisitIP netip.Addr

func NewLinkVisitIP(ip string) (LinkVisitIP, error) {
	netipAddrPort, err := netip.ParseAddrPort(ip)
	if err == nil {
		return LinkVisitIP(netipAddrPort.Addr()), nil
	}
	
	netipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return LinkVisitIP{}, err
	}

	return LinkVisitIP(netipAddr), nil
}

func (ip LinkVisitIP) String() string {
	return netip.Addr(ip).String()
}
