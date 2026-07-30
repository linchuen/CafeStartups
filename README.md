# Café Startups Web

Go + React 的《Café Startups》數位桌遊 MVP。

## Phase 0

目前完成規格與技術骨架：

- Go HTTP server：`GET /health`、`GET/POST /api/games`
- React + TypeScript + Vite 首頁與建立房間流程
- 卡牌 JSON Schema 與 MVP fixture
- Go 基礎測試與前端 typecheck/build scripts

## 啟動

### Backend

```powershell
go run ./cmd/server
```

API 會在 `http://localhost:8080` 啟動。

### Frontend

```powershell
cd client
npm install
npm run dev
```

前端會在 `http://localhost:5173` 啟動。

## 驗證

```powershell
go test ./...
cd client
npm run typecheck
npm run build
```

規則與後續階段請參閱 [Agent.md](./Agent.md) 與 [MVP 階段式需求](./docs/MVP-階段式需求.md)。
