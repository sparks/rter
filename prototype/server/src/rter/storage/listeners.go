package storage

type DatabaseListener interface {
	InsertEvent(v interface{})
	UpdateEvent(v interface{})
	DeleteEvent(v interface{})
}

type ListenerSlice []DatabaseListener

var listeners ListenerSlice = make([]DatabaseListener, 0)

func (s *ListenerSlice) Add(l DatabaseListener) {
	*s = append(*s, l)
}

func (s *ListenerSlice) Remove(l DatabaseListener) bool {
	for i, c := range *s {
		if c == l {
			*s = append((*s)[0:i], (*s)[i+1:len(*s)]...)
			return true
		}
	}

	return false
}

func (s *ListenerSlice) NotifyInsert(val interface{}) {
	for _, c := range *s {
		c.InsertEvent(val)
	}
}

func (s *ListenerSlice) NotifyUpdate(val interface{}) {
	for _, c := range *s {
		c.UpdateEvent(val)
	}
}

func (s *ListenerSlice) NotifyDelete(val interface{}) {
	for _, c := range *s {
		c.DeleteEvent(val)
	}
}

func AddListener(l DatabaseListener) {
	listeners.Add(l)
}

func RemoveListener(l DatabaseListener) bool {
	return listeners.Remove(l)
}
