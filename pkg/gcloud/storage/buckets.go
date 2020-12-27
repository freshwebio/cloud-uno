package storage

type Buckets interface {
	Delete()
	Get()
	GetIAMPolicy()
	Create()
	List()
	ListChannels()
	LockRetentionPolicy()
	Patch()
	SetIAMPolicy()
	TestIamPermissions()
	Update()
}
