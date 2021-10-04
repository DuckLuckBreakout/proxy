# proxy



## Запуск

```bash
docker-compose up
```

## Функционал:

* Проксирование HTTP запросов; 
* Проксирование HTTPS запросов; 
* Повторная отправка проксированных запросов;
* Сканер уязвимости.
  * Вариант 4. XSS - во все GET/POST параметры попробовать подставить по очереди
`vulnerable'"><img src onerror=alert()>`. В ответе искать эту же строчку, если нашлась, писать, что данный GET/POST параметр уязвим.

## API

```bash
http://127.0.0.1:8000/requests     - все запросы
http://127.0.0.1:8000/requests/:id - запрос c id
http://127.0.0.1:8000/repeat/:id   - повторить запрос c id
http://127.0.0.1:8000/scan/:id     - сканирование запроса c id
```

