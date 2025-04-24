# Нормализация поисковых запросов, часть 1
## Цель
Создать микросервис, обслуживающий клиентов по gRPC протоколу, - Words Normalizer. 

Mикросервис нормализации слов должен работать в соответсвии с предложенным proto-файлом.
Сервис должен принимать на вход строку (на английском) и возвращать назад нормализованный вид
в виде слайса слов. То есть при посылке "follower brings bunch of questions" сервер должен отдать
["follow", "bring", "bunch", "question"] - слова в слайсе в любом порядке.

Приложение должно отсеивать часто употребляемые слова типа of/a/the/, местоимения
и глагольные частицы (will).

Для нормализации необходимо использовать библиотеку
[snowball](github.com/kljensen/snowball)

Сервис должен возвращать gRPC ошибку при получении сообщения больше 4 KiB - 
[ResourceExhausted](https://pkg.go.dev/google.golang.org/grpc/codes#pkg-constants).

Сервис должен собираться и запускаться через предоставленный compose файл,
а также проходить интеграционные тесты - запуск специального тест контейнера.

Для тестирования сервиса можно использовать [gRPC curl](https://github.com/fullstorydev/grpcurl)
или [gRPC ui](https://github.com/fullstorydev/grpcui)

## Критерии приемки

1. Микросервис компилируeтся в docker image, запускаeтся через compose файл и проходит тесты.
2. Используется библиотека snowball, код логики нормализации находится
в папке search-services/words/words.
3. Сервер конфигурируeтся через cleanenv пакет и должeн уметь запускаться как с config.yaml файлом
через флаг -config, так и через переменные среды, в этом задании - WORDS_GRPC_PORT.
4. Используется golang 1.23+

## Материалы для ознакомления

- [Как создавать модули](https://go.dev/doc/tutorial/create-module)
- [Библиотека для нормализации](github.com/kljensen/snowball)
- [Библиотека для нормализации, English + stopwords](https://github.com/kljensen/snowball/tree/master/english)
