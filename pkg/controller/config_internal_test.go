package controller

import (
	"testing"
)

func Test_parseNameWithVersion(t *testing.T) {
	t.Parallel()
	data := []struct {
		title      string
		name       string
		expName    string
		expVersion string
	}{
		{
			title:      "no version",
			name:       "foo",
			expName:    "foo",
			expVersion: "",
		},
		{
			title:      "with version",
			name:       "foo@v1.0.0",
			expName:    "foo",
			expVersion: "v1.0.0",
		},
		{
			title:      "invalid name @v1.0.0",
			name:       "@v1.0.0",
			expName:    "",
			expVersion: "v1.0.0",
		},
		{
			title:      "invalid name foo@",
			name:       "foo@",
			expName:    "foo",
			expVersion: "",
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			name, version := parseNameWithVersion(d.name)
			if name != d.expName {
				t.Fatalf("name is got %s, wanted %s", name, d.expName)
			}
			if version != d.expVersion {
				t.Fatalf("version is got %s, wanted %s", version, d.expVersion)
			}
		})
	}
}

func TestCondition_match(t *testing.T) {
	t.Parallel()
	data := []struct {
		title     string
		os        string
		arch      string
		condition *Condition
		expMatch  bool
	}{
		{
			title: "match os and arch when using On",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				On: []ConditionSpec{
					{
						GOOS:   "foo",
						GOARCH: "amd64",
					},
					{
						GOOS:   "bar",
						GOARCH: "amd64",
					},
				},
			},
			expMatch: true,
		},
		{
			title: "unmatch os and arch when using On",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				On: []ConditionSpec{
					{
						GOOS:   "bar",
						GOARCH: "amd64",
					},
					{
						GOOS:   "baz",
						GOARCH: "amd64",
					},
				},
			},
			expMatch: false,
		},
		{
			title: "match os when using On",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				On: []ConditionSpec{
					{
						GOOS:   "foo",
					},
					{
						GOOS:   "bar",
					},
				},
			},
			expMatch: true,
		},
		{
			title: "unmatch os when using On",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				On: []ConditionSpec{
					{
						GOOS:   "bar",
					},
					{
						GOOS:   "baz",
					},
				},
			},
			expMatch: false,
		},
		{
			title: "match arch when using On",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				On: []ConditionSpec{
					{
						GOARCH: "amd64",
					},
					{
						GOARCH: "arm64",
					},
				},
			},
			expMatch: true,
		},
		{
			title: "unmatch arch when using On",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				On: []ConditionSpec{
					{
						GOARCH: "arm64",
					},
					{
						GOARCH: "arm",
					},
				},
			},
			expMatch: false,
		},
		{
			title: "match os and arch when using Ignoring",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				Ignoring: []ConditionSpec{
					{
						GOOS:   "foo",
						GOARCH: "amd64",
					},
					{
						GOOS:   "bar",
						GOARCH: "amd64",
					},
				},
			},
			expMatch: false,
		},
		{
			title: "unmatch os and arch when using Ignoring",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				Ignoring: []ConditionSpec{
					{
						GOOS:   "bar",
						GOARCH: "amd64",
					},
					{
						GOOS:   "baz",
						GOARCH: "amd64",
					},
				},
			},
			expMatch: true,
		},
		{
			title: "match os when using Ignoring",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				Ignoring: []ConditionSpec{
					{
						GOOS:   "foo",
					},
					{
						GOOS:   "bar",
					},
				},
			},
			expMatch: false,
		},
		{
			title: "unmatch os when using Ignoring",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				Ignoring: []ConditionSpec{
					{
						GOOS:   "bar",
					},
					{
						GOOS:   "baz",
					},
				},
			},
			expMatch: true,
		},
		{
			title: "match arch when using Ignoring",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				Ignoring: []ConditionSpec{
					{
						GOARCH: "amd64",
					},
					{
						GOARCH: "arm64",
					},
				},
			},
			expMatch: false,
		},
		{
			title: "unmatch os and arch when using Ignoring",
			os:    "foo",
			arch:  "amd64",
			condition: &Condition{
				Ignoring: []ConditionSpec{
					{
						GOARCH: "arm64",
					},
					{
						GOARCH: "arm",
					},
				},
			},
			expMatch: true,
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			match := d.condition.match(d.os, d.arch)
			if match != d.expMatch {
				t.Fatalf("wanted %v, go %v", d.expMatch, match)
			}
		})
	}
}
