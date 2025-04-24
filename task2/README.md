# Hello gRPC server!
## Цель
Создать микросервис, обслуживающий клиентов по gRPC протоколу, - Pet Name Generator. 

Вам дан proto файл, методы которого нужно реализовать в коде:

1. Generate - генерируем случайное имя по заданным в запросе количеству слов и разделителе
2. GenerateMany - то же самое, но дополнительно задается количество имен. В этом запросе сервер
отвечает потоком имен.

Для генерации случайных имен необходимо использовать библиотеку
[golang-petname](https://github.com/dustinkirkland/golang-petname)

Сервис должен возвращать gRPC ошибку
[InvalidArgument](https://pkg.go.dev/google.golang.org/grpc/codes#pkg-constants), если
пользователь послал количество слов или имен <= 0.

Сервис должен собираться и запускаться через предоставленный compose файл,
а также проходить интеграционные тесты - запуск специального тест контейнера.

Для тестирования сервиса можно использовать [gRPC curl](https://github.com/fullstorydev/grpcurl)
или [gRPC ui](https://github.com/fullstorydev/grpcui)

Полезные curl-ы:

```bash
grpcurl -plaintext  localhost:28081 list
grpcurl -plaintext  localhost:28081 list petname.PetnameGenerator
grpcurl -plaintext  localhost:28081 petname.PetnameGenerator.Ping
grpcurl -plaintext -d '{"words": 3, "separator": " "}' localhost:28081 petname.PetnameGenerator.Generate
```

## Критерии приемки

1. Микросервис компилируeтся в docker image, запускаeтся через compose файл и проходит тесты.
2. Используется библиотека petname.
3. Сервер конфигурируeтся через cleanenv пакет и должeн уметь запускаться как с config.yaml файлом
через флаг -config, так и через переменные среды, в этом задании - PETNAME_GRPC_PORT.
4. Используется golang 1.23+

## Материалы для ознакомления

- [Petname Generator](https://github.com/dustinkirkland/golang-petname)
- [Quick start](https://grpc.io/docs/languages/go/quickstart/)
- [Basics](https://grpc.io/docs/languages/go/basics/)
- [Codes](https://pkg.go.dev/google.golang.org/grpc/codes)
- [gRPC url](https://github.com/fullstorydev/grpcurl)
- [gRPC UI](https://github.com/fullstorydev/grpcui)
