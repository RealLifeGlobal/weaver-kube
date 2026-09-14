// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package impl

import (
	"testing"
	"time"
)

func TestParseBuildTimeout(t *testing.T) {
	for _, c := range []struct {
		name    string
		in      string
		want    time.Duration
		wantErr bool
	}{
		{"empty uses the default", "", 15 * time.Minute, false},
		{"explicit duration", "10m", 10 * time.Minute, false},
		{"compound duration", "1h30m", 90 * time.Minute, false},
		{"zero is rejected", "0s", 0, true},
		{"negative is rejected", "-1m", 0, true},
		{"bare number is rejected", "120", 0, true},
		{"garbage is rejected", "soon", 0, true},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseBuildTimeout(c.in)
			if (err != nil) != c.wantErr {
				t.Fatalf("parseBuildTimeout(%q) error = %v, wantErr %v", c.in, err, c.wantErr)
			}
			if got != c.want {
				t.Errorf("parseBuildTimeout(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestDefaultBuilderImage(t *testing.T) {
	for _, c := range []struct {
		name string
		in   string
		want string
	}{
		{"patch release", "go1.26.5", "golang:1.26"},
		{"first release of a minor", "go1.27", "golang:1.27"},
		{"release candidate", "go1.27rc1", "golang:1.27"},
		{"devel toolchain falls back", "devel go1.28-8c2e1a4 Mon Sep 1 2026", fallbackBuilderImage},
		{"empty falls back", "", fallbackBuilderImage},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := defaultBuilderImage(c.in); got != c.want {
				t.Errorf("defaultBuilderImage(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
