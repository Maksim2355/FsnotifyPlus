package event

//Базовый интерфейс для всех событий. Содержит поле EventName. Все сущности, реализуемые им должны
//вернуть вернуть одну из строк в EventList
type Event interface {
	EventName() string
}

var eventChangeFiles = "ChangeFileEvent"

var eventList = []string{
	eventChangeFiles,
}

//Проверяет, есть ли название события в общем списке событий
func isValidEvent(eventName string) bool {
	for _, event := range eventList {
		if event == eventName {
			return true
		}
	}
	return false
}

//Преобразуем базовое событие в событие об изменении файла
func CastToChangeFilesEvent(e Event) *ChangeFileEvent {
	if !isValidEvent(e.EventName()) {
		return nil
	}
	if v, ok := e.(ChangeFileEvent); ok {
		return &v
	} else {
		return nil
	}
}
