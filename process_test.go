package cortexbot

import (
	"fmt"
	"testing"

	"github.com/ilyaglow/go-cortex"
)

func TestObservableConstructor(t *testing.T) {
	var cases = []struct {
		data     string
		tlp      cortex.TLP
		pap      cortex.PAP
		dataType string
	}{
		{"1.1.1.1", cortex.TLPWhite, cortex.PAPWhite, "ip"},
		{"https://subdomain.domain.com/route/query?param=val", cortex.TLPGreen, cortex.PAPGreen, "url"},
		{"subdomain.domain.com", cortex.TLPAmber, cortex.PAPAmber, "domain"},
		{"email.address@domain.com", cortex.TLPRed, cortex.PAPRed, "mail"},
		{"email+address@domain.com", cortex.TLPRed, cortex.PAPRed, "mail"},
		{"a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", cortex.TLPWhite, cortex.PAPWhite, "hash"},
		{"randomdata", cortex.TLPGreen, cortex.PAPGreen, "other"},
	}

	for _, c := range cases {
		a := newArtifact(c.data, c.tlp, c.pap)
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
