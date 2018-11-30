package nsutil

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
)

// If we are selecting resources by their name and the namespace in which they
// are installed, displayName and displayNamespace are identical to the
// resourceRef. However, meshes are selected by the ns in which they were
// installed, so we need both representations
// NOTE: if we add select helper utils for other resources we should make a
// general "select by resource ref" util
type ResSelect struct {
	displayName      string
	displayNamespace string
	resourceRef      core.ResourceRef
}

type ResMap map[string]ResSelect

func generateMeshSelectOptions(nsrMap options.NsResourceMap) ([]string, ResMap) {

	var meshOptions []string
	// map the key to the mesh select object
	// key is namespace, name
	meshMap := make(ResMap)

	for installNs, nsr := range nsrMap {
		for _, meshRef := range nsr.MeshesByInstallNs {
			selectMenuString := fmt.Sprintf("%v, %v", installNs, meshRef.Name)
			meshOptions = append(meshOptions, selectMenuString)
			meshMap[selectMenuString] = ResSelect{
				displayName:      meshRef.Name,
				displayNamespace: installNs,
				resourceRef: core.ResourceRef{
					Name:      meshRef.Name,
					Namespace: meshRef.Namespace,
				},
			}
		}
	}
	return meshOptions, meshMap
}

func generateCommonResourceSelectOptions(typeName string, nsrMap options.NsResourceMap) ([]string, ResMap) {

	var resOptions []string
	// map the key to the res select object
	// key is namespace, name
	resMap := make(ResMap)

	for namespace, nsr := range nsrMap {
		var resArray []string
		switch typeName {
		case "secret":
			resArray = nsr.IstioSecrets
		case "upstream":
			resArray = nsr.Upstreams
		case "awssecret":
			// TODO: the secret mappings here use the proto name in the resource map, and translate to a user facing name based on the use case. cleanup?
			resArray = nsr.GlooSecrets
		default:
			panic(fmt.Errorf("resource type %v not recognized", typeName))
		}
		for _, res := range resArray {
			selectMenuString := fmt.Sprintf("%v, %v", namespace, res)
			resOptions = append(resOptions, selectMenuString)
			resMap[selectMenuString] = ResSelect{
				displayName:      res,
				displayNamespace: namespace,
				resourceRef: core.ResourceRef{
					Name:      res,
					Namespace: namespace,
				},
			}
		}
	}
	return resOptions, resMap
}
