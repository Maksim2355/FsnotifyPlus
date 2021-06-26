package file_watcher

import (
	"FsnotifyPlus/data"
	"FsnotifyPlus/event"
	"FsnotifyPlus/hub"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"time"
)

//Стандартный интервал между отправки широковещательных рассылок
const standardDuration = 2 * time.Second

//Стандартное наполнение буфера, при достижении максимального размера, будет эмититься событие
const standardBuffFilling = 3

//Сущность для отслеживание изменения файлов в директории и отправки событий в хаб
type Watcher struct {
	Watcher          *fsnotify.Watcher //Стандартный watcher из fsnotify
	WatchDir         string            //Корневая директория для отслеживания файла
	IsRecursiveWatch bool              //Стоит ли отслеживать все вложенные директории
	hub              *hub.Hub          //Шина, куда эмитим события изменения файла
	buff             data.FilesTracked //Буфер событий. Если буфер заполняется полностью, то происходит мгновенная отправка, иначе сообщения рассылаются через определенный интервал
	sendInterval     time.Duration     //Интервал между отправкой сообщений
}

func NewWatcher(dir string, isRecursiveWatch bool, hub *hub.Hub) *Watcher {
	return &Watcher{
		WatchDir:         dir,
		IsRecursiveWatch: isRecursiveWatch,
		hub:              hub,
		buff:             data.NewFilesTracked(standardBuffFilling),
		sendInterval:     standardDuration,
	}
}

//Метод запускает отслеживание файлов в указанной директории @WatchDir
func (w *Watcher) RunWatch() {
	w.initWatcher()
	for {
		tick := time.NewTicker(w.sendInterval)
		select {
		case ev := <-w.Watcher.Events:
			w.handleChangeFileEvent(ev)
		case <-tick.C:
			w.sendingChangeFileEvent()
		}
	}
}

//Добавляет измененный файл в буфер, если есть дубликат, то удаляет. Если мы добавили файл и буфер полный, то отправляем эвент
func (w *Watcher) addFileChangingInBuffAndSentEventIfNeeded(file data.FileTracked) {
	isAddedElementDuplicate := false
	for _, v := range w.buff {
		if v.FileName == file.FileName {
			isAddedElementDuplicate = true
			break
		}
	}
	if !isAddedElementDuplicate {
		w.buff = append(w.buff, &file)
	}
	if len(w.buff) >= standardBuffFilling {
		w.sendingChangeFileEvent()
	}
}

//Отправка эвентов
func (w *Watcher) sendingChangeFileEvent() {
	if len(w.buff) > 0 {
		eventChangeFile := createNewFileChangingEvent(w.buff)
		w.hub.Broadcast <- eventChangeFile
		w.buff = w.buff[:0]
	}
}

func (w *Watcher) initWatcher() {
	fileWatcher, err := fsnotify.NewWatcher()
	handleErrorAndLogging(err)
	w.Watcher = fileWatcher
	addWatcherFiles(w)
}

//Добавление директорий в список отслеживаемых
func addWatcherFiles(w *Watcher) {
	if w.IsRecursiveWatch {
		filepath.Walk(w.WatchDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				path, err = filepath.Abs(path)
				if err != nil {
					handleErrorAndLogging(err)
					return err
				}
				err = w.Watcher.Add(path)
				if err != nil {
					return err
				}
			}
			return nil
		})
	} else {
		err := w.Watcher.Add(w.WatchDir)
		handleErrorAndLogging(err)
	}
}

func handleErrorAndLogging(err error) {
	if err != nil {
		//TODO добавить обработку или логгирование. При надобности удалить
	}
}

func createNewFileChangingEvent(filesChanging data.FilesTracked) event.ChangeFileEvent {
	return event.NewChangeFileEvent(filesChanging)
}
