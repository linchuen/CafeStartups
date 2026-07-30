import { useState } from 'react'

type Game = { id: string; roomCode: string; status: string }

export function App() {
  const [screen, setScreen] = useState<'home' | 'lobby'>('home')
  const [game, setGame] = useState<Game | null>(null)
  const [error, setError] = useState('')

  async function createGame() {
    setError('')
    try {
      const response = await fetch('http://localhost:8080/api/games', { method: 'POST' })
      if (!response.ok) throw new Error('無法建立房間，請確認 Go server 已啟動。')
      setGame((await response.json()) as Game)
      setScreen('lobby')
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : '建立房間失敗。')
    }
  }

  return (
    <main className="app-shell">
      <section className="hero-card">
        <p className="eyebrow">LEAN STARTUP · DIGITAL BOARD GAME</p>
        <h1>Café Startups</h1>
        {screen === 'home' ? (
          <>
            <p className="lead">用一杯咖啡，驗證你的創業假設。</p>
            <div className="action-row">
              <button className="primary" onClick={createGame}>建立遊戲房間</button>
              <button className="secondary" onClick={() => setError('加入房間功能將在 Phase 2 實作。')}>加入房間</button>
            </div>
            {error && <p className="error">{error}</p>}
          </>
        ) : (
          <div className="lobby">
            <p className="lead">房間已建立，邀請夥伴加入吧。</p>
            <div className="room-code">{game?.roomCode}</div>
            <p className="muted">目前狀態：等待玩家（Phase 0 骨架）</p>
            <button className="secondary" onClick={() => setScreen('home')}>返回首頁</button>
          </div>
        )}
      </section>
      <footer>Phase 0 · 規格與技術骨架</footer>
    </main>
  )
}
