package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var (
	DelayDuration = time.Duration(3 * time.Second)
	RetryCount    = -1
)

func Round8(x float64) float64 {
	s := fmt.Sprintf("%.8f", x)
	y, _ := strconv.ParseFloat(s, 64)
	return y
}

func Round4(x float64) float64 {
	s := fmt.Sprintf("%.4f", x)
	y, _ := strconv.ParseFloat(s, 64)
	return y
}

func Float64Round(x float64, prec ...int) float64 {
	precision := 4
	if len(prec) == 1 {
		precision = prec[0]
	}
	format := "%." + strconv.Itoa(precision) + "f"
	s := fmt.Sprintf(format, x)
	y, _ := strconv.ParseFloat(s, 64)
	return y
}

func SetRetryCount(count int) {
	RetryCount = count
}

func SetDelay(delay int) {
	DelayDuration = time.Duration(time.Duration(delay) * time.Millisecond)
}
func RE(method interface{}, params ...interface{}) interface{} {

	invokeM := reflect.ValueOf(method)
	if invokeM.Kind() != reflect.Func {
		panic("method not a function")
		return nil
	}

	var value []reflect.Value = make([]reflect.Value, len(params))
	var i int = 0
	for ; i < len(params); i++ {
		value[i] = reflect.ValueOf(params[i])
	}

	var retV interface{}
	var retryC int = 0

	loop := true
	for {
		retValues := invokeM.Call(value)
		for _, vl := range retValues {
			if vl.Type().String() == "error" {
				if !vl.IsNil() {
					if RetryCount != -1 {
						retryC++
						if retryC <= RetryCount {
						} else {
							loop = false
						}
					}
				} else {
					loop = false
				}
			} else {
				retV = vl.Interface()
			}
		}
		if loop {
			time.Sleep(DelayDuration)
		} else {
			return retV
		}
	}
}
