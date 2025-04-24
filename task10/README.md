# Web-интерфейс к поисковому сервису.
## Цель
Разработать frontend микросервиc (генератор HTML страниц), с которым пользователь сможет
взаимодейсвовать через веб-браузер. Альтернативно можно попробовать реализовать Telegram-бота.

Новый микросервис должен позволить клиентам искать картинки комиксов по введенной фразе.
В качестве дополнительной возможности предлагается реализовать страницу администрирования
сервиса (update/drop/stats/status APIs) c авторизацией пользователя.

Результатом должен стать небольшой видеоролик, в котором представлено взаимодействие
пользователя с поисковым кластером с помощью веб-браузера или Telegram-клиента.

Сервисы должны собираться и запускаться через модифицированный compose файл,
а также проходить интеграционные тесты - запуск специального тест контейнера.

## Критерии приемки

1. Все модульные и интеграционные тесты проходят.
2. Видео-ролик на 1-2 минуты снят и внедрен в ABOUT.md файл.

## Материалы для ознакомления

Html templates & cookies

- [Go Web Examples. Templates](https://gowebexamples.com/templates/)
- [Определение и использовние шаблонов](https://metanit.com/go/web/2.1.php)
- [Implementing JWT Authentication In Go](https://permify.co/post/jwt-authentication-go/)
- [A complete guide to working with Cookies in Go](https://www.alexedwards.net/blog/working-with-cookies-in-go)
- [How to Redirect to a URL in Go](https://freshman.tech/snippets/go/http-redirect/)
