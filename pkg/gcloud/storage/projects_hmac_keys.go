package storage

type ProjectsHMACKeys interface {
	Create()
	Delete()
	Get()
	List()
	Update()
}
