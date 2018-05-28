package cortexbot

import (
	"fmt"
	"testing"

	"gopkg.ilya.app/ilyaglow/go-cortex.v2"
)

func TestConstructJobFromIP(t *testing.T) {
	ip := "1.1.1.1"
	tlp := 1

	j, err := newArtifact(ip, tlp)
	if err != nil {
		t.Error("Failed to construct Cortex Job from an IP")
	}

	if j.(*cortex.Task).DataType != "ip" {
		t.Error("Failed datatype in Cortex Job constructed from an IP")
	}
}

func TestConstructJobFromLink(t *testing.T) {
	link := "https://subdomain.domain.com/route/query?param=val"
	tlp := 0

	j, err := newArtifact(link, tlp)
	if err != nil {
		t.Error("Failed to construct Cortex Job from a URL")
	}

	if j.(*cortex.Task).DataType != "url" {
		t.Error("Failed datatype in Cortex Job constructed from a URL")
	}
}

func TestConstructJobFromDomain(t *testing.T) {
	domain := "subdomain.domain.com"
	tlp := 2

	j, err := newArtifact(domain, tlp)
	if err != nil {
		t.Error("Failed to construct Cortex Job from a domain")
	}

	if j.(*cortex.Task).DataType != "domain" {
		t.Error("Failed datatype in Cortex Job constructed from a domain")
	}
}

func TestConstructJobFromHash(t *testing.T) {
	hash := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	tlp := 3

	j, err := newArtifact(hash, tlp)
	if err != nil {
		t.Error("Failed to construct Cortex Job from a hash")
	}

	if j.(*cortex.Task).DataType != "hash" {
		t.Error("Failed datatype in Cortex Job constructed from a hash")
	}
}

func TestConstructJobFromUnknown(t *testing.T) {
	unknown := "unknown_data"
	tlp := 2

	_, err := newArtifact(unknown, tlp)
	if err == nil {
		t.Error("Unknown data didn't trigger the error")
	}
}

func TestBuildTaxonomies(t *testing.T) {
	var txs []cortex.Taxonomy

	tx1 := cortex.Taxonomy{
		Predicate: "Predicate1",
		Namespace: "Namespace1",
		Value:     "Value1",
		Level:     "safe",
	}
	txs = append(txs, tx1)

	tx2 := cortex.Taxonomy{
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
