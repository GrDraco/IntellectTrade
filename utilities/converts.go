package utilities

import (
    "reflect"
    "strconv"
)

func ToString(value interface{}) string {
    if value == nil {
        return ""
    }
    switch reflect.TypeOf(value).Kind() {
    case reflect.String:
        return value.(string)
    case reflect.Int64:
        return IntToString(value.(int64))
    case reflect.Float64:
        return FloatToString(value.(float64), -1)
    }
    return "error ToString"
}

func IntToString(value int64) string {
    return strconv.FormatInt(value, 10)
}

func FloatToString(input_num float64, modulo int) string {
    return strconv.FormatFloat(input_num, 'f', modulo, 64)
}

func ToInt(value interface{}) int64 {
    if value == nil {
        return 0
    }
    if reflect.TypeOf(value).Kind() == reflect.String {
        if value.(string) == "" {
            return 0
        } else {
            if res, err := strconv.ParseInt(value.(string), 10, 64); err == nil {
                return res
        	}
        }
    }
    if reflect.TypeOf(value).Kind() == reflect.Int64 {
        return value.(int64)
    }
    if reflect.TypeOf(value).Kind() == reflect.Float64 {
        return int64(value.(float64))
    }
    return 0
}
func ToFloat(value interface{}) float64 {
    if value == nil {
        return 0
    }
    if reflect.TypeOf(value).Kind() == reflect.String {
        if value.(string) == "" {
            return 0
        } else {
            if res, err := strconv.ParseFloat(value.(string), 64); err == nil {
                return res
        	}
        }
    }
    if reflect.TypeOf(value).Kind() == reflect.Float64 {
        return value.(float64)
    }
    if reflect.TypeOf(value).Kind() == reflect.Int64 {
        return float64(value.(int64))
    }
    return 0
}
