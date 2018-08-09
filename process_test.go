package cortexbot

import (
	"fmt"
	"testing"

	"github.com/ilyaglow/go-cortex"
)

func TestObservableConstructor(t *testing.T) {
	var cases = []struct {
		data     string
		tlp      int
		dataType string
	}{
		{"1.1.1.1", 0, "ip"},
		{"https://subdomain.domain.com/route/query?param=val", 1, "url"},
		{"subdomain.domain.com", 2, "domain"},
		{"email.address@domain.com", 3, "mail"},
		{"email+address@domain.com", 3, "mail"},
		{"a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", 0, "hash"},
		{"randomdata", 1, "other"},
	}

	for _, c := range cases {
		a := newArtifact(c.data, c.tlp)
		if a.Description() != c.data {
			t.Errorf("need %s, got %s", c.data, a.Description())
		}

		if a.Type() != c.dataType {
			t.Errorf("need %s, got %s", c.dataType, a.Type())
		}
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

	expected := fmt.Sprintf("%s:%s = %s\n%s:%s = %s",
		tx1.Namespace, tx1.Predicate, tx1.Value,
		tx2.Namespace, tx2.Predicate, tx2.Value,
	)
	if buildTaxonomies(txs) != expected {
		t.Errorf("Wrong taxonomies format:\n%s\n%s", expected, buildTaxonomies(txs))
	}
}
