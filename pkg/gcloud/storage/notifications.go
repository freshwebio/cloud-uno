package storage

type Notifications interface {
	Delete()
	Get()
	Create()
	List()
	Patch()
	Update()
}
