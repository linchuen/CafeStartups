import { useCallback, useEffect, useMemo, useState } from 'react'

const API = 'http://localhost:8080'

type Screen = 'home' | 'lobby' | 'game'
type Card = { id: string; name: string; kind: string; period: number; cost: { cash: number; icons: string[] }; icons: string[] }
type Player = { id: string; displayName: string; bot?: boolean; ready: boolean; cash: number; loans: number; handCount: number }
type GameState = { id: string; roomCode: string; status: string; seed: string; gameVersion: number; period: number; phase: string; round: number; players: Player[]; me?: { id: string; hand: Card[]; tableau: Card[]; discardCount: number; cash: number; loans: number } }
type ApiError = Error & { code?: string }

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
      setRoom(next); if (next.status === 'playing') setScreen('game')
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

  const startGame = async () => {
    if (!room) return
    setBusy(true); setError('')
    try { const next = await request<GameState>(`/api/games/${encodeURIComponent(room.id)}/start`, { method: 'POST', headers: { 'X-Session-Token': token } }); setRoom(next); setScreen('game') }
    catch (cause) { setError(cause instanceof Error ? cause.message : '開始遊戲失敗') } finally { setBusy(false) }
  }

  const leave = () => { setRoom(null); setScreen('home'); setToken(''); setGameId(''); setPlayerId(''); localStorage.removeItem('cafe-session'); localStorage.removeItem('cafe-game-id'); localStorage.removeItem('cafe-player-id') }

  return <main className="app-shell">
    <header className="topbar"><span className="brand-mark">CS</span><span>Café Startups</span>{room && <span className="sync-pill">v{room.gameVersion} · 已同步</span>}</header>
    {screen === 'home' && <Home name={name} setName={setName} roomInput={roomInput} setRoomInput={setRoomInput} seed={seed} setSeed={setSeed} createRoom={createRoom} joinRoom={() => joinRoom(roomInput)} busy={busy} error={error} />}
    {screen === 'lobby' && room && <Lobby room={room} host={host} busy={busy} allReady={allReady} toggleReady={toggleReady} startGame={startGame} leave={leave} error={error} />}
    {screen === 'game' && room && <GameTable room={room} selectedCard={selectedCard} setSelectedCard={setSelectedCard} command={command} busy={busy} error={error} leave={leave} />}
  </main>
}

function Home(props: { name: string; setName: (v: string) => void; roomInput: string; setRoomInput: (v: string) => void; seed: string; setSeed: (v: string) => void; createRoom: () => void; joinRoom: () => void; busy: boolean; error: string }) {
  return <section className="landing layout-grid"><div className="intro"><p className="eyebrow">LEAN STARTUP · DIGITAL BOARD GAME</p><h1>把假設<br /><em>煮成一杯</em><br />好生意。</h1><p className="lead">選擇你的策略、傳遞你的牌，從一間小咖啡店開始驗證商業模式。</p><div className="stat-row"><span><strong>3</strong> 時期</span><span><strong>84</strong> 張 MVP 卡</span><span><strong>1–4</strong> 真人</span></div></div><div className="panel home-panel"><h2>開始一局</h2><label>你的名稱<input value={props.name} onChange={(e) => props.setName(e.target.value)} maxLength={20} /></label><button className="primary full" onClick={props.createRoom} disabled={props.busy}>{props.busy ? '建立中…' : '建立新房間'}</button><p className="hint">單人開始時，系統會自動加入隨機電腦玩家。</p><div className="divider"><span>或</span></div><label>房間代碼<input value={props.roomInput} onChange={(e) => props.setRoomInput(e.target.value.toUpperCase())} placeholder="例如 A1B2C3" maxLength={32} /></label><button className="secondary full" onClick={props.joinRoom} disabled={props.busy || !props.roomInput.trim()}>加入房間</button><details><summary>測試 seed（建立房間用）</summary><input value={props.seed} onChange={(e) => props.setSeed(e.target.value)} /></details>{props.error && <p className="error">{props.error}</p>}</div></section>
}

