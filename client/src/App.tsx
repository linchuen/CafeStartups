import { useCallback, useEffect, useState } from 'react'
import type { ApiError, GameState } from './modules/game/model/gameTypes'
import { GameDashboard } from './modules/game/ui/dashboard/GameDashboard'
import { CashFlowTable } from './modules/game/ui/dashboard/CashFlowTable'
import { Home } from './modules/game/ui/lobby/Home'
import { Lobby } from './modules/game/ui/lobby/Lobby'
import { ReferenceBoard } from './modules/game/ui/reference/ReferenceBoard'

const API = 'http://localhost:8080'

type Screen = 'home' | 'lobby' | 'game'
async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API}${path}`, { headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) }, ...init })
  const body = await response.json().catch(() => ({}))
  if (!response.ok) { const error = new Error(body.code ?? '請求失敗') as ApiError; error.code = body.code; throw error }
  return body as T
}

export function App() {
  const [screen, setScreen] = useState<Screen>('home')
  const [name, setName] = useState('咖啡創業家')
  // Leave the seed blank by default so the server assigns a fresh seed per game.
  // A fixed seed can still be entered manually to reproduce a game.
  const [seed, setSeed] = useState('')
  const [room, setRoom] = useState<GameState | null>(null)
  // Each visit to the home screen starts a fresh local game.
  const [token, setToken] = useState('')
  const [gameId, setGameId] = useState('')
  const [playerId, setPlayerId] = useState('')
  const [selectedCard, setSelectedCard] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

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
      window.alert(message.code ?? message.message)
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
    {screen === 'game' && room && <><ReferenceBoard period={room.period} phase={room.phase} round={room.round} seed={room.seed} demandCards={room.demandCards} /><GameDashboard room={room} selectedCard={selectedCard} setSelectedCard={setSelectedCard} command={command} busy={busy} error={error} leave={leave} setupInitialCards={setupInitialCards} /><CashFlowTable room={room} /></>}
  </main>
}
