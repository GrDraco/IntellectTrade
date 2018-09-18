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
    Name string                                             `json:"name"`
    IsDefault bool                                          `json:"is_default"`
    ParamsExchanges map[string]map[string]map[string]string `json:"params_exchanges"`
    StartedExchanges map[string]map[string]bool             `json:"started_exchanges"`
    StartedStartegies map[string]bool                       `json:"started_strategies"`
}
// Функция создания нового объекта настроеек
func NewSettings(name string) *Settings {
    settings := new(Settings)
    settings.Name = name
    settings.IsDefault = false
    settings.ParamsExchanges = make(map[string]map[string]map[string]string)
    settings.StartedExchanges = make(map[string]map[string]bool)
    settings.StartedStartegies = make(map[string]bool)
    return settings
}
func (settings *Settings) Init() {
    if settings.ParamsExchanges == nil {
        settings.ParamsExchanges = make(map[string]map[string]map[string]string)
    }
    if settings.StartedExchanges == nil {
        settings.StartedExchanges = make(map[string]map[string]bool)
    }
    if settings.StartedStartegies == nil {
        settings.StartedStartegies = make(map[string]bool)
    }
}

func (settings *Settings) Load() error {
    return settings.LoadFrom(settings.Name)
}

func (settings *Settings) LoadFrom(name string) error {
    newSettings, err := ReadSettings(PATH_SETTINGS + name + ".json")
    settings.Name = newSettings.Name
    settings.IsDefault = newSettings.IsDefault
    settings.ParamsExchanges = newSettings.ParamsExchanges
    settings.StartedExchanges = newSettings.StartedExchanges
    settings.StartedStartegies = newSettings.StartedStartegies
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
    // Выставляем параметры согласно настроек
    for setExchange, entities := range settings.ParamsExchanges {
        exchange := terminal.Exchanges[setExchange]
        if exchange != nil {
            for entity, properties := range entities {
                if properties != nil {
                    params := make(map[string]interface{})
                    for param, value := range properties {
                        params[param] = value
                    }
                    res := exchange.SetValues(entity, params)
                    if !res {
                        return false
                    }
                }
            }
        }
    }
    // Запускаем биржи согласно настроек
    for setExchange, entities := range settings.StartedExchanges {
        exchange := terminal.Exchanges[setExchange]
        if exchange != nil {
            for entity, value := range entities {
                if value {
                    _, success := exchange.Turn(entity)
                    if !success {
                        return false
                    }
                }
            }
        }
    }
    // Запускаем стратегии согласно настроек
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
