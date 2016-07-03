package jerrors

// errors
var (
	ErrNotFound   = JError("not found")
	ErrNoNewBlock = JError("there is no new block ahead")
)

// JError is jackpot custom error
type JError string

func (j JError) Error() string {
	return string(j)
}
