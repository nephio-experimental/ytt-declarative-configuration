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
