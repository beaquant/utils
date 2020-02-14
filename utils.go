package utils

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

var (
	defaultPrecision = 4
)

func ToFloat64(v interface{}) float64 {
	if v == nil {
		return 0.0
	}

	switch v.(type) {
	case float64:
		return v.(float64)
	case string:
		vStr := v.(string)
		vF, _ := strconv.ParseFloat(vStr, 64)
		return vF
	default:
		panic("to float64 error.")
	}
}

func ToInt(v interface{}) int {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case string:
		vStr := v.(string)
		vInt, _ := strconv.Atoi(vStr)
		return vInt
	case int:
		return v.(int)
	case float64:
		vF := v.(float64)
		return int(vF)
	default:
		panic("to int error.")
	}
}

func ToUint64(v interface{}) uint64 {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case int:
		return uint64(v.(int))
	case float64:
		return uint64((v.(float64)))
	case string:
		uV, _ := strconv.ParseUint(v.(string), 10, 64)
		return uV
	default:
		panic("to uint64 error.")
	}
}

func ToInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case float64:
		return int64(v.(float64))
	default:
		vv := fmt.Sprint(v)

		if vv == "" {
			return 0
		}

		vvv, err := strconv.ParseInt(vv, 0, 64)
		if err != nil {
			return 0
		}

		return vvv
	}
}

func Float64Round2(x float64, prec ...int) float64 {
	precision := defaultPrecision
	if len(prec) == 1 {
		precision = prec[0]
	}
	s := Float64RoundString(x, precision)
	y, _ := strconv.ParseFloat(s, 64)
	return y
}

func Float64Round3(x float64, prec ...int) float64 {
	precision := defaultPrecision
	if len(prec) == 1 {
		precision = prec[0]
	}
	format := "%." + strconv.Itoa(precision) + "f"
	s := fmt.Sprintf(format, x)
	y, _ := strconv.ParseFloat(s, 64)
	return y
}

func Float64Round(f float64, prec ...int) float64 {
	precision := defaultPrecision
	if len(prec) == 1 {
		precision = prec[0]
	}
	x := math.Pow10(precision)
	return math.Trunc(f*x) / x
}

func Float64RoundString(x float64, prec ...int) string {
	precision := defaultPrecision
	if len(prec) == 1 {
		precision = prec[0]
	}
	s := strconv.FormatFloat(x, 'f', precision, 64)
	return s
}

func Float64RoundString2(x float64, prec ...int) string {
	precision := defaultPrecision
	if len(prec) == 1 {
		precision = prec[0]
	}
	format := "%." + strconv.Itoa(precision) + "f"
	s := fmt.Sprintf(format, x)
	return s
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func Float64ToString2(f float64) string {
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
