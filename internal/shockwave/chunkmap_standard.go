package shockwave

import (
	"encoding/json"
	"fmt"
)

type StandardChunkMap struct {
	Shockwave *Shockwave
	Resources []*ShockwaveResource
}

func (chunkMap *StandardChunkMap) SetShockwave(shockwave *Shockwave) error {
	if shockwave == nil {
		return fmt.Errorf("shockwave cannot be nil")
	}
	chunkMap.Shockwave = shockwave
	return nil
}

func (chunkMap *StandardChunkMap) AddResource(resource *ShockwaveResource) (*ShockwaveResource, error) {
	chunkMap.Resources = append(chunkMap.Resources, resource)

	return resource, nil
}

func (chunkMap *StandardChunkMap) GetAllResources() []*ShockwaveResource {
	return chunkMap.Resources
}

func (chunkMap *StandardChunkMap) GetResourceById(id int32) *ShockwaveResource {
	for i := range chunkMap.Resources {
		var resource = chunkMap.Resources[i]
		if resource.ResourceId == id {
			return resource
		}
	}
	return nil
}

func (chunkMap *StandardChunkMap) GetResourcesByTag(tag string) []*ShockwaveResource {
	var resources []*ShockwaveResource
	for i := range chunkMap.Resources {
		var resource = chunkMap.Resources[i]
		if resource.ChunkType == tag {
			resources = append(resources, resource)
		}
	}
	return resources
}

func (chunkMap *StandardChunkMap) ToJson() (string, error) {
	bytes, err := json.MarshalIndent(chunkMap.Resources, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal resources to json: %s", err)
	}
	return string(bytes), nil
}