function Lobby(props: { room: GameState; host: boolean; busy: boolean; allReady: boolean; toggleReady: () => void; startGame: () => void; leave: () => void; error: string }) {
  return <section className="lobby-page layout-grid"><div className="panel room-card"><p className="eyebrow">WAITING ROOM</p><h1>{props.room.roomCode}</h1><p className="muted">單人可直接開始；開始後不足席位會由隨機電腦玩家補足。</p><div className="player-list">{props.room.players.map((player, index) => <div className="player-row" key={player.id}><span className="avatar">{player.displayName.slice(0, 1)}</span><span><strong>{player.displayName}</strong>{player.bot && <small>電腦</small>}{index === 0 && <small>房主</small>}</span><span className={player.ready ? 'ready' : 'waiting'}>{player.ready ? '已準備' : '等待中'}</span></div>)}{Array.from({ length: 4 - props.room.players.length }).map((_, i) => <div className="player-row empty" key={i}><span className="avatar">·</span><span>等待玩家加入</span></div>)}</div><div className="lobby-actions"><button className="secondary" onClick={props.leave}>離開</button><button className="primary" onClick={props.toggleReady} disabled={props.busy}>{props.room.players.find((p) => p.id === props.room.me?.id)?.ready ? '取消準備' : '準備完成'}</button>{props.host && <button className="accent" onClick={props.startGame} disabled={props.busy || !props.allReady}>開始遊戲</button>}</div>{props.host && !props.allReady && <p className="hint">準備完成後即可開始，其他席位會自動加入電腦玩家。</p>}{props.error && <p className="error">{props.error}</p>}</div><aside className="rules-card"><p className="eyebrow">HOW TO PLAY</p><h2>一場實驗，六次選擇</h2><p>每個時期會拿到 7 張管理牌。選牌、傳牌，最後把你的策略打進咖啡館。</p><ol><li>設定你的關鍵指標</li><li>每回合選一張牌並傳遞</li><li>打出策略或棄牌換取現金</li></ol></aside></section>
}

function GameTable(props: { room: GameState; selectedCard: string; setSelectedCard: (id: string) => void; command: (type: string, extra?: Record<string, unknown>) => void; busy: boolean; error: string; leave: () => void }) {
  const me = props.room.me
  const selected = me?.hand.find((card) => card.id === props.selectedCard)
  const phaseLabel: Record<string, string> = { hypothesis: '假設', experiment: '實驗', learning: '學習', finished: '結算' }
  return <section className="table-page"><div className="game-header"><div><p className="eyebrow">PERIOD {props.room.period} · ROUND {props.room.round}/6</p><h1>{phaseLabel[props.room.phase] ?? props.room.phase}</h1></div><div className="progress"><span className="period-dot active">1</span><i /><span className={props.room.period >= 2 ? 'period-dot active' : 'period-dot'}>2</span><i /><span className={props.room.period >= 3 ? 'period-dot active' : 'period-dot'}>3</span></div><button className="text-button" onClick={props.leave}>離開</button></div><div className="table-grid"><aside className="sidebar"><section className="panel resource-panel"><p className="eyebrow">你的咖啡館</p><div className="money">${me?.cash ?? 0}<small> 萬</small></div><div className="resource-line"><span>貸款</span><strong>{me?.loans ?? 0} / 6</strong></div><div className="resource-line"><span>顧客</span><strong>—</strong></div></section><section className="panel people-panel"><p className="eyebrow">玩家</p>{props.room.players.map((player) => <div className="mini-player" key={player.id}><span className="avatar small">{player.displayName.slice(0, 1)}</span><span>{player.displayName}{player.bot ? '（電腦）' : ''}{player.id === me?.id && '（你）'}</span><b>{player.handCount} 張</b></div>)}</section><button className="loan-button" onClick={() => props.command('TAKE_LOAN')} disabled={props.busy}>＋ 取得一筆貸款</button></aside><main className="board"><div className="board-banner"><div><span className="tag">目前階段</span><strong>{props.room.phase === 'experiment' ? '選擇你的實驗策略' : '等待系統進入下一步'}</strong></div><span className="waiting-label">{props.room.players.filter((p) => !p.ready).length ? '玩家行動中' : '同步中'}</span></div><div className="market panel"><div><p className="eyebrow">MARKET BOARD</p><h2>需求市場</h2></div><div className="market-cards"><span>饕客 <b>10</b></span><span>一般客 <b>10</b></span><span>奧客 <b>0</b></span></div></div><div className="hand-header"><div><p className="eyebrow">YOUR HAND · {me?.hand.length ?? 0} CARDS</p><h2>選擇一張牌</h2></div>{selected && <div className="selected-actions"><button className="secondary" onClick={() => props.command('DISCARD_SELECTED_CARD')} disabled={props.busy}>棄牌 +20</button><button className="primary" onClick={() => props.command('PLAY_SELECTED_CARD')} disabled={props.busy}>打出 ${selected.cost.cash}萬</button></div>}</div><div className="card-grid">{me?.hand.map((card) => <button key={card.id} className={`game-card ${selected?.id === card.id ? 'selected' : ''}`} onClick={() => { if (selected?.id === card.id) props.setSelectedCard(''); else { props.setSelectedCard(card.id); props.command('SELECT_CARD', { cardId: card.id }) } }}><span className="card-kind">{card.kind}</span><h3>{card.name}</h3><span className="card-cost">${card.cost.cash}萬</span><div className="card-icons">{card.icons.map((icon) => <span key={icon}>{icon}</span>)}</div></button>)}</div><button className="pass-button" onClick={() => props.command('PASS_HAND')} disabled={props.busy}>完成選擇並傳牌 <span>→</span></button>{props.error && <p className="error board-error">{props.error}</p>}</main></div></section>
}
