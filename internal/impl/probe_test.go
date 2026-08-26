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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestPickProbe(t *testing.T) {
	// Two distinguishable probes, shaped like the ones this override exists
	// for: an app-wide TCP check every group can pass, and an httpGet check
	// that only the group hosting the listener could ever satisfy.
	appProbe := &corev1.Probe{ProbeHandler: corev1.ProbeHandler{
		TCPSocket: &corev1.TCPSocketAction{Port: intstr.FromInt32(10000)},
	}}
	groupProbe := &corev1.Probe{ProbeHandler: corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{Path: "/health", Port: intstr.FromInt32(20001)},
	}}

	for _, c := range []struct {
		name  string
		group *corev1.Probe
		app   *corev1.Probe
		want  *corev1.Probe
	}{
		{"group overrides app", groupProbe, appProbe, groupProbe},
		{"nil group inherits app", nil, appProbe, appProbe},
		{"group applies with no app probe", groupProbe, nil, groupProbe},
		{"both nil yields nil", nil, nil, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := pickProbe(c.group, c.app); got != c.want {
				t.Errorf("pickProbe() = %v, want %v", got, c.want)
			}
		})
	}
}
