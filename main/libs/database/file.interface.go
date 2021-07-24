package database

type IFileDatabase interface {
	Init() *FileDatabase
	CreateFile(file *File) error
	GetFileByProjectAndId(project string, id string) *File
	DeleteFilesByProjectName(projectName string) error
	ClearCollection()
}
