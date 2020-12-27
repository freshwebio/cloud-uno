package storage

// BucketAccessControls represents a service
// that deals with managing bucket access controls.
type BucketAccessControls interface {
	Delete()
	Get()
	Create()
	List()
	Patch()
	Update()
}
