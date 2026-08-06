# Café Startups 技術架構與開發規範

## 文件職責

- `Agent.md`：本專案的技術架構、程式分層與開發規範。
- `docs/Agent-可執行規則.md`：遊戲規則、狀態條件、公式與驗收基準的唯一來源。
- `docs/後端分層設計.md`：後端分層與責任邊界的補充說明。
- `docs/後端時序圖.md`：主要遊戲流程與服務互動時序。

實作遊戲功能時，先依 `docs/Agent-可執行規則.md` 判定規則，再依本文件決定程式放置位置。若程式、測試與規則文件不一致，應標記 `RULE_REVIEW_REQUIRED`，不得自行猜測規則。

## 專案目標

使用 Go、React 與 TypeScript 建立可在本機瀏覽器執行的單機版 Café Startups。

目前 MVP 範圍：

- 1 位真人玩家。
- 由固定 seed 的 Bot 補足至 2–4 位玩家。
- 本機 HTTP/JSON 通訊。
- 先完成核心遊戲流程、結算、卡牌與 UI，再擴充區域網路與線上模式。

## 技術堆疊

- Backend：Go 1.22+。
- Frontend：React、TypeScript、Vite、MUI。
- 通訊：本機 HTTP/JSON。
- 測試：Go unit test、server lifecycle test、TypeScript typecheck、production build。
- 遊戲資料：JSON fixture 與 schema；正式規則以規則文件為準。

## 後端分層

### `cmd/server`

應用程式入口，只負責組合依賴與啟動 HTTP server。不可在此放遊戲規則。

### `internal/server`

負責 HTTP 邊界：

- route 與 request parsing。
- session、room token、game version 的傳遞。
- application command 的 DTO 轉換。
- domain error 對應 HTTP status 與 error code。
- domain state 轉換為前端可讀的 view。

HTTP handler 不直接實作遊戲規則，也不直接修改 domain 欄位。

### `internal/application`

負責應用流程協調：

- 建立與管理遊戲 room。
- 將 command 分派給 domain。
- 協調期與 phase 的推進。
- 執行 Bot 行動。
- 組合 catalog 與遊戲初始化資料。

Application layer 不應重新實作 domain 的成本、顧客、結算或合法性判斷。

### `internal/domain`

唯一負責遊戲狀態與規則執行的層級，包括：

- Game、Player、Card 等核心狀態。
- phase、period、round 的合法轉移。
- 選牌、打牌、棄牌、傳牌。
- 卡片效果與成本判斷。
- 顧客分配、需求滿足、營收與現金流。
- 貸款、利息、還款與排名。
- KPI 與最終分數。

所有可由玩家行動觸發的狀態變更，都必須經過 domain method 的合法性檢查。

### `internal/catalog`

負責載入與轉換 JSON fixture，將資料交給 domain 驗證。不得在 loader 內加入遊戲流程或結算邏輯。

## 前端分層

### `client/src/modules/game/model`

放置 API view、卡牌資料與遊戲狀態的 TypeScript 型別。

### `client/src/modules/game/api`

負責 HTTP request、command、錯誤處理與 server state 更新。

### `client/src/modules/game/ui`

負責頁面與元件呈現：

- lobby：建立遊戲與初始設定。
- dashboard：遊戲狀態、玩家資源、手牌與行動。
- cards：卡牌面板、成本、圖示與卡片分類。
- market：需求市場與市場資訊。
- reference：規則參考板。

前端只呈現 server state，不自行決定遊戲結果。涉及成本、顧客、營收、分數或 phase 的判斷，必須以後端結果為準。

## 資料與卡牌規範

- 卡牌共用格式以 `data/card.schema.json` 為準。
- 需求卡使用 `kind: "demand"`，不另維護第二套卡片 schema。
- 卡片資料的 fixture 必須標示 `source`。
- 需要新增或變更卡片欄位時，先同步 schema、domain 型別、前端型別與測試。
- 卡片的資料描述與實際效果必須分開驗證；文字描述不能取代 domain 規則。

## 狀態與通訊規範

- server 是遊戲狀態的唯一來源。
- command 必須包含必要的 game version，避免舊狀態覆寫新狀態。
- command 失敗時不得部分套用狀態變更。
- domain 錯誤使用穩定 error code，例如 `INVALID_ACTION`、`INVALID_PHASE`、`CARD_NOT_FOUND`、`INSUFFICIENT_CASH`、`LOAN_LIMIT`。
- 前端遇到未知錯誤碼時，仍需保留可讀的錯誤訊息。
- 同一 command 的重送不得造成重複扣款、重複結算或重複推進回合。

## 開發流程

1. 先閱讀 `docs/Agent-可執行規則.md` 的相關章節。
2. 確認變更屬於 domain、application、server、catalog 或 frontend 哪一層。
3. 先補 domain 規則與測試，再接 application、API 與 UI。
4. 不在 UI 複製後端公式。
5. 若需求與現有規則衝突，停止猜測並標記 `RULE_REVIEW_REQUIRED`。
6. 完成後更新必要的測試與文件。

## 驗證要求

後端變更至少執行：

```text
gofmt
go test ./...
```

前端變更至少執行：

```text
npm run typecheck
npm run build
```

涉及遊戲流程的修改，必須涵蓋至少一條完整流程測試：建立遊戲、初始設定、進入實驗、完成選牌、學習結算與進入下一期。

涉及金額、卡牌成本、顧客或分數的修改，必須補上邊界條件測試，包括零值、資料不足、重複圖示、現金不足與非法 phase。

## 不可違反的架構原則

- 不將遊戲規則散落在 HTTP handler、React component 或 CSS。
- 不在前端自行扣現金、判定顧客滿意或計算最終分數。
- 不建立與 `docs/Agent-可執行規則.md` 平行且互相矛盾的規則來源。
- 不為了通過單一測試而繞過 domain 合法性檢查。
- 不把 MVP fixture 未驗證的內容宣稱為正式規則。
