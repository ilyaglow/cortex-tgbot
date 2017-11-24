package cortexbot

import (
	"fmt"
	"testing"

	gocortex "github.com/ilyaglow/go-cortex"
)

func TestConstructJobFromIP(t *testing.T) {
	ip := "1.1.1.1"
	tlp := 1

	j, err := constructJob(ip, tlp)
	if err != nil {
		t.Error("Failed to construct Cortex Job from an IP")
	}

	if j.Attributes.DataType != "ip" {
		t.Error("Failed datatype in Cortex Job constructed from an IP")
	}
}

func TestConstructJobFromLink(t *testing.T) {
	link := "https://subdomain.domain.com/route/query?param=val"
	tlp := 0

	j, err := constructJob(link, tlp)
	if err != nil {
		t.Error("Failed to construct Cortex Job from a URL")
	}

	if j.Attributes.DataType != "url" {
		t.Error("Failed datatype in Cortex Job constructed from a URL")
	}
}

func TestConstructJobFromDomain(t *testing.T) {
	domain := "subdomain.domain.com"
	tlp := 2

	j, err := constructJob(domain, tlp)
	if err != nil {
		t.Error("Failed to construct Cortex Job from a domain")
	}

	if j.Attributes.DataType != "domain" {
		t.Error("Failed datatype in Cortex Job constructed from a domain")
	}
}

func TestConstructJobFromHash(t *testing.T) {
	hash := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	tlp := 3

	j, err := constructJob(hash, tlp)
	if err != nil {
		t.Error("Failed to construct Cortex Job from a hash")
	}

	if j.Attributes.DataType != "hash" {
		t.Error("Failed datatype in Cortex Job constructed from a hash")
	}
}

func TestConstructJobFromUnknown(t *testing.T) {
	unknown := "unknown_data"
	tlp := 2

	_, err := constructJob(unknown, tlp)
	if err == nil {
		t.Error("Unknown data didn't trigger the error")
	}
}

func TestBuildTaxonomies(t *testing.T) {
	var txs []gocortex.Taxonomy

	tx1 := gocortex.Taxonomy{
		Predicate: "Predicate1",
		Namespace: "Namespace1",
		Value:     "Value1",
		Level:     "safe",
	}
	txs = append(txs, tx1)

	tx2 := gocortex.Taxonomy{
		Predicate: "Predicate2",
		Namespace: "Namespace2",
		Value:     "Value2",
		Level:     "info",
	}
	txs = append(txs, tx2)

	expected := fmt.Sprintf("%s:%s = %s, %s:%s = %s",
		tx1.Namespace, tx1.Predicate, tx1.Value,
		tx2.Namespace, tx2.Predicate, tx2.Value,
	)
	if buildTaxonomies(txs) != expected {
		t.Errorf("Wrong taxonomies format:\n%s\n%s", expected, buildTaxonomies(txs))
	}
}
