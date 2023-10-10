package shockwave

type ChunkMap interface {
	SetShockwave(shockwave *Shockwave) error
	AddResource(resource *ShockwaveResource) (*ShockwaveResource, error)
	GetAllResources() []*ShockwaveResource
	GetResourceById(id int32) *ShockwaveResource
	GetResourcesByTag(tag string) []*ShockwaveResource
	ToJson() (string, error)
}
