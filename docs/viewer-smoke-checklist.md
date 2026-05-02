# Splatmaker Viewer — минимальный smoke-check

## A. Backend health
1. `GET /healthz` отвечает 200
2. Нет ошибок в логах контейнера на старте

## B. Auth gate
1. Без логина URL viewer недоступен напрямую
2. ALB/Cognito отправляет на login
3. После логина возвращает в viewer

## C. Jobs list
1. Открывается страница списка
2. Есть элементы (или корректное empty-state)
3. Фильтр по статусу меняет список

## D. Job details
1. Открывается страница по `job_id`
2. Видны статус и базовые поля
3. Есть секция result URLs

## E. Result URLs
1. Ссылки открываются
2. Истекают по TTL (ожидаемо)
3. Нет доступа к чужим/несуществующим ключам

## F. Регрессии на scope viewer-only
Проверить, что отсутствует:
- кнопки/формы запуска джобов
- UI редактирования pipeline settings
- backend endpoints submit/cancel/retry

## G. Быстрый verdict
- PASS: A..F зелёные
- FAIL: любая критичная ошибка в Auth/List/Details/URLs
