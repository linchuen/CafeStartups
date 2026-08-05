import { useCallback, useEffect, useState } from 'react'
import { CardFace } from './CardFace'
import type { CardFaceData } from './CardFace'
import { arrangeDemandCards, DemandCard, DemandMarket } from './DemandMarket'
import { GameDashboard } from './modules/game/ui/GameDashboard'

const API = 'http://localhost:8080'

type Screen = 'home' | 'lobby' | 'game'
type Card = CardFaceData
type Player = { id: string; displayName: string; bot?: boolean; ready?: boolean; cash: number; loans: number; revenue?: number; score?: number; selectedKPIs?: string[]; brandAwareness?: number; products?: number; values?: number; resources?: number; handCount: number }
type CashFlowStatement = { period: number; beginningCash: number; operatingRevenue: number; otherIncome: number; operatingExpenses: number; interestPaid: number; principalRepayment: number; newLoans: number; endingCash: number }
type GameState = { id: string; status: string; seed: string; gameVersion: number; period: number; phase: string; round: number; demandBoard?: Record<string, number>; partnerOptions?: Card[]; starterShopOptions?: Card[]; players: Player[]; me?: { id: string; hand: Card[]; tableau: Card[]; discardCount: number; partner?: Card; starterShop?: Card; initialCardsSelected?: boolean; cash: number; loans: number; customers?: { kind: string; demand: string; unitPrice: number; count: number }[]; revenue?: number; score?: number; selectedKPIs?: string[]; kpiSelectionPeriod?: number; cashFlow?: CashFlowStatement[]; cashFlowRounds?: CashFlowStatement[]; brandAwareness?: number; products?: number; values?: number; resources?: number } }
type ApiError = Error & { code?: string }
const KPI_OPTIONS = [{ id: 'brand_awareness', label: '品牌知名度' }, { id: 'products', label: '特色產品' }, { id: 'values', label: '價值主張' }, { id: 'resources', label: '關鍵資源' }]
const REFERENCE_KPIS = [
  { label: '饕客滿意度', value: '每張 +5 分', tone: 'gourmet', symbol: '★' },
  { label: '盈餘', value: '每 30 → +1', tone: 'cash', symbol: '30' },
  { label: '通路', value: '每張 +4 分', tone: 'channel', symbol: '▰' },
  { label: '一般客滿意度', value: '每張 +4 分', tone: 'regular', symbol: '★' },
  { label: '知名度', value: '每 1 點 +1', tone: 'awareness', symbol: '✦' },
  { label: '品質', value: '每張 +3 分', tone: 'quality', symbol: '■' },
  { label: '全面顧客滿意度', value: '每張 +2 分', tone: 'total', symbol: '▣' },
  { label: '產品品項', value: '每張 +3 分', tone: 'product', symbol: '■' },
  { label: '成本', value: '每張 +3 分', tone: 'cost', symbol: '■' },
]
const REFERENCE_PERIODS = [
  { id: 1, label: '試營運', gourmet: '+10', regular: '+10', customers: [3, 2, 1, 1] },
  { id: 2, label: '正式營運', gourmet: '+20', regular: '+10', customers: [4, 3, 2, 1] },
  { id: 3, label: '擴大營運', gourmet: '+30', regular: '+10', customers: [5, 3, 2, 1] },
]

