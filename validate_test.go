package cortexbot

import (
	"testing"
)

func TestIsHash(t *testing.T) {
	hashes := []string{
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		"098f6bcd4621d373cade4e832627b4f6",
		"a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
	}

	nonvalid := []string{
		"qwerty.com",
		"1.1.1.1",
		"http://long-link-to-test-regex.domain-name.com",
	}

	for _, v := range hashes {
		if !IsHash(v) {
			t.Errorf("Valid hash %s is not a hash", v)
		}
	}

	for _, n := range nonvalid {
		if IsHash(n) {
			t.Errorf("Non valid data %s is a hash", n)
		}
	}
}

func TestIsDNSName(t *testing.T) {
	domains := []string{
		"google.com",
		"long-subdomain0.thedomain.shop",
	}

	nonvalid := []string{
		"non-valid-too-long-domain-name------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------.com",
		"http://the-link.com",
		"localhostdomain-that-we-do-not-allow",
	}

	for _, v := range domains {
		if !IsDNSName(v) {
			t.Errorf("Valid domain %s is not a domain", v)
		}
	}

	for _, n := range nonvalid {
		if IsDNSName(n) {
			t.Errorf("Non valid data %s is a domain", n)
		}
	}
}
