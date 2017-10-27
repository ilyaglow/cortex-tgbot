package cortexbot

import (
	"testing"
)

func TestConstructJobFromIP(t *testing.T) {
	ip := "1.1.1.1"

	j, err := constructJob(ip)
	if err != nil {
		t.Error("Failed to construct Cortex Job from an IP")
	}

	if j.Attributes.DataType != "ip" {
		t.Error("Failed datatype in Cortex Job constructed from an IP")
	}
}

func TestConstructJobFromLink(t *testing.T) {
	link := "https://subdomain.domain.com/route/query?param=val"

	j, err := constructJob(link)
	if err != nil {
		t.Error("Failed to construct Cortex Job from a URL")
	}

	if j.Attributes.DataType != "url" {
		t.Error("Failed datatype in Cortex Job constructed from a URL")
	}
}

func TestConstructJobFromDomain(t *testing.T) {
	domain := "subdomain.domain.com"

	j, err := constructJob(domain)
	if err != nil {
		t.Error("Failed to construct Cortex Job from a domain")
	}

	if j.Attributes.DataType != "domain" {
		t.Error("Failed datatype in Cortex Job constructed from a domain")
	}
}

func TestConstructJobFromHash(t *testing.T) {
	hash := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"

	j, err := constructJob(hash)
	if err != nil {
		t.Error("Failed to construct Cortex Job from a hash")
	}

	if j.Attributes.DataType != "hash" {
		t.Error("Failed datatype in Cortex Job constructed from a hash")
	}
}

func TestConstructJobFromUnknown(t *testing.T) {
	unknown := "unknown_data"

	_, err := constructJob(unknown)
	if err == nil {
		t.Error("Unknown data didn't trigger the error")
	}
}
