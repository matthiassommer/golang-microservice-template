package utils

//HandleError uses given error to print int to console
func HandleError(err error) error {
	Log.Error(err)
	return err
}
