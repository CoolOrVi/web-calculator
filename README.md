# web-calculator
Этот проект - веб-сервис подсчета арифметических выражений.

## Функциональность
* Вычисление выражений с многозначными числами
* Поддерживаются базовые операции(+, - , / , *) и операции приоретизации ( и )
* Калькулятор __не поддерживает__:
    * Выражения с пропущенным знаком умножения перед скобкой
    * Операции отличные от перечисленных выше (извлечение корня, возведение в степень)

## Запуск
1. Клонировать репозиторий с помощью команды:\
__git clone https://github.com/coolorvi/web-calculator.git__
2. Перейти в папку проекта и запустить проект командой: __go run ./cmd/main.go__

## Примеры использования
* Запрос с успешым вычислением выражения(код 200):\
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"(1984 + (2030 - 1918)) * 404"
}' Ответ сервера:
{"result":"846784.000000"}
-----
* Запрос с невалидным выражением(код 422):\
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"333 - 33a"
}' Ответ сервера: {"error":"Expression is not valid"}
----
* Запрос с ошибкой сервера(код 500):  
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"10 / 0"
}' Ответ сервера: {"error":"Internal server error"}


