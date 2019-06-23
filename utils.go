package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
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

func Float64RoundString(x float64, prec ...int) string {
	precision := 4
	if len(prec) == 1 {
		precision = prec[0]
	}
	format := "%." + strconv.Itoa(precision) + "f"
	s := fmt.Sprintf(format, x)
	return s
}

func Float64ToString(f float64) string {
	return fmt.Sprint(f)
}

func StringToFloat64(s string) float64 {
	y, _ := strconv.ParseFloat(s, 64)
	return y
}

type RE struct {
	delayDuration time.Duration
	retryCount    int
}

func NewRE() *RE {
	return &RE{delayDuration: time.Duration(3 * time.Second), retryCount: -1}
}

func (r *RE) SetRetryCount(count int) {
	r.retryCount = count
}

func (r *RE) SetDelay(delay int) {
	r.delayDuration = time.Duration(time.Duration(delay) * time.Millisecond)
}

func (r *RE) RE(method interface{}, params ...interface{}) interface{} {

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
					if r.retryCount != -1 {
						retryC++
						if retryC <= r.retryCount {
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
			time.Sleep(r.delayDuration)
		} else {
			return retV
		}
	}
}
