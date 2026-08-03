import { useCallback, useEffect, useMemo, useState } from 'react'

const API = 'http://localhost:8080'

type Screen = 'home' | 'lobby' | 'game'
type Card = { id: string; name: string; kind: string; period: number; description?: string; cost: { cash: number; icons: string[] }; icons: string[] }
type Player = { id: string; displayName: string; bot?: boolean; ready: boolean; cash: number; loans: number; revenue?: number; score?: number; selectedKPIs?: string[]; brandAwareness?: number; products?: number; values?: number; resources?: number; handCount: number }
type GameState = { id: string; roomCode: string; status: string; seed: string; gameVersion: number; period: number; phase: string; round: number; demandBoard?: Record<string, number>; partnerOptions?: Card[]; starterShopOptions?: Card[]; players: Player[]; me?: { id: string; hand: Card[]; tableau: Card[]; discardCount: number; partner?: Card; starterShop?: Card; cash: number; loans: number; customers?: { kind: string; demand: string; unitPrice: number; count: number }[]; revenue?: number; score?: number; selectedKPIs?: string[]; brandAwareness?: number; products?: number; values?: number; resources?: number } }
type ApiError = Error & { code?: string }
const KPI_OPTIONS = [{ id: 'brand_awareness', label: '品牌知名度' }, { id: 'products', label: '特色產品' }, { id: 'values', label: '價值主張' }, { id: 'resources', label: '關鍵資源' }]
const REFERENCE_KPIS = [
  { label: '饕客滿意度', value: '每張 +5 分', tone: 'gourmet', symbol: '▣' },
  { label: '盈餘', value: '每 30 → +1', tone: 'cash', symbol: '30' },
  { label: '通路', value: '每張 +4 分', tone: 'channel', symbol: '▰' },
  { label: '一般客滿意度', value: '每張 +4 分', tone: 'regular', symbol: '▣' },
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

function ReferenceBoard(props: { period: number }) {
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
        <div className="reference-period-head"><span className="reference-customer-spacer">客群需求與來客</span><span className="reference-initial"><strong>0</strong><small>創業</small></span>{REFERENCE_PERIODS.map((item) => <span className={props.period === item.id ? 'is-current' : ''} key={item.id}><strong>{item.id}</strong><small>{item.label}</small></span>)}</div>
        <ReferenceDemandRow type="gourmet" label="饕客" base="$10" additions={REFERENCE_PERIODS.map((item) => item.gourmet)} activePeriod={props.period} />
        <ReferenceDemandRow type="regular" label="一般客" base="$10" additions={REFERENCE_PERIODS.map((item) => item.regular)} activePeriod={props.period} />
        <div className="reference-bottom">
          <div className="reference-summary"><span className="reference-cube blue">◆</span><span>關鍵資源</span><strong>$0</strong><div className="reference-metric-order"><b>?</b><b>→</b><b>▣</b><b>→</b><b>■</b><b>→</b><b>◆</b></div></div>
          <div className="reference-ranking"><div className="reference-ranking-title"><span>市場排名</span><small>抽取市場袋顧客數</small></div><div className="reference-ranking-grid"><div className="reference-rank-labels"><span>1st</span><span>2nd</span><span>3rd</span><span>4th</span></div>{REFERENCE_PERIODS.map((item) => <div className={props.period === item.id ? 'reference-rank-column is-current' : 'reference-rank-column'} key={item.id}>{item.customers.map((count, index) => <span key={`${item.id}-${index}`}><b>{count}</b><i>●</i></span>)}</div>)}</div></div>
        </div>
      </div>
    </div>
  </details>
}

function ReferenceDemandRow(props: { type: 'gourmet' | 'regular'; label: string; base: string; additions: string[]; activePeriod: number }) {
  const questionMarks = (period: number) => props.type === 'gourmet' ? (period === 0 ? '?' : '??') : (period <= 2 ? '?' : '??')
  return <div className={`reference-demand-row reference-demand-${props.type}`}>
    <div className="reference-customer-label"><span className="reference-customer-figures"><i /><i /></span><strong>{props.label}</strong><small>{props.type === 'gourmet' ? '高消費客群' : '一般消費客群'}</small></div>
    <div className="reference-demand-slot reference-demand-base"><b>{questionMarks(0)}</b><small>{props.base}</small></div>
    {props.additions.map((addition, index) => <div className={`reference-demand-slot ${props.activePeriod === index + 1 ? 'is-current' : ''}`} key={`${props.type}-${index}`}><b>{questionMarks(index + 1)}</b><small>✓ {addition}</small></div>)}
  </div>
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API}${path}`, { headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) }, ...init })
  const body = await response.json().catch(() => ({}))
  if (!response.ok) { const error = new Error(body.code ?? '請求失敗') as ApiError; error.code = body.code; throw error }
  return body as T
}

export function App() {
  const [screen, setScreen] = useState<Screen>('home')
  const [name, setName] = useState('咖啡創業家')
  const [roomInput, setRoomInput] = useState('')
  const [seed, setSeed] = useState('phase-3-demo')
  const [room, setRoom] = useState<GameState | null>(null)
  const [token, setToken] = useState(() => localStorage.getItem('cafe-session') ?? '')
  const [gameId, setGameId] = useState(() => localStorage.getItem('cafe-game-id') ?? '')
  const [playerId, setPlayerId] = useState(() => localStorage.getItem('cafe-player-id') ?? '')
  const [host, setHost] = useState(false)
  const [selectedCard, setSelectedCard] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [kpis, setKpis] = useState<[string, string]>(['brand_awareness', 'products'])

  const saveSession = (id: string, session: string, player: string) => {
    setGameId(id); setToken(session); setPlayerId(player)
    localStorage.setItem('cafe-game-id', id); localStorage.setItem('cafe-session', session); localStorage.setItem('cafe-player-id', player)
  }

  const joinRoom = async (key: string, asHost = false) => {
    setBusy(true); setError('')
    try {
      const joined = await request<{ token: string; playerId: string; state: GameState }>(`/api/games/${encodeURIComponent(key)}/join`, { method: 'POST', body: JSON.stringify({ displayName: name.trim() || '咖啡創業家' }) })
      saveSession(joined.state.id, joined.token, joined.playerId); setHost(asHost); setRoom(joined.state); setScreen('lobby')
    } catch (cause) { setError(cause instanceof Error ? cause.message : '加入房間失敗') } finally { setBusy(false) }
  }

  const createRoom = async () => {
    setBusy(true); setError('')
    try {
      const created = await request<{ id: string }>('/api/games', { method: 'POST', body: JSON.stringify({ seed }) })
      await joinRoom(created.id, true)
    } catch (cause) { setError(cause instanceof Error ? cause.message : '建立房間失敗'); setBusy(false) }
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

  const ownPlayer = useMemo(() => room?.players.find((player) => player.id === playerId), [room, playerId])
  const allReady = Boolean(room?.players.length && room.players.every((player) => player.ready))

  const toggleReady = async () => {
    if (!room) return
    setBusy(true); setError('')
    try { const next = await request<GameState>(`/api/games/${encodeURIComponent(room.id)}/ready`, { method: 'POST', body: JSON.stringify({ token, ready: !ownPlayer?.ready }) }); setRoom(next) }
    catch (cause) { setError(cause instanceof Error ? cause.message : '準備狀態更新失敗') } finally { setBusy(false) }
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
    try { const next = await request<GameState>(`/api/games/${encodeURIComponent(room.id)}/start`, { method: 'POST', headers: { 'X-Session-Token': token }, body: JSON.stringify({ kpis }) }); setRoom(next); setScreen('game') }
    catch (cause) { setError(cause instanceof Error ? cause.message : '開始遊戲失敗') } finally { setBusy(false) }
  }

  const leave = () => { setRoom(null); setScreen('home'); setToken(''); setGameId(''); setPlayerId(''); localStorage.removeItem('cafe-session'); localStorage.removeItem('cafe-game-id'); localStorage.removeItem('cafe-player-id') }

  return <main className="app-shell">
    <header className="topbar"><span className="brand-mark">CS</span><span>Café Startups</span>{room && <span className="sync-pill">v{room.gameVersion} · 已同步</span>}</header>
    {screen === 'home' && <Home name={name} setName={setName} roomInput={roomInput} setRoomInput={setRoomInput} seed={seed} setSeed={setSeed} createRoom={createRoom} joinRoom={() => joinRoom(roomInput)} busy={busy} error={error} />}
    {screen === 'lobby' && room && <Lobby room={room} host={host} busy={busy} allReady={allReady} kpis={kpis} setKpis={setKpis} toggleReady={toggleReady} setupInitialCards={setupInitialCards} startGame={startGame} leave={leave} error={error} />}
    {screen === 'game' && room && <GameTable room={room} selectedCard={selectedCard} setSelectedCard={setSelectedCard} command={command} busy={busy} error={error} leave={leave} />}
  </main>
}

function Home(props: { name: string; setName: (v: string) => void; roomInput: string; setRoomInput: (v: string) => void; seed: string; setSeed: (v: string) => void; createRoom: () => void; joinRoom: () => void; busy: boolean; error: string }) {
  return <section className="landing layout-grid"><div className="intro"><p className="eyebrow">LEAN STARTUP · DIGITAL BOARD GAME</p><h1>把假設<br /><em>煮成一杯</em><br />好生意。</h1><p className="lead">選擇你的策略、傳遞你的牌，從一間小咖啡店開始驗證商業模式。</p><div className="stat-row"><span><strong>3</strong> 時期</span><span><strong>84</strong> 張 MVP 卡</span><span><strong>1–4</strong> 真人</span></div></div><div className="panel home-panel"><h2>開始一局</h2><label>你的名稱<input value={props.name} onChange={(e) => props.setName(e.target.value)} maxLength={20} /></label><button className="primary full" onClick={props.createRoom} disabled={props.busy}>{props.busy ? '建立中…' : '建立新房間'}</button><p className="hint">單人開始時，系統會自動加入隨機電腦玩家。</p><div className="divider"><span>或</span></div><label>房間代碼<input value={props.roomInput} onChange={(e) => props.setRoomInput(e.target.value.toUpperCase())} placeholder="例如 A1B2C3" maxLength={32} /></label><button className="secondary full" onClick={props.joinRoom} disabled={props.busy || !props.roomInput.trim()}>加入房間</button><details><summary>測試 seed（建立房間用）</summary><input value={props.seed} onChange={(e) => props.setSeed(e.target.value)} /></details>{props.error && <p className="error">{props.error}</p>}</div></section>
}

function InitialCardPicker(props: { room: GameState; busy: boolean; setupInitialCards: (partnerId: string, starterShopId: string) => void }) {
  const selectedPartner = props.room.me?.partner?.id ?? props.room.partnerOptions?.[0]?.id ?? ''
  const selectedShop = props.room.me?.starterShop?.id ?? props.room.starterShopOptions?.[0]?.id ?? ''
  return <section className="initial-card-picker" aria-label="創業卡片選擇">
    <div className="initial-card-heading"><div><p className="eyebrow">START YOUR CAFÉ</p><h2>選擇創辦人與創業店</h2></div><span>示意卡片 · MVP</span></div>
    <p className="initial-card-note">創業前先選擇 1 張創辦人卡與 1 張創業店卡；選擇會儲存在後端房間狀態。</p>
    <div className="initial-card-group"><strong>創辦人卡</strong><div className="initial-card-grid">{(props.room.partnerOptions ?? []).map((card) => <button type="button" className={`initial-card ${selectedPartner === card.id ? 'selected' : ''}`} key={card.id} onClick={() => props.setupInitialCards(card.id, selectedShop)} disabled={props.busy}><span className="initial-card-symbol">♟</span><span><b>{card.name}</b><small>{card.description}</small></span>{selectedPartner === card.id && <em>已選</em>}</button>)}</div></div>
    <div className="initial-card-group"><strong>創業店卡</strong><div className="initial-card-grid">{(props.room.starterShopOptions ?? []).map((card) => <button type="button" className={`initial-card shop-card ${selectedShop === card.id ? 'selected' : ''}`} key={card.id} onClick={() => props.setupInitialCards(selectedPartner, card.id)} disabled={props.busy}><span className="initial-card-symbol">⌂</span><span><b>{card.name}</b><small>{card.description}</small></span>{selectedShop === card.id && <em>已選</em>}</button>)}</div></div>
  </section>
}

function Lobby(props: { room: GameState; host: boolean; busy: boolean; allReady: boolean; kpis: [string, string]; setKpis: (value: [string, string]) => void; toggleReady: () => void; setupInitialCards: (partnerId: string, starterShopId: string) => void; startGame: () => void; leave: () => void; error: string }) {
  return <section className="lobby-page layout-grid"><div className="panel room-card"><p className="eyebrow">WAITING ROOM</p><h1>{props.room.roomCode}</h1><p className="muted">單人可直接開始；開始後不足席位會由隨機電腦玩家補足。</p><div className="player-list">{props.room.players.map((player, index) => <div className="player-row" key={player.id}><span className="avatar">{player.displayName.slice(0, 1)}</span><span><strong>{player.displayName}</strong>{player.bot && <small>電腦</small>}{index === 0 && <small>房主</small>}</span><span className={player.ready ? 'ready' : 'waiting'}>{player.ready ? '已準備' : '等待中'}</span></div>)}{Array.from({ length: 4 - props.room.players.length }).map((_, i) => <div className="player-row empty" key={i}><span className="avatar">·</span><span>等待玩家加入</span></div>)}</div><InitialCardPicker room={props.room} busy={props.busy} setupInitialCards={props.setupInitialCards} />{props.host && <div className="kpi-picker"><p className="eyebrow">YOUR KEY METRICS</p><label>關鍵指標一<select value={props.kpis[0]} onChange={(e) => { const value = e.target.value; props.setKpis([value, props.kpis[1] === value ? props.kpis[0] : props.kpis[1]]) }}>{KPI_OPTIONS.map((option) => <option key={option.id} value={option.id}>{option.label}</option>)}</select></label><label>關鍵指標二<select value={props.kpis[1]} onChange={(e) => { const value = e.target.value; props.setKpis([props.kpis[0] === value ? props.kpis[1] : props.kpis[0], value]) }}>{KPI_OPTIONS.map((option) => <option key={option.id} value={option.id}>{option.label}</option>)}</select></label></div>}<div className="lobby-actions"><button className="secondary" onClick={props.leave}>離開</button><button className="primary" onClick={props.toggleReady} disabled={props.busy}>{props.room.players.find((p) => p.id === props.room.me?.id)?.ready ? '取消準備' : '準備完成'}</button>{props.host && <button className="accent" onClick={props.startGame} disabled={props.busy || !props.allReady}>開始遊戲</button>}</div>{props.host && !props.allReady && <p className="hint">準備完成後即可開始，其他席位會自動加入電腦玩家。</p>}{props.error && <p className="error">{props.error}</p>}</div><aside className="rules-card"><p className="eyebrow">HOW TO PLAY</p><h2>一場實驗，六次選擇</h2><p>每個時期會拿到 7 張管理牌。選牌、傳牌，最後把你的策略打進咖啡館。</p><ol><li>設定你的關鍵指標</li><li>每回合選一張牌並傳遞</li><li>打出策略或棄牌換取現金</li></ol></aside></section>
}

function GameTable(props: { room: GameState; selectedCard: string; setSelectedCard: (id: string) => void; command: (type: string, extra?: Record<string, unknown>) => void; busy: boolean; error: string; leave: () => void }) {
  const me = props.room.me
  const selected = me?.hand.find((card) => card.id === props.selectedCard)
  const finalRanking = [...props.room.players].sort((a, b) => (b.score ?? 0) - (a.score ?? 0))
  const phaseLabel: Record<string, string> = { hypothesis: '假設', experiment: '實驗', learning: '學習', finished: '結算' }
  return <section className="table-page"><div className="game-header"><div><p className="eyebrow">PERIOD {props.room.period} · ROUND {props.room.round}/6</p><h1>{phaseLabel[props.room.phase] ?? props.room.phase}</h1></div><div className="progress"><span className="period-dot active">1</span><i /><span className={props.room.period >= 2 ? 'period-dot active' : 'period-dot'}>2</span><i /><span className={props.room.period >= 3 ? 'period-dot active' : 'period-dot'}>3</span></div><button className="text-button" onClick={props.leave}>離開</button></div><ReferenceBoard period={props.room.period} /><div className="table-grid"><aside className="sidebar"><section className="panel resource-panel"><p className="eyebrow">你的咖啡館</p><div className="money">${me?.cash ?? 0}<small> 萬</small></div><div className="resource-line"><span>貸款</span><strong>{me?.loans ?? 0} / 6</strong></div><div className="resource-line"><span>營收</span><strong>${me?.revenue ?? 0} 萬</strong></div><div className="resource-line"><span>分數</span><strong>{me?.score ?? '—'}</strong></div><div className="resource-line"><span>品牌 / 產品</span><strong>{me?.brandAwareness ?? 0} / {me?.products ?? 0}</strong></div><div className="resource-line"><span>價值 / 資源</span><strong>{me?.values ?? 0} / {me?.resources ?? 0}</strong></div></section><section className="panel people-panel"><p className="eyebrow">玩家</p>{props.room.players.map((player) => <div className="mini-player" key={player.id}><span className="avatar small">{player.displayName.slice(0, 1)}</span><span>{player.displayName}{player.bot ? '（電腦）' : ''}{player.id === me?.id && '（你）'}</span><b>{player.handCount} 張</b></div>)}</section><button className="loan-button" onClick={() => props.command('TAKE_LOAN')} disabled={props.busy}>＋ 取得一筆貸款</button></aside><main className="board"><div className="board-banner"><div><span className="tag">目前階段</span><strong>{props.room.phase === 'experiment' ? '選擇你的實驗策略' : props.room.phase === 'learning' ? '結算本期市場與營收' : props.room.phase === 'finished' ? '遊戲完成' : '等待系統進入下一步'}</strong></div><span className="waiting-label">{props.room.players.filter((p) => !p.ready).length ? '玩家行動中' : '同步中'}</span></div>{props.room.phase === 'learning' && <div className="panel"><p>本期實驗完成，系統會自動處理顧客、營收與利息。</p><button className="accent" onClick={() => props.command('RESOLVE_LEARNING')} disabled={props.busy}>結算本期並繼續</button></div>}{props.room.phase === 'finished' && <div className="panel final-scoreboard"><p className="eyebrow">FINAL RESULTS</p><h2>三個時期完成，最終排名</h2><div className="scoreboard-header"><span>排名 / 玩家</span><span>分數</span><span>現金</span><span>營收</span><span>貸款</span></div>{finalRanking.map((player, index) => <div className={`scoreboard-row ${player.id === me?.id ? 'is-me' : ''}`} key={player.id}><strong>#{index + 1} {player.displayName}{player.bot ? '（電腦）' : '（你）'}</strong><b>{player.score ?? 0}</b><span>${player.cash} 萬</span><span>${player.revenue ?? 0} 萬</span><span>{player.loans}</span></div>)}<p className="hint">分數由最終現金與關鍵指標計算。</p></div>}<div className="market panel"><div><p className="eyebrow">MARKET BOARD</p><h2>需求市場</h2></div><div className="market-cards"><span>饕客 <b>{props.room.demandBoard?.gourmet ?? 0}</b></span><span>一般客 <b>{props.room.demandBoard?.regular ?? 0}</b></span><span>奧客 <b>{props.room.demandBoard?.difficult ?? 0}</b></span></div></div>{props.room.phase === 'experiment' && <><div className="hand-header"><div><p className="eyebrow">YOUR HAND · {me?.hand.length ?? 0} CARDS</p><h2>選擇一張牌</h2></div>{selected && <div className="selected-actions"><button className="secondary" onClick={() => props.command('DISCARD_SELECTED_CARD')} disabled={props.busy}>棄牌 +20</button><button className="primary" onClick={() => props.command('PLAY_SELECTED_CARD')} disabled={props.busy}>打出 ${selected.cost.cash}萬</button></div>}</div><div className="card-grid">{me?.hand.map((card) => <button key={card.id} className={`game-card ${selected?.id === card.id ? 'selected' : ''}`} onClick={() => { if (selected?.id === card.id) props.setSelectedCard(''); else { props.setSelectedCard(card.id); props.command('SELECT_CARD', { cardId: card.id }) } }}><span className="card-kind">{card.kind}</span><h3>{card.name}</h3><span className="card-cost">${card.cost.cash}萬</span><div className="card-icons">{card.icons.map((icon) => <span key={icon}>{icon}</span>)}</div></button>)}</div><button className="pass-button" onClick={() => props.command('PASS_HAND')} disabled={props.busy}>完成選擇並傳牌 <span>→</span></button></>}{props.error && <p className="error board-error">{props.error}</p>}</main></div></section>
}
