package jerrors

// errors
var (
	ErrNotFound = JError("not found")
)

// JError is jackpot custom error
type JError string

func (j JError) Error() string {
	return string(j)
}
