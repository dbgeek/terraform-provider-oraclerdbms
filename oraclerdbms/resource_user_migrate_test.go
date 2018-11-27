package oraclerdbms

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestUserMigrateState(t *testing.T) {
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
				"password": "xxxxxxxxxxxx",
			},
			Expected: map[string]string{},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.Attributes,
		}
		is, err := resourceOracleRdbmsUserMigrate(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}
		if _, ok := tc.Attributes["password"]; ok {
			t.Errorf("Attributes contains password\n")
		}
		t.Logf("attrubutes: %v\n", is)
	}
}
