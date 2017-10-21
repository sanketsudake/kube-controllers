// Copyright (c) 2017 Tigera, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package converter

import (
	"fmt"

	api "github.com/projectcalico/libcalico-go/lib/apis/v2"
	"github.com/projectcalico/libcalico-go/lib/backend/k8s/conversion"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

// ProfileNameFormat Format used by policy controller to name Calico profiles
const ProfileNameFormat = "k8s_ns."

type namespaceConverter struct {
}

// NewNamespaceConverter Constructor for namespaceConverter
func NewNamespaceConverter() Converter {
	return &namespaceConverter{}
}
func (nc *namespaceConverter) Convert(k8sObj interface{}) (interface{}, error) {
	var c conversion.Converter
	namespace, ok := k8sObj.(*v1.Namespace)
	if !ok {
		tombstone, ok := k8sObj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return nil, fmt.Errorf("couldn't get object from tombstone %+v", k8sObj)
		}
		namespace, ok = tombstone.Obj.(*v1.Namespace)
		if !ok {
			return nil, fmt.Errorf("tombstone contained object that is not a Namespace %+v", k8sObj)
		}
	}
	kvp, err := c.NamespaceToProfile(namespace)
	if err != nil {
		return nil, err
	}
	profile := kvp.Value.(*api.Profile)

	return *profile, nil
}

// GetKey returns name of the Profile as its key.  For Profiles
// backed by Kubernetes namespaces and managed by this controller, the name
// is of format `k8s_ns.name`.
func (nc *namespaceConverter) GetKey(obj interface{}) string {
	profile := obj.(api.Profile)
	return profile.Name
}

func (p *namespaceConverter) DeleteArgsFromKey(key string) (string, string) {
	// Not namespaced, so just return the key, which is the profile name.
	return "", key
}