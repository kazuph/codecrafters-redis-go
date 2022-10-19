package main

import (
	"log"
	"time"
)

type Mem struct {
	data map[string]ValueWithExpiry
}

type ValueWithExpiry struct {
	value     string
	expiresAt time.Time
}

func NewMem() *Mem {
	return &Mem{
		data: make(map[string]ValueWithExpiry),
	}
}

func (m *Mem) Get(key string) (string, bool) {
	log.Printf("%#v\n", m.data)
	valueWithExpiry, ok := m.data[key]

	if !ok {
		log.Println("GET: !ok")
		return "", false
	}

	if !valueWithExpiry.expiresAt.IsZero() && time.Now().Before(valueWithExpiry.expiresAt) {
		delete(m.data, key)
		log.Println("GET: expired")
		return "", false
	}

	return valueWithExpiry.value, true
}

func (m *Mem) Set(key, value string) {
	m.data[key] = ValueWithExpiry{value: value}
}

func (m *Mem) SetWithExpiry(key, value string, expiry time.Duration) {
	m.data[key] = ValueWithExpiry{
		value:     value,
		expiresAt: time.Now().Add(expiry),
	}
}
