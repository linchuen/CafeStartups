import type { CSSProperties } from 'react'

export type DemandCardData = { id: string; row: 'ordinary' | 'advanced'; quantity: 1 | 2; icons: string[]; source?: string }

const DEMAND_TYPES = [
  { id: 'coffee', label: '\u5496\u5561', symbol: '\u2615' },
  { id: 'dessert', label: '\u751c\u9ede', symbol: '\u25c7' },
  { id: 'beans', label: '\u5496\u5561\u8c46', symbol: '\u25c9' },
  { id: 'taste', label: '\u53ef\u53e3\u7f8e\u5473', symbol: '\u2665' },
  { id: 'service', label: '\u8cbc\u5fc3\u670d\u52d9', symbol: '\u2665' },
  { id: 'price', label: '\u5408\u7406\u50f9\u683c', symbol: '$' },
] as const

export const DEMAND_CARDS: DemandCardData[] = [
  ...DEMAND_TYPES.map((type, index) => ({ id: `ordinary-${String(index + 1).padStart(2, '0')}`, row: 'ordinary' as const, quantity: 1 as const, icons: [type.id], source: 'mvp-fixture' })),
  ...DEMAND_TYPES.map((type, index) => ({ id: `advanced-${String(index + 1).padStart(2, '0')}`, row: 'advanced' as const, quantity: 2 as const, icons: [type.id, DEMAND_TYPES[(index + 1) % DEMAND_TYPES.length].id], source: 'mvp-fixture' })),
]

const iconOf = (id: string) => DEMAND_TYPES.find((type) => type.id === id) ?? DEMAND_TYPES[0]

// A seeded shuffle keeps every player on the same random board while the seed
// can be replaced by the server's game seed when the API exposes it.
const shuffle = <T,>(items: T[], seed: number) => {
  const result = [...items]
  let value = seed || 1
  for (let index = result.length - 1; index > 0; index -= 1) {
    value = (value * 1664525 + 1013904223) >>> 0
    const target = value % (index + 1)
    ;[result[index], result[target]] = [result[target], result[index]]
  }
  return result
}

export function DemandCard({ card, revealed }: { card: DemandCardData; revealed: boolean }) {
  const style = { '--demand-accent': card.row === 'ordinary' ? '#c8b8a7' : '#aa9885' } as CSSProperties
  if (!revealed) return <div className="demand-card demand-card-back" style={style} aria-label="\u5c1a\u672a\u63ed\u793a\u9700\u6c42\u5361"><span>{card.row === 'ordinary' ? '?' : '??'}</span></div>
  return <div className={`demand-card demand-card-face ${card.row}`} style={style}><div className="demand-card-heading"><span>{card.row === 'ordinary' ? '\u666e\u901a\u9700\u6c42' : '\u9032\u968e\u9700\u6c42'}</span></div><div className="demand-card-icons">{card.icons.map((icon) => <span key={icon} title={iconOf(icon).label}>{iconOf(icon).symbol}</span>)}</div><div className="demand-card-labels">{card.icons.map((icon) => <small key={icon}>{iconOf(icon).label}</small>)}</div><div className="demand-card-bonus">\u9700\u6c42\u6578\u91cf<strong>{card.quantity}</strong></div></div>
}

export function DemandMarket({ period, round, embedded = false, seed = 1 }: { period: number; round: number; embedded?: boolean; seed?: number }) {
  const revealedColumn = Math.max(0, Math.min(5, round))
  const positioned = (row: DemandCardData['row']) => shuffle(DEMAND_CARDS.filter((card) => card.row === row), seed + (row === 'advanced' ? 97 : 0))
  const ordinary = positioned('ordinary')
  const advanced = positioned('advanced')
  return <section className={`demand-market ${embedded ? 'demand-market-embedded' : 'panel'}`} aria-label="\u9700\u6c42\u5361"><div className="demand-market-heading"><div><p className="eyebrow">DEMAND CARDS</p><h2>\u9700\u6c42\u5361</h2></div><span>\u7b2c {period} \u671f\u30fb\u7b2c {round} \u56de\u5408</span></div><div className="demand-card-board"><div className="demand-card-board-labels"><span>\u5361\u80cc</span><span>\u666e\u901a\u9700\u6c42\u5361</span><span>\u9032\u968e\u9700\u6c42\u5361</span></div><div className="demand-card-board-grid">{DEMAND_TYPES.map((type, column) => <div className="demand-card-column" key={type.id}><div className="demand-card-type-label"><span>{type.symbol}</span><small>{type.label}</small></div><DemandCard card={ordinary[column]} revealed={column <= revealedColumn} /><DemandCard card={advanced[column]} revealed={column <= revealedColumn} /></div>)}</div></div></section>
}
