import type { CSSProperties } from 'react'
import { CardCost } from './modules/game/ui/CardCost'

export type CardFaceData = {
  id: string
  name: string
  kind: string
  period: number
  description?: string
  effect?: string
  function?: string
  colorKey?: string
  startingCashBonus?: number
  starterShopId?: string
  demand?: Record<string, number>
  minPlayers?: number
  cost: { cash: number; icons: string[] }
  icons: string[]
  marketChange?: Record<string, number>
  customerCount?: Record<string, number>
}

const KIND_META: Record<string, { label: string; color: string; icon: string }> = {
  resource: { label: '關鍵資源', color: '#3976a6', icon: '◆' },
  product: { label: '特色產品', color: '#b98a25', icon: '●' },
  value: { label: '價值主張', color: '#b44f52', icon: '♥' },
  channel: { label: '門市通路', color: '#3f7d66', icon: '⌂' },
  marketing: { label: '行銷活動', color: '#7a5ba5', icon: '✦' },
  partner: { label: '合夥人', color: '#765341', icon: '♟' },
  starter_shop: { label: '創始店', color: '#3f7478', icon: '⌂' },
}

const MARKET_LABELS: Record<string, string> = { gourmet: '饕客', regular: '一般客', difficult: '奧客' }

function metaFor(kind: string) {
  return KIND_META[kind] ?? { label: kind, color: '#765341', icon: '●' }
}

function effectText(card: CardFaceData) {
  const market = Object.entries(card.marketChange ?? {})
    .filter(([, value]) => value !== 0)
    .map(([key, value]) => `${MARKET_LABELS[key] ?? key} ${value > 0 ? '+' : ''}${value}`)
  return market.length > 0 ? `市場變動：${market.join('、')}` : undefined
}

export function CardFace({ card, selected = false }: { card: CardFaceData; selected?: boolean }) {
  const meta = metaFor(card.kind)
  const style = { '--card-accent': meta.color } as CSSProperties
  const effect = effectText(card)

  return <>
    <div className="game-card-topline" style={style}>
      <span className="game-card-kind">{meta.icon} {meta.label}</span>
      <span className="game-card-period">第 {card.period} 期</span>
    </div>
    <div className="game-card-art" style={style}><span>{meta.icon}</span></div>
    <h3>{card.name}</h3>
    <p className="game-card-description">{card.description || effect || '打出此卡以發展咖啡館。'}</p>
    <div className="game-card-footer"><CardCost card={card} color={meta.color} /></div>
    {effect && card.description && <div className="game-card-market">{effect}</div>}
    {selected && <span className="game-card-selected">已選取</span>}
  </>
}
