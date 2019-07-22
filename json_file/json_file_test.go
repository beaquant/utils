package json_file

import "testing"

type C1 struct {
	Status int           `json:"status"`
	Data   []interface{} `json:"data"`
}

func TestJsonFile_Load(t *testing.T) {
	c := &C1{}
	t.Log(Load("json_file.json", c))
	t.Log(c)
}

func TestJsonFile_Save(t *testing.T) {
	c := &C1{}
	Load("json_file.json", c)
	t.Log(Save("json_file_save.json", c))
}
