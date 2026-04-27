# frontend

Next.js + TypeScript client for Splatmaker.

## Stack

- Next.js (App Router)
- TypeScript
- Redux Toolkit
- RTK Query
- Client components only (for now)

## Screens

- `/login` — username/password form
- `/reconstructions` — list of reconstructions with status
- `/reconstructions/[id]` — details placeholder

## Run

```bash
cd frontend
npm install
npm run dev
```

By default frontend calls API at `http://localhost:8080`.

Override with env var:

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```
