package utilities

import (
    "reflect"
    "strings"
    // "fmt"
)

func AppendMap(to map[string]string, from map[string]string) map[string]string {
    if to == nil || from == nil {
        return to
    }
    for key, value := range from {
        to[key] = value
    }
    return to
}

func ReplaceValues(data interface{}, newValues interface{}) interface{} {
    if data == nil {
        return nil
    }
    collection := data.(map[string]interface{})
    for key, value := range collection {
        if reflect.TypeOf(value).Kind() == reflect.Map {
            value = ReplaceValues(value, newValues)
        }
        searchValue := SearchValue(newValues, key)
        //fmt.Println("search", key, searchValue)
        if searchValue != nil {
            collection[key] = searchValue
        }
    }
    return collection
}

func SearchValue(data interface{}, fieldName string) interface{} {
    if data == nil {
        return nil
    }
    collection := data.(map[string]interface{})
    for key, value := range collection {
        if reflect.TypeOf(value).Kind() == reflect.Map {
            return SearchValue(value, fieldName)
        }
        if strings.ToLower(key) == strings.ToLower(fieldName) {
            return value
        }
    }
    return nil
}

func ValuesToString(values map[string]string) string {
    str := ""
    for key, value := range values {
        str = str + key + "(" + value + ") "
    }
    return str
}

func GetValues(data interface{}) map[string]string {
    collection := data.(map[string]interface{})
    res := make(map[string]string)
    for key, value := range collection {
        if reflect.TypeOf(value).Kind() == reflect.Map {
            return GetValues(value)
        }
        if reflect.TypeOf(value).Kind() == reflect.String {
            res[strings.ToLower(key)] = value.(string)
        }
        if reflect.TypeOf(value).Kind() == reflect.Float64 {
            res[strings.ToLower(key)] = FloatToString(value.(float64), 8)
        }
    }
    return res
}

func GetValue(data interface{}, fieldsName []string) interface{} {
    if data == nil {
        return nil
    }
    if len(fieldsName) == 0 {
        return nil
    }
    if reflect.TypeOf(data).Kind() == reflect.Slice {
        for _, d := range data.([]interface{}) {
            return GetValue(d, fieldsName)
        }
    } else {
        if len(fieldsName) > 1 {
            return GetValue(data.(map[string]interface{})[fieldsName[0]], fieldsName[1:])
        }
        return data.(map[string]interface{})[fieldsName[0]]
    }
    return nil
}

func ArrayToString(array []interface{}) string {
    if array == nil {
        return ""
    }
    var str string
    for _, value := range array {
        if reflect.TypeOf(value).Kind() == reflect.String {
            str += ToString(value) + " "
        }
        if reflect.TypeOf(value).Kind() == reflect.Map {
            collection := value.(map[string]interface{})
            for key, valueMap := range collection {
                if reflect.TypeOf(valueMap).Kind() == reflect.String {
                    str += key + "(" + ToString(valueMap) + ") "
                }
                if reflect.TypeOf(valueMap).Kind() == reflect.Map {
                    str += ArrayToString([]interface{} { valueMap })
                }
            }
        }
    }
    return str
}
