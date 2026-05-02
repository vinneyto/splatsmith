# Splatmaker Viewer — env/config cheat sheet

Минимальный набор переменных для backend в режиме AWS:

```env
SPLATMAKER_MODE=aws
SPLATMAKER_API_ADDR=:8080
SPLATMAKER_AWS_REGION=<aws-region>
SPLATMAKER_AWS_JOBS_TABLE=<jobs-table-name>
SPLATMAKER_AWS_RESULT_BUCKET=<results-bucket-name>
```

## Что за что отвечает
- `SPLATMAKER_MODE`: должен быть `aws` для DynamoDB/S3 адаптеров
- `SPLATMAKER_API_ADDR`: bind address API контейнера
- `SPLATMAKER_AWS_REGION`: регион, где лежат таблица и bucket
- `SPLATMAKER_AWS_JOBS_TABLE`: таблица джобов пайплайна
- `SPLATMAKER_AWS_RESULT_BUCKET`: bucket с файлами результатов

## IAM для task role (минимум)
DynamoDB:
- `dynamodb:Scan`
- `dynamodb:GetItem`
- (опционально) `dynamodb:Query`

S3:
- `s3:GetObject` на `arn:aws:s3:::<bucket>/*`

## Не нужно для viewer
- Любые переменные/секреты, связанные с запуском пайплайна
- PipelineSettings в env
- write-права на DynamoDB/S3

## Standalone режим (для локальной проверки)
```env
SPLATMAKER_MODE=standalone
SPLATMAKER_API_ADDR=:8080
# и локальные standalone-пути из текущего конфига проекта
```
