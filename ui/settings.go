package ui

import (
    "os"
    "errors"
    "strings"
    "io/ioutil"
    "encoding/json"

    "../market"
)

const (
    PATH_SETTINGS = "./settings/"

    MSG_SETTINGS_NO_NAME = "Отсутсвет имя настроеек"
)

type Settings struct {
    Name string                                                         `json:"name"`
    IsDefault bool                                                      `json:"is_default"`
    ParamsMain map[string]string                                        `json:"params_main"`
    ParamsExchanges map[string]map[string]map[string]map[string]string  `json:"params_exchanges"`
    StartedExchanges map[string]map[string]map[string]bool                         `json:"started_exchanges"`
    ParamsStartegies map[string]map[string]interface{}                  `json:"params_strategies"`
}
// Функция создания нового объекта настроеек
func NewSettings(name string) *Settings {
    settings := new(Settings)
    settings.Name = name
    settings.IsDefault = false
    settings.ParamsMain = make(map[string]string)
    settings.ParamsExchanges = make(map[string]map[string]map[string]map[string]string)
    settings.StartedExchanges = make(map[string]map[string]map[string]bool)
    settings.ParamsStartegies = make(map[string]map[string]interface{})
    return settings
}
func (settings *Settings) Init() {
    if settings.ParamsMain == nil {
        settings.ParamsMain = make(map[string]string)
    }
    if settings.ParamsExchanges == nil {
        settings.ParamsExchanges = make(map[string]map[string]map[string]map[string]string)
    }
    if settings.StartedExchanges == nil {
        settings.StartedExchanges = make(map[string]map[string]map[string]bool)
    }
    if settings.ParamsStartegies == nil {
        settings.ParamsStartegies = make(map[string]map[string]interface{})
    }
}

func (settings *Settings) Load() error {
    return settings.LoadFrom(settings.Name)
}

func (settings *Settings) LoadFrom(name string) error {
    newSettings, err := ReadSettings(PATH_SETTINGS + name + ".json")
    settings.Name = newSettings.Name
    settings.IsDefault = newSettings.IsDefault
    settings.ParamsMain = newSettings.ParamsMain
    settings.ParamsExchanges = newSettings.ParamsExchanges
    settings.StartedExchanges = newSettings.StartedExchanges
    settings.ParamsStartegies = newSettings.ParamsStartegies
    if err != nil {
        return err
    }
    return nil
}
// Запись насироек на диск
func (settings *Settings) Save() error {
    return settings.SaveAs(settings.Name)
}
func (settings *Settings) SaveAs(name string) error {
    if name == "" {
        return errors.New(MSG_SETTINGS_NO_NAME)
    }
    lastName := settings.Name
    settings.Name = name
    if err := WriteSettings(settings, PATH_SETTINGS + name + ".json"); err != nil {
        return err
    }
    if settings.IsDefault && lastName != name {
        var settingsTemp *Settings
        var err error
        if settingsTemp, err = ReadSettings(PATH_SETTINGS + lastName + ".json"); err != nil {
            return err
        }
        settingsTemp.IsDefault = false
        if err := WriteSettings(settings, PATH_SETTINGS + lastName + ".json"); err != nil {
            return err
        }

    }
    return nil
}
// Чтение настроек из файла на диске
func (settings *Settings) Apply(terminal *market.Terminal) bool {
    if settings.ParamsExchanges == nil {
        return false
    }
    // Запускаем стратегии согласно настроек
    for setStartegy, entities := range settings.ParamsStartegies {
        startegy := terminal.Strategies[setStartegy]
        if startegy != nil {
            for property, value := range entities {
                if value != nil {
                    if !startegy.SetProperty(property, value) {
                        return false
                    }
                }
            }
        }
    }
    // Выставляем параметры согласно настроек
    for setExchange, entities := range settings.ParamsExchanges {
        exchange := terminal.Exchanges[setExchange]
        if exchange != nil {
            for entity, providers := range entities {
                for provider, properties := range providers {
                    if properties != nil {
                        params := make(map[string]interface{})
                        for param, value := range properties {
                            params[param] = value
                        }
                        if !exchange.SetValues(entity, provider, params) {
                            return false
                        }
                    }
                }
            }
        }
    }
    // Запускаем биржи согласно настроек
    for setExchange, entities := range settings.StartedExchanges {
        exchange := terminal.Exchanges[setExchange]
        if exchange != nil {
            for entity, providers := range entities {
                for provider, value := range providers {
                    _, success := exchange.Turn(entity, provider, value)
                    if !success {
                        return false
                    }
                }
            }
        }
    }
    return true
}
//
func GetSettingsPaths() (paths map[string]string, err error) {
    var files []os.FileInfo
    files, err = ioutil.ReadDir(PATH_SETTINGS)
    if err != nil {
        return
    }
    paths = make(map[string]string)
    for _, f := range files {
        paths[strings.Split(f.Name(), ".")[0]] = PATH_SETTINGS + f.Name()
    }
    return
}
// Чтение настроек из файла на диске
func ReadSettings(path string) (settings *Settings, err error) {
    var jsonFile []byte
    jsonFile, err = ioutil.ReadFile(path)
    if err != nil {
		return
	}
    settings = new(Settings)
    json.Unmarshal(jsonFile, settings)
    return
}
// Запись настроек в файл
func WriteSettings(settings *Settings, path string) error {
    var settingsJson []byte
    var err error
    settingsJson, err = json.Marshal(settings)
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(path, settingsJson, 0644)
    if err != nil {
        return err
    }
    return nil
}
// Получение всех имеющихся настроеек
func GetSettings() (allSettings map[string]*Settings, err error) {
    var paths map[string]string
    var settings *Settings
    paths, err = GetSettingsPaths()
    if err != nil {
		return
	}
    allSettings = make(map[string]*Settings)
    for _, path := range paths {
        settings, err = ReadSettings(path)
        if err != nil {
    		return
    	}
        allSettings[settings.Name] = settings
    }
    return
}
