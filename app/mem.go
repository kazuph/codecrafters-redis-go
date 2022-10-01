package main

type Mem struct {
	data map[string]string
}

func NewMem() *Mem {
	return &Mem{
		data: make(map[string]string),
	}
}

func (m *Mem) Set(key, value string) {
	m.data[key] = value
}

func (m *Mem) Get(key string) string {
	return m.data[key]
}
