package dberrors

type CustomError struct {
    Err string
}


func (e *CustomError) Error() string {
	return e.Err
}