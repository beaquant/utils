package json_file

import "testing"

type C1 struct {
	Status int           `json:"status"`
	Data   []interface{} `json:"data"`
}

func TestConfig_Read(t *testing.T) {
	c := &C1{}
	jf := &JsonFile{}
	t.Log(jf.Read("json_file.json", c))
	t.Log(c)
}
