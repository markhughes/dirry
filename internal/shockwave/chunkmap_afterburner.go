package shockwave

import (
	"encoding/json"
	"fmt"

	"github.com/markhughes/dirry/internal/utils"
)

type AfterburnerChunkMap struct {
	Shockwave *Shockwave
	Resources []*ShockwaveResource
}

func (chunkMap *AfterburnerChunkMap) SetShockwave(shockwave *Shockwave) error {
	chunkMap.Shockwave = shockwave
	return nil
}

func (chunkMap *AfterburnerChunkMap) AddResource(resource *ShockwaveResource) (*ShockwaveResource, error) {
	if chunkMap.Resources == nil {
		chunkMap.Resources = make([]*ShockwaveResource, 0)
	}

	chunkMap.Resources = append(chunkMap.Resources, resource)

	return resource, nil
}

func (chunkMap *AfterburnerChunkMap) GetAllResources() []*ShockwaveResource {
	return chunkMap.Resources
}

func (chunkMap *AfterburnerChunkMap) GetResourceById(id int32) *ShockwaveResource {
	for i := range chunkMap.Resources {
		if chunkMap.Resources[i].ResourceId == id {
			utils.DebugMsg("shockwave", "Found resource: %v\n", chunkMap.Resources[i].ChunkType)
			return chunkMap.Resources[i]
		}
	}

	utils.DebugMsg("shockwave", "Did not find resource with id: %v\n", id)

	return nil
}

func (chunkMap *AfterburnerChunkMap) GetResourcesByTag(tag string) []*ShockwaveResource {
	// fmt.Printf("GetResourcesByTag: %s\n", tag)

	var resources []*ShockwaveResource = make([]*ShockwaveResource, 0)
	for i := range chunkMap.Resources {
		if chunkMap.Resources[i].ChunkType == tag {
			utils.DebugMsg("shockwave", "Found resource: %v\n", chunkMap.Resources[i].ChunkType)
			resources = append(resources, chunkMap.Resources[i])
		}
	}

	if len(resources) == 0 {
		utils.DebugMsg("shockwave", "Did not find resource with tag: %v\n", tag)
	}

	return resources
}

func (chunkMap *AfterburnerChunkMap) ToJson() (string, error) {
	bytes, err := json.MarshalIndent(chunkMap.Resources, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal resources to json: %s", err)
	}
	return string(bytes), nil
}
