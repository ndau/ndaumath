package signature

//go:generate msgp

//msgp:tuple IdentifiedData

// AlgorithmID is an identifier uniquely associated with each supported signature algorithm
type AlgorithmID uint8

// IdentifiedData is a byte slice associated with an algorithm
type IdentifiedData struct {
	Algorithm AlgorithmID
	Data      []byte
}
