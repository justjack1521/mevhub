package translate

type Marshaller[I any, O any] interface {
	Marshall(data I) (out O, err error)
}

type Unmarshaller[I any, O any] interface {
	Unmarshall(data I) (out O, err error)
}

type Translator[I any, O any] interface {
	Marshaller[I, O]
	Unmarshaller[O, I]
}
