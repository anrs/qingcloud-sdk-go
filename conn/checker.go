package conn

type Condition struct {
	Required []interface{}
}

func CheckParams(params map[string]interface{}, cond *Condition) bool {
	return true
}
