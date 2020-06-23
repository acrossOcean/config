package config

type ReadOrder int

const (
	ReadFromParam ReadOrder = iota
	ReadFromEnv
	ReadFromFile
)

func DefaultReadOrder() []ReadOrder {
	return []ReadOrder{ReadFromParam, ReadFromFile, ReadFromEnv}
}
