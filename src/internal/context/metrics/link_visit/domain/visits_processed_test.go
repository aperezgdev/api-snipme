package domain

import (
	"testing"
)

func TestNewLinkVisitsProcessedDomainEvent(t *testing.T) {
	t.Run("Creates event correctly", func(t *testing.T) {
		t.Parallel()
		linkId := "00000000-0000-0000-0000-000000000001"
		totalViews := uint(100)
		uniqueViews := uint(75)
		ipsFrequency := []IpsVisits{
			{Ip: "192.168.1.1", TotalViews: 50, UniqueViews: 40},
			{Ip: "192.168.1.2", TotalViews: 50, UniqueViews: 35},
		}

		event := NewLinkVisitsProcessedDomainEvent(linkId, totalViews, uniqueViews, ipsFrequency)

		if event.AggregateID() != linkId {
			t.Errorf("Expected AggregateID %s, got %s", linkId, event.AggregateID())
		}
		if event.TotalViews != totalViews {
			t.Errorf("Expected TotalViews %d, got %d", totalViews, event.TotalViews)
		}
		if event.UniqueViews != uniqueViews {
			t.Errorf("Expected UniqueViews %d, got %d", uniqueViews, event.UniqueViews)
		}
		if len(event.IpsFrequency) != len(ipsFrequency) {
			t.Errorf("Expected %d IpsFrequency, got %d", len(ipsFrequency), len(event.IpsFrequency))
		}
		if event.Name() != LinkVisitsProcessedEventName {
			t.Errorf("Expected Name %s, got %s", LinkVisitsProcessedEventName, event.Name())
		}
		if event.OccurredOn().IsZero() {
			t.Error("Expected OccurredOn to be set")
		}
	})

	t.Run("Creates event with empty IpsFrequency", func(t *testing.T) {
		t.Parallel()
		linkId := "00000000-0000-0000-0000-000000000001"
		totalViews := uint(0)
		uniqueViews := uint(0)
		ipsFrequency := []IpsVisits{}

		event := NewLinkVisitsProcessedDomainEvent(linkId, totalViews, uniqueViews, ipsFrequency)

		if event.TotalViews != 0 {
			t.Errorf("Expected TotalViews 0, got %d", event.TotalViews)
		}
		if event.UniqueViews != 0 {
			t.Errorf("Expected UniqueViews 0, got %d", event.UniqueViews)
		}
		if len(event.IpsFrequency) != 0 {
			t.Errorf("Expected empty IpsFrequency, got %d items", len(event.IpsFrequency))
		}
	})

	t.Run("Creates event with multiple IPs", func(t *testing.T) {
		t.Parallel()
		linkId := "00000000-0000-0000-0000-000000000001"
		totalViews := uint(300)
		uniqueViews := uint(150)
		ipsFrequency := []IpsVisits{
			{Ip: "192.168.1.1", TotalViews: 100, UniqueViews: 50},
			{Ip: "192.168.1.2", TotalViews: 100, UniqueViews: 50},
			{Ip: "192.168.1.3", TotalViews: 100, UniqueViews: 50},
		}

		event := NewLinkVisitsProcessedDomainEvent(linkId, totalViews, uniqueViews, ipsFrequency)

		if len(event.IpsFrequency) != 3 {
			t.Errorf("Expected 3 IpsFrequency, got %d", len(event.IpsFrequency))
		}
		for i, ipVisit := range event.IpsFrequency {
			if ipVisit.Ip != ipsFrequency[i].Ip {
				t.Errorf("Expected IP %s at index %d, got %s", ipsFrequency[i].Ip, i, ipVisit.Ip)
			}
			if ipVisit.TotalViews != ipsFrequency[i].TotalViews {
				t.Errorf("Expected TotalViews %d at index %d, got %d", ipsFrequency[i].TotalViews, i, ipVisit.TotalViews)
			}
		}
	})
}

func TestLinkVisitsProcessed_Name(t *testing.T) {
	t.Run("Returns correct event name", func(t *testing.T) {
		t.Parallel()
		event := LinkVisitsProcessed{}
		if event.Name() != LinkVisitsProcessedEventName {
			t.Errorf("Expected Name %s, got %s", LinkVisitsProcessedEventName, event.Name())
		}
		if event.Name() != "LinkVisitsProcessed" {
			t.Errorf("Expected Name 'LinkVisitsProcessed', got %s", event.Name())
		}
	})
}
