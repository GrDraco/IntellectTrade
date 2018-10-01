package connections

import (
    // "fmt"
    "net/http"
    "errors"
    "strings"
    //"reflect"
    "encoding/json"
    // "time"
    "io/ioutil"
    "../../../utilities"
)

type Api struct {
    // Наследования
    BaseConnection
    // Свойства
    client *http.Client
    request *http.Request
    err error
}
//api.manifest.RequestJSON = utilities.ReplaceValues(api.manifest.RequestJSON, data[0])
func NewApi(manifest *Manifest) *Api {
    // Выделение памяти под сокет
    api := &Api {}
    api.client = &http.Client {}
    // Инициализация функци работы запроса данных
    api.init = func() error {
        var url = api.manifest.URL//"?symbol=BTC-USDT"
        // Формируем url параметрами указаными в запросе манифеста
        values := utilities.GetValues(api.manifest.RequestJSON)
        if len(values) > 0 {
            url = url + "?"
        }
        for key, value := range values {
            place := "{" + key + "}"
            if i, _ := utilities.SearchIndex(url, place); i > -1 {
                url = strings.Replace(url, place, value, -1)
            } else {
                url = url + key + "=" + value + "&"
            }
        }
        // Формируем запрос
        api.request, api.err = http.NewRequest("GET", url, nil)
        if api.err != nil {
    		return api.err
    	}
        return nil
    }
    api.send = func() error {
        if api.request == nil {
            return errors.New(api.manifest.Messages["CONNECTION_NOT_PARAMS"])
        }
        var resp *http.Response
        resp, api.err = api.client.Do(api.request)
        if api.err != nil {
    		return api.err
    	}
        var body []byte
    	body, api.err = ioutil.ReadAll(resp.Body)
        if api.err != nil {
    		return api.err
    	}
        // Сторку конвертируем в объект и записываем в манифест
        json.Unmarshal(body, &api.manifest.Response.JSON)
        // Закрываем чтение ответа
        resp.Body.Close()
        return nil
    }
    api.close = func() error {
        return nil
    }
    // Активируем базовый функционал
    if api.activate(manifest) {
        return api
    }
    return nil
}
