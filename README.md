
# Сервис проверки клиентов по Перечню в целях ПОД/ФТ

Данный сервер является рузельтатом работы магистерской диссертации.

Для работы требуется директория dist с собранным [графическим интерфейсом](http://sabaka.net), база данных клиентов sql и Перечень в формате xml.




## Сборка

Для сборки необходим установленный язык [Go](https://go.dev/).

Для того, чтобы собрать исполняемый файл под Windows используйте команду:

```bash
  make build
```

## Запуск

После запустите файл main.exe. Если добавить флаг -b при запуске, сразу откроется браузер с графическим интерфейсом.

```bash
  ./main.exe -b
```