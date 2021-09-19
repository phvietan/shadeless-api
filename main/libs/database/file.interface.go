package database

type IFileDatabase interface {
	IDatabase

	Init() *FileDatabase
	GetFileByProjectAndId(project string, id string) *File
	ClearCollection()
}
