package constants

const (
    MSG_PLACE_ERROR = "{error}"
    MSG_PLACE_URL = "{url}"
    MSG_PLACE_NAME = "{name}"
    MSG_PLACE_PROVIDER = "{provider}"
    MSG_PLACE_PARAMS = "{params}"
    MSG_PLACE_STATUS = "{status}"
    MSG_PLACE_N = "{N}"

    MSG_SUCCESS = "SUCCESS"
    MSG_FAILED = "FAILED"
    MSG_TIMEOUT = "Превышено время ожидания."
    MSG_SET_PARAMS = "Установлены параметры: " + MSG_PLACE_PARAMS
    MSG_PARAMS_NOT_SET = "Параметры: " + MSG_PLACE_PARAMS + " не установлены"
    MSG_PARAMS_REQUIRED = "Требуются параметры: " + MSG_PLACE_PARAMS
    MSG_PARAMS_NOT_ENOUGH = "Пареметров недостаточно"
    MSG_CONNECTION_NOT_EXIST = "Коннекшен не существует"
    MSG_CONNECTION_SKIP_RESPONSE = "Пропускаем запрос № " + MSG_PLACE_N
    MSG_CONNECTION_NOT_PARAMS = "Коннекшен " + MSG_PLACE_NAME + " не имеет параметров запроса"
    MSG_CONNECTION_DISCONNECTED_TO = "Коннекшен " + MSG_PLACE_NAME + " отключен от" + MSG_PLACE_URL
    MSG_CONNECTION_CONNECTED_TO = "Коннекшен " + MSG_PLACE_NAME + " подключен к" + MSG_PLACE_URL
    MSG_CONNECTION_NOT_STARTED = "Коннекшен " + MSG_PLACE_NAME + " не запущен " + MSG_PLACE_PROVIDER
    MSG_CONNECTION_STARTED = "Коннекшен " + MSG_PLACE_NAME + " запущен " + MSG_PLACE_PROVIDER
    MSG_CONNECTION_NOT_STOPED = "Коннекшен " + MSG_PLACE_NAME + " не остановлен " + MSG_PLACE_PROVIDER
    MSG_CONNECTION_STOPED = "Коннекшен " + MSG_PLACE_NAME + " остановлен " + MSG_PLACE_PROVIDER
    MSG_CONNECTION_ACTIVATED = "Коннекшен " + MSG_PLACE_NAME + " активирован"
    MSG_CONNECTION_DEACTIVATED = "Коннекшен " + MSG_PLACE_NAME + " деактивирован"
    MSG_CONNECTION_NOT_ACTIVATED = "Коннекшен " + MSG_PLACE_NAME + " не активирован"
    MSG_MANIFESTS_ERROR = "Терминал. Биржи не инициализированы, ошибка чтения манифестов: " + MSG_PLACE_ERROR + " Исправьте ошибку и перезапустите программу."
)
