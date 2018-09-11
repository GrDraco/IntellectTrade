package utilities

import (
    "reflect"
    "strconv"
)

func ToString(value interface{}) string {
    if value == nil {
        return ""
    }
    return value.(string)
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
    return 0
}
