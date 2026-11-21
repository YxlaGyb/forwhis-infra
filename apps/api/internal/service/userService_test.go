package service

import "testing"

func TestRegistrationRole(t *testing.T) {
	cases := []struct {
		existing int64
		want     string
	}{
		{0, "superadmin"},
		{1, "user"},
		{42, "user"},
	}
	for _, tc := range cases {
		if got := registrationRole(tc.existing); got != tc.want {
			t.Errorf("registrationRole(%d) = %q, want %q", tc.existing, got, tc.want)
		}
	}
}
