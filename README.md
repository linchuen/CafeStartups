# Café Startups

Go + React 的《Café Startups》單機桌遊 MVP。MVP 只讓一位真人在本機遊玩，其餘席位由不考慮策略、只做合法隨機選擇的 MVP 電腦玩家補足；玩法確認正確後，才開發區網模式，最後才開發線上模式。

## Phase 0

目前完成規格與技術骨架：

- Go HTTP server：`GET /health`、`GET/POST /api/games`
- React + TypeScript + Vite 首頁與建立房間流程
- 卡牌 JSON Schema 與 MVP fixture
- Go 基礎測試與前端 typecheck/build scripts
- 單機遊玩設計：預設由 MVP 隨機電腦玩家補足其他席位
- 產品路線：單機桌遊 → 區域網路 → 線上模式

## 桌面版啟動方式

預期以 Wails 啟動桌面應用程式，開發模式使用：

```powershell
wails dev
```

目前 repository 尚未加入 Wails 專案設定與 Go binding，因此此指令尚不能在本專案執行。現有版本仍是分開啟動 Go API 與 Vite 前端：

```powershell
# Terminal 1：API
go run ./cmd/server
```

```powershell
# Terminal 2：前端
cd client
npm install
npm run dev
```

API 位於 `http://localhost:8080`，Vite 開發頁面位於 `http://localhost:5173`。

## 驗證

```powershell
go test ./...
cd client
npm run typecheck
npm run build
```

規則與後續階段請參閱 [Agent.md](./Agent.md)、[Agent 可執行規則](./docs/Agent-可執行規則.md) 與 [MVP 階段式需求](./docs/MVP-階段式需求.md)。
