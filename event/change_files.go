package event

import "FsnotifyPlus/data"

// Структура события, которая хранит информацию об изменившихся файлов. Реализует интерфейс Event
type ChangeFileEvent struct {
	ChangedFiles data.FilesTracked //Массив файлов, которые изменились в результате события
	name         string
}

func NewChangeFileEvent(changedFiles data.FilesTracked) ChangeFileEvent {
	return ChangeFileEvent{
		ChangedFiles: changedFiles,
		name:         eventChangeFiles,
	}
}

func (e ChangeFileEvent) EventName() string {
	return e.name
}
