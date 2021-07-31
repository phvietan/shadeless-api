package database

type IFileDatabase interface {
	Init() *FileDatabase
	CreateFile(file *File) error
	GetFileByProjectAndId(project string, id string) *File
	UpdateProjectName(oldName string, newName string) error
	DeleteFilesByProjectName(projectName string) error
	ClearCollection()
}
