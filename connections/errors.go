package connections

type TypeAssertionError struct{}

func (TypeAssertionError) Error() string {
	return "Couldn't conver to type!"
}
