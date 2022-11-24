package err

type Error interface {
	Reason() int
	ShortErr() string
	LongErr() string
}
