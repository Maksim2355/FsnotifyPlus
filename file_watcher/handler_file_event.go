package file_watcher

import (
	"FsnotifyPlus/data"
	"github.com/fsnotify/fsnotify"
	"os"
)

func (w *Watcher) handleChangeFileEvent(ev fsnotify.Event) {
	if ev.Op&fsnotify.Create == fsnotify.Create {
		fi, err := os.Stat(ev.Name)
		if err != nil {
			return
		}
		if fi.IsDir() {
			err = w.Watcher.Add(ev.Name)
			if err == nil {
				//TODO добавить обработку или кастомный логгер
			}
		} else {
			w.addFileChangingInBuffAndSentEventIfNeeded(data.NewFileTracked(ev.Name, data.ADD))
		}
		return
	}
	//Отслеживаем удаление файла
	if ev.Op&fsnotify.Remove == fsnotify.Remove {
		fi, err := os.Stat(ev.Name)
		if err == nil && fi.IsDir() {
			err = w.Watcher.Remove(ev.Name)
			handleErrorAndLogging(err)
		} else {
			w.addFileChangingInBuffAndSentEventIfNeeded(data.NewFileTracked(ev.Name, data.REMOVE))
		}
		return
	}
	if ev.Op&fsnotify.Write == fsnotify.Write {
		w.addFileChangingInBuffAndSentEventIfNeeded(data.NewFileTracked(ev.Name, data.CHANGE))
	}
}
