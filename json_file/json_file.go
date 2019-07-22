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

func readFile(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func confirmConfigJSON(file []byte, result interface{}) error {
	if !strings.Contains(reflect.TypeOf(result).String(), "*") {
		return errors.New(errNotAPointer)
	}
	return jsonDecode(file, &result)
}

func jsonDecode(data []byte, to interface{}) error {
	if !strings.Contains(reflect.ValueOf(to).Type().String(), "*") {
		return errors.New("json decode error - memory address not supplied")
	}
	return json.Unmarshal(data, to)
}

func GetExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func GetOSPathSlash() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}

func Load(configPath string, result interface{}) error {
	file, err := readFile(configPath)
	if err != nil {
		fmt.Errorf("common.ReadFile err:%s", err)
		return err
	}

	err = confirmConfigJSON(file, result)
	if err != nil {
		fmt.Errorf("confirmConfigJSON:%s", err)
		return err
	}

	return nil
}

func Save(configPath string, cfg interface{}) error {
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
