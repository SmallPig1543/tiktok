package model

type EsModel interface {
	Index() string
	Mapping() string
}
