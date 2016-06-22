package main

import "testing"

func TestCidrList_List(t *testing.T) {
	var listTests = []struct {
		cidrs    cidrList
		expected []string
	}{
		{cidrList{"192.168.1.0/24", "1234::1222:0/16"}, []string{"192.168.1.0/24", "1234::1222:0/16"}},
	}

	for row, test := range listTests {
		list := test.cidrs.List()

		for i, cidr := range list {
			if cidr != test.expected[i] {
				t.Errorf("Row: %d returned unexpected CIDR, got: %s, wanted: %s", row, cidr, test.expected[i])
			}
		}
	}
}

func TestCidrList_String(t *testing.T) {
	var stringTests = []struct {
		cidrs    cidrList
		expected string
	}{
		{cidrList{"192.168.1.0/24", "1234::1222:0/16"}, "192.168.1.0/24 1234::1222:0/16"},
	}

	for row, test := range stringTests {
		str := test.cidrs.String()

		if str != test.expected {
			t.Errorf("Row: %d returned unexpected string, got: %s, wanted: %s", row, str, test.expected)
		}
	}
}

func TestCidrList_Set(t *testing.T) {
	var setTests = []struct {
		cidrs    []string
		expected cidrList
	}{
		{[]string{"192.168.1.0/24", "1234::1222:0/16"}, cidrList{"192.168.1.0/24", "1234::1222:0/16"}},
	}

	var cidrs cidrList
	for row, test := range setTests {
		for _, cidr := range test.cidrs {
			if err := cidrs.Set(cidr); err != nil {
				t.Errorf("Row: %d returned unexpected error: %v", row, err) // This should be impossible, but what the heck.
			}
		}

		for i, cidr := range test.expected {
			if cidrs[i] != cidr {
				t.Errorf("Row: %d returned unexpected CIDR, got: %s, wanted: %s", row, cidrs[i], cidr)
			}
		}
	}
}
