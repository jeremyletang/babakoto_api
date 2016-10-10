package jsend

const (
	success = "success"
	fail    = "fail"
	error   = "error"
)

func New(data interface{}) map[string]interface{} {
	jsend := map[string]interface{}{}
	jsend["status"] = success
	jsend["data"] = data
	return jsend
}

func WithName(data interface{}, name string) map[string]interface{} {
	data_ := map[string]interface{}{}
	data_[name] = data
	return New(data_)
}

func Fail(data interface{}) map[string]interface{} {
	jsend := map[string]interface{}{}
	jsend["status"] = fail
	jsend["data"] = data
	return jsend
}

func FailWithName(data interface{}, name string) map[string]interface{} {
	data_ := map[string]interface{}{}
	data_[name] = data
	return Fail(data_)
}

// only message is mandatory
func Error(message string) map[string]interface{} {
	jsend := map[string]interface{}{}
	jsend["status"] = error
	jsend["message"] = message
	return jsend
}

func ErrorWithData(message string, data interface{}) map[string]interface{} {
	jsend := Error(message)
	jsend["data"] = data
	return jsend
}

func ErrorWithDataAndName(message string, data interface{}, name string) map[string]interface{} {
	data_ := map[string]interface{}{}
	data_[name] = data
	return ErrorWithData(message, data_)
}
