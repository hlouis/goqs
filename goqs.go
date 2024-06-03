package goqs

type QSType map[interface{}]interface{}

type QSMap[T int | string] map[T]any
