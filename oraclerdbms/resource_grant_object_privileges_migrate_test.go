package oraclerdbms

import (
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestGrantObjectPrivilegeMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		ID           string
		Attributes   map[string]string
		Expected     map[string]string
		Meta         interface{}
	}{
		"v0_1": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"privs_sha256":   "xxxxxxxxxxxx",
				"objects_sha256": "yyyyyyyyyy",
			},
			Expected: map[string]string{"privs_sha256": "xxxxxxxxxxxx",
				"objects_sha256": "yyyyyyyyyy"},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.Attributes,
		}
		is, err := resourceOracleRdbmsGrantObjectPrivilegeMigrate(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}
		t.Logf("attrubutes: %v\n", is)
	}
}
