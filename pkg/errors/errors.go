package errors

type Error interface {
	Reason() int
	ShortErr() string
	LongErr() string
}

type errorType struct {
	reason int
	short  string
	long   string
}

func (e errorType) Reason() int {
	return e.reason
}

func (e errorType) ShortErr() string {
	return e.short
}

func (e errorType) LongErr() string {
	return e.long
}

func CreateError(reason int, short, long string) Error {
	return &errorType{
		reason: reason,
		short:  short,
		long:   long,
	}
}
