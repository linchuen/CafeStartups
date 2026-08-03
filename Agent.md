# Café Startups Web Agent Guide

## 1. 專案目標

使用 Go + React 將《Café Startups》製作成可在瀏覽器進行的桌遊 MVP。產品首先要能讓 1 位真人玩家即可開始一局三時期遊戲；其餘席位由簡易電腦玩家補足至 2–4 位。電腦玩家只需執行合法的隨機選擇，不需要策略推理或高階 AI，並清楚呈現每個玩家的手牌、咖啡館牌組、現金、貸款、顧客與分數。

規則來源：工作區根目錄的 `Cafe Startups 規則書_2025.pdf`；Agent 實作時必須同時閱讀 `docs/Agent-可執行規則.md`。若文件與規則書衝突，以規則書為準；若 MVP 為了可玩性採用簡化規則，必須在需求文件與 UI 中標示為「MVP 簡化」。

## 2. Agent 工作原則

1. 先讀 `Agent.md`、`docs/Agent-可執行規則.md` 與 `docs/MVP-階段式需求.md`，再修改程式。
2. 先確認目前階段的驗收條件，只實作該階段必要範圍。
3. 遊戲規則集中在 Go domain/service 層，不在 React 元件中自行計算。
4. 所有會改變遊戲狀態的操作都必須經過伺服器驗證；前端只是顯示與送出意圖。
5. 不以玩家端輸入的金額、成本、排名或分數作為可信資料；伺服器重新計算。
6. 牌堆洗牌、抽牌、抽顧客標記等隨機行為由伺服器負責，並保存足以重播/除錯的 seed 或事件紀錄。
7. 每次完成一個需求後，至少執行格式化、單元測試與建置檢查，並在交付訊息列出結果。
8. 不為了 MVP 預先加入登入、支付、社交、原生 App、策略型 AI 對手或完整工作坊功能；MVP 僅加入可讓單人遊玩的隨機電腦玩家。
9. 避免修改規則書與既有使用者檔案；新增程式與文件要保持小而可審查。

## 3. 技術基線

### Backend

- Go 1.22+。
- HTTP/JSON API；WebSocket 用於房間狀態推送與回合同步。
- 建議分層：`internal/domain`（規則與狀態）、`internal/application`（用例）、`internal/transport/http`、`internal/transport/ws`、`cmd/server`。
- MVP 優先使用記憶體儲存；若需要持久化，使用 SQLite，並以 repository 介面隔離。
- API 錯誤使用穩定的 error code，例如 `INVALID_ACTION`、`NOT_YOUR_TURN`、`INSUFFICIENT_CASH`。

### Frontend

- React + TypeScript + Vite。
- 狀態分為 server state 與 view state；遊戲真實狀態以伺服器推送為準。
- 建議頁面：建立/加入房間、遊戲桌、規則說明、結算頁。
- UI 先採 responsive desktop/tablet 版面；手機版只需不破壞主要操作。

### 開發品質

- Go：`gofmt`、`go vet`、`go test ./...`。
- Frontend：ESLint、TypeScript typecheck、測試與 production build。
- Domain 測試需覆蓋：階段轉移、傳牌、成本、貸款上限、收入、利息、排名、結算。
- 使用固定 seed 的測試 fixture，避免隨機造成不穩定測試。

## 4. 遊戲規則實作基線

- 遊戲有 3 個時期：第 1 期試營運、第 2 期正式營運、第 3 期擴大營運；每期有假設、實驗、學習三階段。
- 創業設定階段，玩家可在大廳選擇 1 張創辦人卡與 1 張創業店卡；目前卡面為 `mvp-fixture` 示意資料，選擇由後端保存。
- 玩家取得 150 萬元；正式創始店卡的成本與效果仍須依實體卡面校對後啟用。
- 每期發 7 張經營管理卡；實驗階段的回合計數從第 0 回合開始。第 0 回合代表已發牌但尚未完成任何選牌/傳牌/同步行動，之後依序完成第 1 至第 6 回合，最後 1 張進入中央覆蓋區。
- 第 1、3 時期順時鐘傳牌，第 2 時期逆時鐘傳牌。
- 單人模式預設補足電腦玩家至 4 個席位；電腦玩家的選牌、打牌/棄牌、傳牌與其他需要等待的行動由伺服器以固定 seed 的可重現隨機流程完成。
- 電腦玩家只從當下合法選項中隨機選擇，不評估策略、分數、對手或最佳解；其行為需標示為「MVP 隨機電腦玩家」。
- 行動可打出選定卡片或棄牌；棄牌在 MVP 先固定獲得 20 萬元。
- 貸款每張取得 50 萬元，每期利息 10 萬元，同時持有上限 6 張；不足支付利息時允許為付息增貸。
- 學習階段需處理市場變動、顧客需求、來客數、平均客單價/營收與利息/還款。
- 市佔率排名依序比較：品牌知名度、特色產品數、價值主張數、關鍵資源數、現金；同分時由伺服器以穩定順序決定。
- 基本客單價：饕客與一般客各 10 萬元；奧客為 0。需求滿足後依版面增額。
- 結算為現金得分加 2 個關鍵指標得分；同分比較現金，現金也相同則並列。

若完整卡牌資料尚未建檔，允許先使用「MVP 卡牌資料集」；資料格式與正式卡牌資料要相同，並標註 `source: mvp-fixture`，不可把 fixture 視為正式規則。

## 5. 建議資料邊界

主要 aggregate：`Game`、`Player`、`Card`、`DemandCard`、`MarketBag`、`Round`、`Score`。

每個 state mutation 建議產生事件：`GAME_CREATED`、`PLAYER_JOINED`、`PHASE_STARTED`、`CARD_SELECTED`、`CARDS_PASSED`、`CARD_PLAYED`、`CARD_DISCARDED`、`CUSTOMERS_DISTRIBUTED`、`REVENUE_SETTLED`、`GAME_FINISHED`。

前端不要直接依賴 Go struct 的內部欄位；以 versioned DTO 回傳 `gameVersion`，避免未來規則調整破壞 UI。

## 6. 完成定義（Definition of Done）

- 需求文件中的驗收條件可逐項驗證。
- 伺服器拒絕非法或非當前玩家操作，且前端顯示可理解的錯誤。
- 重整頁面後可從伺服器恢復房間狀態（MVP 可限於同一個執行個體）。
- 至少有一條從建房、加入、完成 3 時期到結算的 E2E happy path。
- README 或啟動說明包含 Go server、React client、測試與 seed/fixture 的執行方式。
- 文件中的「不在本階段」項目不得被誤宣稱已完成。

## 7. 變更流程

1. 在對應需求階段更新狀態與技術決策。
2. 先補 domain 測試，再實作規則。
3. 再接 API/WebSocket，最後接 UI。
4. 對規則有歧義時，在 PR/變更說明留下「規則書頁碼、採用解讀、未來調整方式」。
5. 若需求超出 MVP，先更新需求文件並取得產品方向確認，不直接擴張範圍。
