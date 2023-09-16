package main

import "testing"

func TestLatestPatch(t *testing.T) {
	testCases := []struct {
		name     string
		versions []string
		current  string
		want     string
	}{
		{
			name: "latest found",
			versions: []string{
				"v3.0.0",
				"v3.1.0",
				"v3.2.0",
				"v3.3.0",
				"v3.4.0",
				"v3.5.0",
				"v3.5.1",
				"v3.5.2",
				"v3.18.0",
			},
			current: "v3.5.1",
			want:    "v3.5.2",
		},
		{
			name: "latest is the current",
			versions: []string{
				"v3.0.0",
				"v3.1.0",
				"v3.2.0",
				"v3.3.0",
				"v3.4.0",
				"v3.5.0",
				"v3.5.1",
				"v3.5.2",
				"v3.18.0",
			},
			current: "v3.4.0",
			want:    "v3.4.0",
		},
		{
			name: "current not found",
			versions: []string{
				"v3.0.0",
				"v3.1.0",
				"v3.2.0",
				"v3.4.0",
				"v3.5.0",
				"v3.5.1",
				"v3.5.2",
				"v3.18.0",
			},
			current: "v3.3.0",
			want:    "v3.3.0",
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got := latestPatch(test.current, test.versions)
			if got != test.want {
				t.Errorf("got %s, want %s", got, test.want)
			}
		})
	}
}