function ReferenceBoard(props: { period: number; phase?: string; round?: number; seed?: string }) {
  return <details className="reference-panel" open>
    <summary><span><span className="reference-panel-icon">?</span>玩家參考面板</span><small>示意內容 · 規則速查 · 可收合</small></summary>
    <div className="reference-board" aria-label="玩家規則參考面板">
      <div className="reference-kpis">
        <div className="reference-kpis-head"><span className="reference-bars">▮▮▮</span><strong>計分指標</strong><b>9</b></div>
        <div className="reference-kpi-grid">
          {REFERENCE_KPIS.map((item) => <div className={`reference-kpi reference-kpi-${item.tone}`} key={item.label}>
            <span className="reference-kpi-icon">{item.symbol}</span>
            <span className="reference-kpi-label">{item.label}</span>
            <small>{item.value}</small>
          </div>)}
        </div>
      </div>
      <div className="reference-rules">
        <div className="reference-brand-lockup"><span className="reference-cup">☕</span><span><strong>CAFÉ STARTUPS</strong><small>BREWING SUCCESSFUL ENTREPRENEURS</small></span></div>
        <div className="reference-period-head"><span className="reference-customer-spacer">客群需求與來客</span><span className={`reference-initial ${props.period === 0 ? 'is-current' : ''}`}><strong>0</strong><small>創業</small></span>{REFERENCE_PERIODS.map((item) => <span className={props.period === item.id ? 'is-current' : ''} key={item.id}><strong>{item.id}</strong><small>{item.label}</small></span>)}</div>
        <div className="reference-demand-and-cards"><div className="reference-demand-table">
        <ReferenceDemandRow type="gourmet" label="饕客" base="$10" additions={REFERENCE_PERIODS.map((item) => item.gourmet)} activePeriod={props.period} revealed={props.phase === 'learning' || props.phase === 'finished'} round={props.round} seed={props.seed} />
        <ReferenceDemandRow type="regular" label="一般客" base="$10" additions={REFERENCE_PERIODS.map((item) => item.regular)} activePeriod={props.period} revealed={props.phase === 'learning' || props.phase === 'finished'} round={props.round} seed={props.seed} />
        </div></div><div className="reference-bottom">
          <div className="reference-summary"><span className="reference-cube blue">◆</span><span>關鍵資源</span><strong>$0</strong><div className="reference-metric-order"><b>?</b><b>→</b><b>▣</b><b>→</b><b>■</b><b>→</b><b>◆</b></div></div>
          <div className="reference-ranking"><div className="reference-ranking-title"><span>市場排名</span><small>抽取市場袋顧客數</small></div><div className="reference-ranking-grid"><div className="reference-rank-labels"><span>1st</span><span>2nd</span><span>3rd</span><span>4th</span></div>{REFERENCE_PERIODS.map((item) => <div className={props.period === item.id ? 'reference-rank-column is-current' : 'reference-rank-column'} key={item.id}>{item.customers.map((count, index) => <span key={`${item.id}-${index}`}><b>{count}</b><i>●</i></span>)}</div>)}</div></div>
        </div>
      </div>
    </div>
  </details>
}

function ReferenceDemandRow(props: { type: 'gourmet' | 'regular'; label: string; base: string; additions: string[]; activePeriod: number; revealed?: boolean; round?: number; seed?: string }) {
  const questionMarks = (period: number) => props.type === 'gourmet' ? (period === 0 ? '?' : '??') : (period <= 2 ? '?' : '??')
  const numericSeed = [...(props.seed ?? 'local-game')].reduce((value, character) => ((value * 31) + character.charCodeAt(0)) >>> 0, 1)
  const cards = arrangeDemandCards(numericSeed)[props.type === 'gourmet' ? 'ordinary' : 'advanced']
  const cardAt = (column: number) => cards[column]
  const slot = (column: number, fallback: string) => {
    const card = cardAt(column)
    return card && props.revealed ? <DemandCard card={card} revealed /> : <><b>{questionMarks(column)}</b><small>{fallback}</small></>
  }
  return <div className={`reference-demand-row reference-demand-${props.type}`}>
    <div className="reference-customer-label"><span className="reference-customer-figures"><i /><i /></span><strong>{props.label}</strong><small>{props.type === 'gourmet' ? '饕客需求' : '一般客需求'}</small></div>
    <div className="reference-demand-slot reference-demand-base">{slot(0, props.base)}</div>
    {props.additions.map((addition, index) => <div className={`reference-demand-slot ${props.activePeriod === index + 1 ? 'is-current' : ''}`} key={`${props.type}-${index}`}>{slot(index + 1, addition)}</div>)}
  </div>
}

