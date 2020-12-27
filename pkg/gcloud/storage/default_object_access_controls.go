package storage

type DefaultObjectAccessControls interface {
	DeleteDefaultObjectAccessControl()
	GetDefaultObjectAccessControl()
	CreateDefaultObjectAccessControl()
	ListDefaultObjectAccessControls()
}
