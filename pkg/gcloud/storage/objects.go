package storage

type Objects interface {
	Compose()
	Copy()
	Delete()
	Get()
	Create()
	List()
	Patch()
	Rewrite()
	Update()
	WatchAll()
}