/* Legacy card-in-slot experiment kept out of the active board layout.
function ReferenceDemandRowWithCards(props: { type: 'gourmet' | 'regular'; label: string; base: string; additions: string[]; activePeriod: number; round?: number }) {
  const cardForPeriod = (period: number) => DEMAND_CARDS.find((card) => card.group === props.type && card.period === period)
  const revealed = (period: number) => period < props.activePeriod || period === props.activePeriod && period <= (props.round ?? 0)
  return <div className={`reference-demand-row reference-demand-${props.type}`}>
    <div className="reference-customer-label"><span className="reference-customer-figures"><i /><i /></span><strong>{props.label}</strong><small>{props.type === 'gourmet' ? '饕客需求' : '一般客需求'}</small></div>
    <div className="reference-demand-slot reference-demand-base"><b>?</b><small>{props.base}</small></div>
    {props.additions.map((addition, index) => { const period = index + 1; const card = cardForPeriod(period); return <div className={`reference-demand-slot reference-demand-card-slot ${props.activePeriod === period ? 'is-current' : ''}`} key={`${props.type}-${index}`}>{card ? <DemandCard card={card} revealed={revealed(period)} /> : <><b>?</b><small>{addition}</small></>}</div> })}
  </div>
}

function LegacyReferenceDemandRow(props: { type: 'gourmet' | 'regular'; label: string; base: string; additions: string[]; activePeriod: number; round?: number }) {
  const questionMarks = (period: number) => props.type === 'gourmet' ? (period === 0 ? '?' : '??') : (period <= 2 ? '?' : '??')
  return <div className={`reference-demand-row reference-demand-${props.type}`}>
    <div className="reference-customer-label"><span className="reference-customer-figures"><i /><i /></span><strong>{props.label}</strong><small>{props.type === 'gourmet' ? '高消費客群' : '一般消費客群'}</small></div>
    <div className="reference-demand-slot reference-demand-base"><b>{questionMarks(0)}</b><small>{props.base}</small></div>
    {props.additions.map((addition, index) => <div className={`reference-demand-slot ${props.activePeriod === index + 1 ? 'is-current' : ''}`} key={`${props.type}-${index}`}><b>{questionMarks(index + 1)}</b><small>✓ {addition}</small></div>)}
  </div>
}

}
*/

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API}${path}`, { headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) }, ...init })
  const body = await response.json().catch(() => ({}))
  if (!response.ok) { const error = new Error(body.code ?? '請求失敗') as ApiError; error.code = body.code; throw error }
  return body as T
}

export function App() {
  const [screen, setScreen] = useState<Screen>('home')
  const [name, setName] = useState('咖啡創業家')
  const [seed, setSeed] = useState('phase-3-demo')
  const [room, setRoom] = useState<GameState | null>(null)
  // Each visit to the home screen starts a fresh local game.
  const [token, setToken] = useState('')
  const [gameId, setGameId] = useState('')
  const [playerId, setPlayerId] = useState('')
  const [selectedCard, setSelectedCard] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [kpis, setKpis] = useState<[string, string]>(['brand_awareness', 'products'])

  const saveSession = (id: string, session: string, player: string) => {
    setGameId(id); setToken(session); setPlayerId(player)
  }

  const createRoom = async () => {
    setBusy(true); setError('')
    try {
      const created = await request<{ id: string; token: string; playerId: string; state: GameState }>('/api/games', { method: 'POST', body: JSON.stringify({ seed, displayName: name.trim() || '咖啡創業家' }) })
      saveSession(created.id, created.token, created.playerId)
      const started = await request<GameState>(`/api/games/${encodeURIComponent(created.id)}/start`, { method: 'POST', headers: { 'X-Session-Token': created.token } })
      setRoom(started); setScreen('game')
    } catch (cause) { setError(cause instanceof Error ? cause.message : '建立房間失敗') } finally { setBusy(false) }
  }

  const refresh = useCallback(async () => {
    if (!gameId || !token) return
    try {
      const next = await request<GameState>(`/api/games/${encodeURIComponent(gameId)}?token=${encodeURIComponent(token)}`)
      setRoom(next); if (next.status === 'playing' || next.status === 'finished') setScreen('game')
    } catch (cause) { if (cause instanceof Error) setError(cause.message) }
  }, [gameId, token])

  useEffect(() => { if (token && gameId) refresh() }, [refresh, token, gameId])
  useEffect(() => { if (!token || !gameId) return; const timer = window.setInterval(refresh, 1500); return () => window.clearInterval(timer) }, [refresh, token, gameId])

  const command = async (type: string, extra: Record<string, unknown> = {}) => {
    if (!room) return
    setBusy(true); setError('')
    try {
      const result = await request<{ state: GameState }>(`/api/games/${encodeURIComponent(room.id)}/commands`, { method: 'POST', body: JSON.stringify({ token, gameVersion: room.gameVersion, commandId: crypto.randomUUID(), type, ...extra }) })
      setRoom(result.state); if (result.state.phase === 'experiment') setScreen('game')
    } catch (cause) {
      const message = cause as ApiError
      setError(message.code === 'VERSION_CONFLICT' ? '狀態已更新，已重新同步。' : message.message)
      await refresh()
    } finally { setBusy(false) }
  }

  const setupInitialCards = async (partnerId: string, starterShopId: string) => {
    if (!room) return
    setBusy(true); setError('')
    try {
      const next = await request<GameState>(`/api/games/${encodeURIComponent(room.id)}/setup`, { method: 'POST', body: JSON.stringify({ token, partnerId, starterShopId }) })
      setRoom(next)
    } catch (cause) { setError(cause instanceof Error ? cause.message : '創業卡片選擇失敗') } finally { setBusy(false) }
  }

  const startGame = async () => {
    if (!room) return
    setBusy(true); setError('')
    try { const next = await request<GameState>(`/api/games/${encodeURIComponent(room.id)}/start`, { method: 'POST', headers: { 'X-Session-Token': token } }); setRoom(next); setScreen('game') }
    catch (cause) { setError(cause instanceof Error ? cause.message : '開始遊戲失敗') } finally { setBusy(false) }
  }

  const leave = () => { setRoom(null); setScreen('home'); setToken(''); setGameId(''); setPlayerId(''); localStorage.removeItem('cafe-session'); localStorage.removeItem('cafe-game-id'); localStorage.removeItem('cafe-player-id') }

  return <main className="app-shell">
    <header className="topbar"><span className="brand-mark">CS</span><span>Café Startups</span>{room && <span className="sync-pill">v{room.gameVersion} · 已同步</span>}</header>
    {screen === 'home' && <Home name={name} setName={setName} seed={seed} setSeed={setSeed} createRoom={createRoom} busy={busy} error={error} />}
    {screen === 'lobby' && room && <Lobby room={room} busy={busy} startGame={startGame} leave={leave} error={error} />}
    {screen === 'game' && room && <><ReferenceBoard period={room.period} round={room.round} seed={room.seed} /><GameDashboard room={room} selectedCard={selectedCard} setSelectedCard={setSelectedCard} command={command} busy={busy} error={error} leave={leave} setupInitialCards={setupInitialCards} /><CashFlowTable room={room} /></>}
  </main>
}

function Home(props: { name: string; setName: (v: string) => void; seed: string; setSeed: (v: string) => void; createRoom: () => void; busy: boolean; error: string }) {
  return <section className="landing layout-grid"><div className="intro"><p className="eyebrow">LOCAL BOARD GAME · MVP</p><h1>把假設<br /><em>煮成一杯</em><br />好生意。</h1><p className="lead">一位玩家在本機驗證咖啡館策略，電腦玩家只負責合法且可重播的隨機行動。</p><div className="stat-row"><span><strong>3</strong> 時期</span><span><strong>84</strong> 張 MVP 卡</span><span><strong>1</strong> 位玩家</span></div></div><div className="panel home-panel"><p className="eyebrow">SINGLE-PLAYER MVP</p><h2>開始單機桌遊</h2><label>你的名稱<input value={props.name} onChange={(e) => props.setName(e.target.value)} maxLength={20} /></label><button className="primary full" onClick={props.createRoom} disabled={props.busy}>{props.busy ? '建立中…' : '開始本機遊戲'}</button><p className="hint">遊戲會在本機自動補足隨機電腦玩家；目前不開放區網或線上加入。</p><details><summary>測試 seed（建立房間用）</summary><input value={props.seed} onChange={(e) => props.setSeed(e.target.value)} /></details>{props.error && <p className="error">{props.error}</p>}</div></section>
}

function Lobby(props: { room: GameState; busy: boolean; startGame: () => void; leave: () => void; error: string }) {
  return <section className="lobby-page layout-grid"><div className="panel room-card"><p className="eyebrow">LOCAL SETUP</p><h1>單機桌遊</h1><p className="muted">你是一位真人玩家；開始後由本機伺服器自動補足隨機電腦玩家。</p><div className="player-list">{props.room.players.map((player) => <div className="player-row" key={player.id}><span className="avatar">{player.displayName.slice(0, 1)}</span><span><strong>{player.displayName}</strong>{player.bot && <small>電腦玩家</small>}</span><span className="ready">本機</span></div>)}</div><div className="lobby-actions"><button className="secondary" onClick={props.leave}>離開</button><button className="accent" onClick={props.startGame} disabled={props.busy}>進入遊戲</button></div>{props.error && <p className="error">{props.error}</p>}</div><aside className="rules-card"><p className="eyebrow">SINGLE-PLAYER MVP</p><h2>先把玩法煮熟</h2><p>進入遊戲後，再選擇合夥人與初始店面；區域網路與線上模式會在玩法確認後再開發。</p><ol><li>進入遊戲面板</li><li>選擇合夥人與初始店面</li><li>完成三個營運時期</li></ol></aside></section>
}

function KpiPicker(props: { room: GameState; kpis: [string, string]; setKpis: (value: [string, string]) => void; command: (type: string, extra?: Record<string, unknown>) => void; busy: boolean }) {
  if (props.room.phase !== 'hypothesis') return null
  const saved = props.room.me?.kpiSelectionPeriod === props.room.period ? props.room.me.selectedKPIs : undefined
  const selected = saved?.length === 2 ? [saved[0], saved[1]] as [string, string] : props.kpis
  const choose = (index: 0 | 1, value: string) => {
    const next: [string, string] = [...selected] as [string, string]
    next[index] = value
    if (next[0] === next[1]) next[1] = KPI_OPTIONS.find((option) => option.id !== value)?.id ?? next[1]
    props.setKpis(next)
  }
  return <section className="panel kpi-selection-panel" aria-label="關鍵指標選擇"><div><p className="eyebrow">KEY METRICS</p><h2>{props.room.period === 2 ? '第 1 期結束，選擇你的關鍵指標' : '第 2 期結束，可重選一次關鍵指標'}</h2><p className="muted">選擇 2 個不同指標；所有玩家完成後才會開始下一期。</p></div><div className="kpi-selection-controls"><label>指標一<select value={selected[0]} onChange={(event) => choose(0, event.target.value)} disabled={props.busy}>{KPI_OPTIONS.map((option) => <option key={option.id} value={option.id}>{option.label}</option>)}</select></label><label>指標二<select value={selected[1]} onChange={(event) => choose(1, event.target.value)} disabled={props.busy}>{KPI_OPTIONS.map((option) => <option key={option.id} value={option.id}>{option.label}</option>)}</select></label><button className="accent" onClick={() => props.command('SET_KPI', { kpis: selected })} disabled={props.busy || selected[0] === selected[1]}>確認指標</button></div></section>
}

function CashFlowTable(props: { room: GameState }) {
  const statements = new Map((props.room.me?.cashFlow ?? []).map((statement) => [statement.period, statement]))
  const latestRounds = new Map((props.room.me?.cashFlowRounds ?? []).map((statement) => [statement.period, statement]))
  const periods = [1, 2, 3]
  const money = (value: number | undefined, signed = false) => value === undefined ? '待結算' : `${value < 0 ? '-$' : signed && value > 0 ? '+$' : '$'}${Math.abs(value)}萬`
  const statementFor = (period: number) => period === props.room.period ? latestRounds.get(period) ?? statements.get(period) : statements.get(period)
  const cell = (period: number, field: keyof CashFlowStatement, signed = false) => money(statementFor(period)?.[field] as number | undefined, signed)
  return <section className="cash-flow-panel panel" aria-label="現金流量表"><div className="cash-flow-heading"><div><p className="eyebrow">CASH FLOW STATEMENT</p><h2>現金流量表</h2><p>每回合結束更新目前期；第 6 回合完成後，再補上營收與利息結算。</p></div><span>單位：萬元</span></div><div className="cash-flow-grid"><div className="cash-flow-corner">項目</div>{periods.map((period) => <div className={`cash-flow-period ${props.room.period === period ? 'is-current' : ''}`} key={period}><strong>{period}</strong><small>{period === 1 ? '試營運' : period === 2 ? '正式營運' : '擴大營運'}</small></div>)}<div className="cash-flow-section">A. 期初現金</div>{periods.map((period) => <div className="cash-flow-value" key={`beginning-${period}`}>{cell(period, 'beginningCash')}</div>)}<div className="cash-flow-section">B. 營業收入</div>{periods.map((period) => <div className="cash-flow-value positive" key={`revenue-${period}`}>{cell(period, 'operatingRevenue', true)}</div>)}<div className="cash-flow-detail">B1. 顧客營收</div>{periods.map((period) => <div className="cash-flow-value detail" key={`customer-${period}`}>{cell(period, 'operatingRevenue', true)}</div>)}<div className="cash-flow-detail">B2. 棄牌補貼</div>{periods.map((period) => <div className="cash-flow-value detail" key={`other-${period}`}>{cell(period, 'otherIncome', true)}</div>)}<div className="cash-flow-section">C. 營運支出</div>{periods.map((period) => <div className="cash-flow-value negative" key={`expense-${period}`}>{cell(period, 'operatingExpenses', true)}</div>)}<div className="cash-flow-section">D. 籌資活動</div>{periods.map((period) => { const statement = statementFor(period); const financing = statement ? statement.newLoans - statement.interestPaid - statement.principalRepayment : undefined; return <div className="cash-flow-value" key={`financing-${period}`}>{money(financing, true)}</div> })}<div className="cash-flow-detail">D1. 支付利息</div>{periods.map((period) => <div className="cash-flow-value detail negative" key={`interest-${period}`}>{cell(period, 'interestPaid', true)}</div>)}<div className="cash-flow-detail">D2. 償還本金</div>{periods.map((period) => <div className="cash-flow-value detail negative" key={`principal-${period}`}>{cell(period, 'principalRepayment', true)}</div>)}<div className="cash-flow-detail">D3. 增貸資金</div>{periods.map((period) => <div className="cash-flow-value detail positive" key={`loans-${period}`}>{cell(period, 'newLoans', true)}</div>)}<div className="cash-flow-total">期末結餘（A+B+C+D）</div>{periods.map((period) => <div className="cash-flow-value cash-flow-total" key={`ending-${period}`}>{cell(period, 'endingCash')}</div>)}</div></section>
}

function PlayerCardsPanel(props: { room: GameState }) {
  const me = props.room.me
  const startingCards = [me?.partner, me?.starterShop].filter((card): card is Card => Boolean(card))
  const hand = me?.hand ?? []
  const tableau = me?.tableau ?? []
  const total = startingCards.length + hand.length + tableau.length

  return <section className="player-cards-panel panel" aria-label="我的卡牌">
    <div className="player-cards-heading">
      <div><p className="eyebrow">YOUR COLLECTION</p><h2>我的卡牌 <span>{total}</span></h2></div>
      <p>開局取得的卡牌與目前持有的手牌都會顯示在這裡。</p>
    </div>
    {total === 0 && <p className="empty player-cards-empty">尚未取得卡牌</p>}
    {startingCards.length > 0 && <div className="player-card-group"><div className="player-card-group-title"><strong>開局卡牌</strong><small>合夥人卡・創始店卡</small></div><div className="player-card-grid player-card-grid-starting">{startingCards.map((card) => <div className="player-card-tile" key={card.id}><CardFace card={card} /></div>)}</div></div>}
    {hand.length > 0 && <div className="player-card-group"><div className="player-card-group-title"><strong>目前手牌</strong><small>{hand.length} 張</small></div><div className="player-card-grid">{hand.map((card) => <div className="player-card-tile" key={card.id}><CardFace card={card} /></div>)}</div></div>}
    {tableau.length > 0 && <div className="player-card-group"><div className="player-card-group-title"><strong>已打出卡牌</strong><small>{tableau.length} 張</small></div><div className="player-card-grid">{tableau.map((card) => <div className="player-card-tile" key={card.id}><CardFace card={card} /></div>)}</div></div>}
  </section>
}

function GameTable(props: { room: GameState; selectedCard: string; setSelectedCard: (id: string) => void; command: (type: string, extra?: Record<string, unknown>) => void; busy: boolean; error: string; leave: () => void }) {
  const me = props.room.me
  const selected = me?.hand.find((card) => card.id === props.selectedCard)
  const finalRanking = [...props.room.players].sort((a, b) => (b.score ?? 0) - (a.score ?? 0))
  const phaseLabel: Record<string, string> = { hypothesis: '假設', experiment: '實驗', learning: '學習', finished: '結算' }
  return <section className="table-page"><div className="game-header"><div><p className="eyebrow">PERIOD {props.room.period} · ROUND {props.room.round}/6</p><h1>{phaseLabel[props.room.phase] ?? props.room.phase}</h1></div><div className="progress"><span className="period-dot active">1</span><i /><span className={props.room.period >= 2 ? 'period-dot active' : 'period-dot'}>2</span><i /><span className={props.room.period >= 3 ? 'period-dot active' : 'period-dot'}>3</span></div><button className="text-button" onClick={props.leave}>離開</button></div><ReferenceBoard period={props.room.period} phase={props.room.phase} /><div className="table-grid"><aside className="sidebar"><section className="panel resource-panel"><p className="eyebrow">你的咖啡館</p><div className="money">${me?.cash ?? 0}<small> 萬</small></div><div className="resource-line"><span>貸款</span><strong>{me?.loans ?? 0} / 6</strong></div><div className="resource-line"><span>營收</span><strong>${me?.revenue ?? 0} 萬</strong></div><div className="resource-line"><span>分數</span><strong>{me?.score ?? '—'}</strong></div><div className="resource-line"><span>品牌 / 產品</span><strong>{me?.brandAwareness ?? 0} / {me?.products ?? 0}</strong></div><div className="resource-line"><span>價值 / 資源</span><strong>{me?.values ?? 0} / {me?.resources ?? 0}</strong></div></section><section className="panel people-panel"><p className="eyebrow">玩家</p>{props.room.players.map((player) => <div className="mini-player" key={player.id}><span className="avatar small">{player.displayName.slice(0, 1)}</span><span>{player.displayName}{player.bot ? '（電腦）' : ''}{player.id === me?.id && '（你）'}</span><b>{player.handCount} 張</b></div>)}</section><button className="loan-button" onClick={() => props.command('TAKE_LOAN')} disabled={props.busy}>＋ 取得一筆貸款</button></aside><main className="board"><div className="board-banner"><div><span className="tag">目前階段</span><strong>{props.room.phase === 'experiment' ? '選擇你的實驗策略' : props.room.phase === 'learning' ? '結算本期市場與營收' : props.room.phase === 'finished' ? '遊戲完成' : '等待系統進入下一步'}</strong></div><span className="waiting-label">{props.room.players.filter((p) => !p.ready).length ? '玩家行動中' : '同步中'}</span></div>{props.room.phase === 'learning' && <div className="panel"><p>本期實驗完成，系統會自動處理顧客、營收與利息。</p><button className="accent" onClick={() => props.command('RESOLVE_LEARNING')} disabled={props.busy}>結算本期並繼續</button></div>}{props.room.phase === 'finished' && <div className="panel final-scoreboard"><p className="eyebrow">FINAL RESULTS</p><h2>三個時期完成，最終排名</h2><div className="scoreboard-header"><span>排名 / 玩家</span><span>分數</span><span>現金</span><span>營收</span><span>貸款</span></div>{finalRanking.map((player, index) => <div className={`scoreboard-row ${player.id === me?.id ? 'is-me' : ''}`} key={player.id}><strong>#{index + 1} {player.displayName}{player.bot ? '（電腦）' : '（你）'}</strong><b>{player.score ?? 0}</b><span>${player.cash} 萬</span><span>${player.revenue ?? 0} 萬</span><span>{player.loans}</span></div>)}<p className="hint">分數由最終現金與關鍵指標計算。</p></div>}<DemandMarket period={props.room.period} round={props.room.round} reveal={props.room.phase === 'learning' || props.room.phase === 'finished'} />{props.room.phase === 'experiment' && <><div className="hand-header"><div><p className="eyebrow">YOUR HAND · {me?.hand.length ?? 0} CARDS</p><h2>選擇一張牌</h2></div>{selected && <div className="selected-actions"><button className="secondary" onClick={() => props.command('DISCARD_SELECTED_CARD')} disabled={props.busy}>棄牌 +20</button><button className="primary" onClick={() => props.command('PLAY_SELECTED_CARD')} disabled={props.busy}>打出 ${selected.cost.cash}萬</button></div>}</div><div className="card-grid">{me?.hand.map((card) => <button key={card.id} className={`game-card ${selected?.id === card.id ? 'selected' : ''}`} onClick={() => { if (selected?.id === card.id) props.setSelectedCard(''); else { props.setSelectedCard(card.id); props.command('SELECT_CARD', { cardId: card.id }) } }}><CardFace card={card} selected={selected?.id === card.id} /></button>)}</div><button className="pass-button" onClick={() => props.command('PASS_HAND')} disabled={props.busy}>完成選擇並傳牌 <span>→</span></button></>}{props.error && <p className="error board-error">{props.error}</p>}</main></div></section>
}
