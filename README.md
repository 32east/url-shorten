# Сокращатель ссылок
Базовый сокращатель ссылок.<br>Ничего лишнего, только одно поле и кнопка.<br>
Чисто для прикола написал за 1 день.

## API
Все API запросы (в данный момент только один) имеют лимит использования.<br>
Если количество запросов превышает 5 за 30 секунд, пользователю возвращается ошибка с просьбой попробовать позже.

- ``POST /api/v1/shorten`` - позволяет сократить ссылку.<br>
  Пример запроса: ``{"url":"https://www.youtube.com/watch?v=OvLgTjz1jmc"}``<br>
  Пример ответа: ``{"success":true,"url":"QQQQQE9"}``

## Как запустить?
Инструкция написана под операционную систему на базе Linux.
1) Сперва вы должны убедиться, что у вас установлен docker-compose:
   - ``apt update & apt upgrade & apt install docker-compose``
2) Загрузите данный репозиторий:
   - ``git clone https://github.com/tttttt30/url-shorten``
3) Перейдите по пути `/web-forum/build/`:
   - ``cd /web-forum/build/``
4) Запустите компиляцию и запуск проекта, благодаря docker-compose:
   - ``docker-compose up -d``
5) После этого, у вас запустятся только 2 контейнера, поскольку не была загружена база данных. Поэтому, после запуска, прописываем:
   - ``docker exec -ti postgres psql -U postgres -h localhost -c 'create database "url-shorten";'``<br>
   - ``docker exec -ti postgres psql -U postgres -h localhost -f /usr/src/app/migrations/0001_mysql_initialize.sql url-shorten``
6) Запустите остальные контейнеры:
   - ``docker-compose up -d``

## Как выглядит дизайн:
![img.png](assets/img.png)
Я же говорил, ничего лишнего ◕‿↼