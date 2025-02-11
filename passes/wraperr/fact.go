package wraperr

type Kind uint32

const (
	KindUnknown Kind = iota
	KindOK
	KindBad
)

type isErrorHandler struct {
	Kind Kind
}

func (f *isErrorHandler) AFact() {}

func (f *isErrorHandler) String() string {
	switch f.Kind {
	case KindUnknown:
		return "unknownFunc"
	case KindOK:
		return "okFunc"
	case KindBad:
		return "badFunc"
	default:
		panic("unreachable")
	}
}
