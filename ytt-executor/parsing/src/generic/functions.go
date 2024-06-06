/*
Copyright 2022. Ericsson AB

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generic

import (
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"sigs.k8s.io/kustomize/kyaml/resid"
)

func ParseResourceMap(Items fn.KubeObjects) map[resid.ResId]*fn.KubeObject {
	resourceMap := map[resid.ResId]*fn.KubeObject{}

	for _, item := range Items {
		resId := NewResIdWithNameAndKind(item.GetName(), item.GetKind())
		resourceMap[resId] = item
	}

	return resourceMap
}

func NewResIdWithNameAndKind(name string, kind string) resid.ResId {
	return resid.ResId{
		Name: name,
		Gvk:  resid.Gvk{Kind: kind},
	}
}
