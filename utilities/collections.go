package utilities

import (
    "errors"
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
            res[key] = value.(string)
        }
        if reflect.TypeOf(value).Kind() == reflect.Float64 {
            res[key] = FloatToString(value.(float64), 8)
        }
    }
    return res
}

func AddValue(to interface{}, nameValue string, dataValue interface{}) interface{} {
    newValue := make(map[string]interface{})
    newValue[nameValue] = dataValue
    return ReplaceValues(to, newValue)
}

func GetValueByStr(data interface{}, fieldsName []string) (interface{}, error) {
    if data == nil {
        return nil, errors.New("In GetValueByStr data of fieldsName is nil")
    }
    // if len(fieldsName) == 0 {
    //     return nil, errors.New("In GetValueByStr count of fieldsName is 0")
    // }
    if len(fieldsName) > 0 {
        // if reflect.TypeOf(data).Kind() != reflect.Map {
        //     return nil, errors.New("Mismatched type in GetValueByStr")
        // }
        switch reflect.TypeOf(data).Kind() {
        case reflect.Slice:
            for _, d := range data.([]interface{}) {
                return GetValueByStr(d, fieldsName)
            }
        case reflect.Map:
            if data.(map[string]interface{})[fieldsName[0]] != nil {
                if reflect.TypeOf(data.(map[string]interface{})[fieldsName[0]]).Kind() == reflect.Map {
                    return GetValueByStr(data.(map[string]interface{})[fieldsName[0]], fieldsName[1:])
                } else {
                    return data.(map[string]interface{})[fieldsName[0]], nil
                }
            }
        }
        return nil, nil
    } else {
        return data, nil
    }
    // if reflect.TypeOf(data).Kind() == reflect.Slice {
    //     for _, d := range data.([]interface{}) {
    //         return GetValueByStr(d, fieldsName)
    //     }
    // } else {
    //     if len(fieldsName) > 1 {
    //         return GetValueByStr(data.(map[string]interface{})[fieldsName[0]], fieldsName[1:])
    //     }
    //     if reflect.TypeOf(data).Kind() == reflect.Map {
    //         return data.(map[string]interface{})[fieldsName[0]]
    //     }
    //     return nil
    // }
    // return nil
}

func GetValueByInt(data interface{}, fieldsIndex []int64) (interface{}, error) {
    if data == nil {
        return nil, errors.New("In GetValueByInt data of fieldsName is nil")
    }
    if len(fieldsIndex) > 0 {
        // if reflect.TypeOf(data).Kind() != reflect.Slice {
        //     return nil, errors.New("Mismatched type in GetValueByInt")
        // }
        if reflect.TypeOf(data).Kind() == reflect.Slice {
            if data.([]interface{})[fieldsIndex[0]] != nil {
                if reflect.TypeOf(data.([]interface{})[fieldsIndex[0]]).Kind() == reflect.Slice {
                    return GetValueByInt(data.([]interface{})[fieldsIndex[0]], fieldsIndex[1:])
                } else {
                    return data.([]interface{})[fieldsIndex[0]], nil
                }
            }
        }
        return nil, nil
    } else {
        return data, nil
    }
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
