package storage

type ObjectAccessControls interface {
	Delete()
	Get()
	Create()
	List()
	Patch()
	Update()
}
