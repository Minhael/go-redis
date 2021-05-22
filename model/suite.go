package model

type Suite interface {
	Execute() (string, error)
}
