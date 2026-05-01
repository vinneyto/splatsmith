## backend (Go module)

Новый модуль `backend`:
- принимает ссылку на `.ply`
- скачивает и парсит ASCII PLY (vertex x y z)
- строит равномерный octree (ячейки одинакового размера, ограничение `cell_size` + `max_depth`)
- отдает структуру дерева и индексы сплатов в ячейках
- транспорт: gRPC (proto в `proto/backend/v1/backend.proto`)

### Быстрый старт

```bash
cd backend
go test ./...
go run ./cmd/backend
```

Сервер по умолчанию слушает `:9090`.
