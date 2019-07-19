package json_file

import "testing"

type C1 struct {
	Status int           `json:"status"`
	Data   []interface{} `json:"data"`
}

func TestJsonFile_Load(t *testing.T) {
	c := &C1{}
	jf := &JsonFile{}
	t.Log(jf.Load("json_file.json", c))
	t.Log(c)
}

func TestJsonFile_Save(t *testing.T) {
	c := &C1{}
	jf := &JsonFile{}
	jf.Load("json_file.json", c)
	t.Log(jf.Save("json_file_save.json", c))
}
