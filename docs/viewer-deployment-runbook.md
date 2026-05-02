# Splatmaker Viewer — deployment runbook (без AWS-вызовов)

Этот runbook описывает последовательность действий. Команды AWS (`aws`, `cdk deploy`) запускаешь ты сам.

## 1) Что должно быть готово заранее
- ECR репозиторий для backend viewer image
- DynamoDB таблица джобов пайплайна
- S3 bucket с результатами (splat-файлы)
- Локально: Docker, Node.js, npm, CDK CLI

## 2) Сборка и публикация backend image
Из корня репозитория:

```bash
# пример (подставь свой registry/repo/tag)
docker build -t <ecr_repo_uri>:<tag> ./api
```

Дальше ты сам делаешь `docker push` в ECR.

## 3) Подготовка CDK
```bash
cd infra/cdk
npm install
npm run build
npm run synth
```

Проверить, что synth проходит без ошибок и стек содержит:
- ECS Fargate service
- ALB
- Cognito auth
- IAM read-only доступ к DynamoDB/S3 для viewer

## 4) Деплой CDK (выполняешь сам)
При деплое нужно передать параметры:
- `ApiEcrRepositoryName`
- `ApiImageTag`
- `JobsTableName`
- `ResultBucketName`

После деплоя сохранить output `ViewerURL`.

## 5) Smoke-check после деплоя
- Открыть `ViewerURL` в браузере
- Убедиться, что работает Cognito login flow
- Проверить список джобов
- Открыть детали джобы
- Проверить, что ссылки на result URLs открываются

## 6) Быстрый rollback
- Вернуть предыдущий image tag в параметре `ApiImageTag`
- Повторно обновить стек

## 7) Типичные проблемы
- Пустой список джобов: проверить table name и IAM права на `dynamodb:Scan/GetItem`
- Нет result URLs: проверить bucket name и права `s3:GetObject`
- 401/redirect loop: проверить Cognito client/domain + listener auth action
