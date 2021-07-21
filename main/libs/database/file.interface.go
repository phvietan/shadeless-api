package database

type IFileDatabase interface {
	Init() *FileDatabase
	CreateFile(file *File) error
	DeleteFilesByProjectName(projectName string) error
}
