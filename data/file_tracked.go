package data

type fileEvent int

const (
	REMOVE fileEvent = iota
	ADD
	CHANGE
)

// Структура, которая хранит список отслеживаемых файлов
type FilesTracked []*FileTracked

type FileTracked struct {
	FileName  string    //Название изменившегося файла
	FileEvent fileEvent //Какое изменения претерпел файл
}

func NewFileTracked(filename string, event fileEvent) FileTracked {
	return FileTracked{
		FileName:  filename,
		FileEvent: event,
	}
}

func NewFilesTracked(size int) FilesTracked {
	return make(FilesTracked, 0, size)
}
