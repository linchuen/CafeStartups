import type { CSSProperties } from 'react'

export type DemandCardData = {
  id: string
  group: 'gourmet' | 'regular'
  period: number
  tier: 'ordinary' | 'advanced'
  icons: string[]
  valueBonus: number
}

const ICONS: Record<string, { label: string; symbol: string }> = {
  coffee: { label: '咖啡', symbol: '☕' },
  dessert: { label: '甜點', symbol: '◆' },
  beans: { label: '咖啡豆', symbol: '●' },
  taste: { label: '可口美味', symbol: '✦' },
  service: { label: '貼心服務', symbol: '♥' },
  price: { label: '合理價格', symbol: '$' },
}

// The first visual slice of the demand-card board. The final deck rules can
// replace this fixture without changing the card component or board layout.
const DEMAND_CARDS: DemandCardData[] = [
  { id: 'gourmet-1', group: 'gourmet', period: 1, tier: 'ordinary', icons: ['coffee'], valueBonus: 10 },
  { id: 'gourmet-2', group: 'gourmet', period: 1, tier: 'advanced', icons: ['coffee', 'taste'], valueBonus: 20 },
  { id: 'gourmet-3', group: 'gourmet', period: 2, tier: 'ordinary', icons: ['beans'], valueBonus: 10 },
  { id: 'gourmet-4', group: 'gourmet', period: 3, tier: 'advanced', icons: ['dessert', 'service'], valueBonus: 30 },
  { id: 'regular-1', group: 'regular', period: 1, tier: 'ordinary', icons: ['price'], valueBonus: 10 },
  { id: 'regular-2', group: 'regular', period: 1, tier: 'advanced', icons: ['coffee', 'service'], valueBonus: 10 },
  { id: 'regular-3', group: 'regular', period: 2, tier: 'ordinary', icons: ['dessert'], valueBonus: 10 },
  { id: 'regular-4', group: 'regular', period: 3, tier: 'advanced', icons: ['beans', 'price'], valueBonus: 10 },
]

function DemandCard({ card, revealed }: { card: DemandCardData; revealed: boolean }) {
  const accent = card.group === 'gourmet' ? '#e36d46' : '#d1aa32'
  const style = { '--demand-accent': accent } as CSSProperties
  if (!revealed) return <div className="demand-card demand-card-back" style={style} aria-label="未翻開的需求卡"><span>CS</span><small>需求卡</small></div>

  return <div className={`demand-card demand-card-face ${card.tier}`} style={style}>
    <div className="demand-card-heading"><span>{card.tier === 'advanced' ? '進階需求' : '普通需求'}</span><b>第 {card.period} 期</b></div>
    <div className="demand-card-icons">{card.icons.map((icon) => <span key={icon} title={ICONS[icon].label}>{ICONS[icon].symbol}</span>)}</div>
    <div className="demand-card-labels">{card.icons.map((icon) => <small key={icon}>{ICONS[icon].label}</small>)}</div>
    <div className="demand-card-bonus">客單價 <strong>+{card.valueBonus}</strong></div>
  </div>
}

export function DemandMarket({ period }: { period: number }) {
  const groups: Array<{ id: DemandCardData['group']; label: string }> = [{ id: 'gourmet', label: '饕客' }, { id: 'regular', label: '一般客' }]
  return <section className="demand-market panel" aria-label="顧客需求區">
    <div className="demand-market-heading"><div><p className="eyebrow">DEMAND CARDS</p><h2>顧客需求區</h2></div><span>普通 1 圖示 · 進階 2 圖示</span></div>
    <div className="demand-market-rows">{groups.map((group) => <div className={`demand-market-row ${group.id}`} key={group.id}>
      <div className="demand-group-label"><strong>{group.label}</strong><small>{group.id === 'gourmet' ? '高消費力客群' : '一般消費客群'}</small></div>
      {DEMAND_CARDS.filter((card) => card.group === group.id).map((card) => <DemandCard key={card.id} card={card} revealed={card.period <= period} />)}
    </div>)}</div>
  </section>
}
