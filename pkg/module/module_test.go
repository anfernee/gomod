package module

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testModule struct {
	name         string
	versions     []string
	latest       string
	expectLatest string
}

var testModules = []testModule{
	{
		name: "github.com/osrg/gobgp/v3",
		versions: []string{
			"v3.0.0",
			"v3.1.0",
			"v3.2.0",
			"v3.3.0",
			"v3.4.0",
			"v3.5.0",
			"v3.18.0",
		},
		latest: `{
			"Version": "v3.18.0",
			"Time": "2023-09-02T12:20:21Z",
			"Origin": {
				"VCS": "git",
				"URL": "https://github.com/osrg/gobgp",
				"Ref": "refs/tags/v3.18.0",
				"Hash": "6047ca44b14e4e202cb8bf32ab4a271c89268ca1"
			}
		}`,
		expectLatest: "v3.18.0",
	},
}

func TestVersions(t *testing.T) {
	for _, testModule := range testModules {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/"+testModule.name+"/@v/list" {
				for _, version := range testModule.versions {
					fmt.Fprintln(w, version)
				}
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		m := New(testModule.name)
		m.proxy = ts.URL

		versions, err := m.Versions()
		if err != nil {
			t.Error(err)
		}

		if diff := cmp.Diff(versions, testModule.versions); diff != "" {
			t.Errorf("Versions() mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestLatest(t *testing.T) {
	for _, testModule := range testModules {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/"+testModule.name+"/@latest" {
				fmt.Fprintln(w, testModule.latest)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		m := New(testModule.name)
		m.proxy = ts.URL

		latest, err := m.Latest()
		if err != nil {
			t.Error(err)
		}

		if testModule.expectLatest != latest {
			t.Errorf("Latest() mismatch (-want +got):\n%s", cmp.Diff(testModule.expectLatest, latest))
		}
	}
}
