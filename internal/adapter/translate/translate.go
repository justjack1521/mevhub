package translate

type Translator[I any, O any] interface {
	Marshall(data I) (out O, err error)
	Unmarshall(data O) (out I, err error)
}
