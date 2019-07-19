package json_file

import (
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

const (
	EncryptedConfigFile          = "json_file.dat"
	ArbitrerConfigFile           = "arbitrer_config.json"
	ExchangeConfigFile           = "exchange_config.json"
	ConfigTestFile               = "example_config.json"
	configFileEncryptionPrompt   = 0
	configFileEncryptionEnabled  = 1
	configFileEncryptionDisabled = -1
	ErrFailureOpeningConfig      = "Fatal error opening %s file. Error: %s"
	errNotAPointer               = "Error: parameter interface is not a pointer"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type JsonFile struct {
}

func (jf *JsonFile) readFile(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (jf *JsonFile) confirmConfigJSON(file []byte, result interface{}) error {
	if !strings.Contains(reflect.TypeOf(result).String(), "*") {
		return errors.New(errNotAPointer)
	}
	return jf.jsonDecode(file, &result)
}

func (jf *JsonFile) jsonDecode(data []byte, to interface{}) error {
	if !strings.Contains(reflect.ValueOf(to).Type().String(), "*") {
		return errors.New("json decode error - memory address not supplied")
	}
	return json.Unmarshal(data, to)
}

func (jf *JsonFile) GetExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func (jf *JsonFile) GetOSPathSlash() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}

func (jf *JsonFile) Load(configPath string, result interface{}) error {
	file, err := jf.readFile(configPath)
	if err != nil {
		fmt.Errorf("common.ReadFile err:%s", err)
		return err
	}

	err = jf.confirmConfigJSON(file, result)
	if err != nil {
		fmt.Errorf("confirmConfigJSON:%s", err)
		return err
	}

	return nil
}

func (jf *JsonFile) Save(configPath string, cfg interface{}) error {
	data, err := json.MarshalIndent(cfg, "", "    ") //这里返回的data值，类型是[]byte
	if err != nil {
		fmt.Errorf("common.SaveFile err:%s", err)
		return err
	}
	_, err = writeBytes(configPath, data)
	if err != nil {
		fmt.Errorf("confirmConfigJSON:%s", err)
		return err
	}
	return nil
}

func writeBytes(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}
