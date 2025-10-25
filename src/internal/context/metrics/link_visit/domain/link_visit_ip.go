package domain

import "net/netip"

type LinkVisitIP netip.AddrPort

func NewLinkVisitIP(ip string) (LinkVisitIP, error) {
	netipAddrPort, err := netip.ParseAddrPort(ip)
	if err != nil {
		return LinkVisitIP{}, err
	}
	return LinkVisitIP(netipAddrPort), nil
}
