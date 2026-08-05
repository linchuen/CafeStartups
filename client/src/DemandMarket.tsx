import type { CSSProperties } from 'react'
import { GameIcon } from './modules/game/ui/GameIcon'

export type DemandCardData = { id: string; row: 'ordinary' | 'advanced'; quantity: 1 | 2; icons: string[]; source?: string }

const DEMAND_TYPES = [
  { id: 'coffee', label: '咖啡', color: '#b87932', pale: '#f4e4ce' },
  { id: 'dessert', label: '甜點', color: '#d18b3d', pale: '#f7e6c8' },
  { id: 'beans', label: '咖啡豆', color: '#8b6845', pale: '#ece1d2' },
  { id: 'taste', label: '可口美味', color: '#bd5c50', pale: '#f3dcd7' },
  { id: 'service', label: '貼心服務', color: '#b94e55', pale: '#f3d8db' },
  { id: 'price', label: '合理價格', color: '#ad4e4b', pale: '#f1d9d4' },
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

export const arrangeDemandCards = (seed = 1) => ({
  ordinary: shuffle(DEMAND_CARDS.filter((card) => card.row === 'ordinary'), seed),
  advanced: shuffle(DEMAND_CARDS.filter((card) => card.row === 'advanced'), seed + 97),
})

export function DemandCard({ card, revealed }: { card: DemandCardData; revealed: boolean }) {
  const firstType = iconOf(card.icons[0])
  const style = { '--demand-accent': firstType.color, '--demand-pale': firstType.pale } as CSSProperties
  if (!revealed) return <div className="demand-card demand-card-back" style={style} aria-label="\u5c1a\u672a\u63ed\u793a\u9700\u6c42\u5361"><span>{card.row === 'ordinary' ? '?' : '??'}</span></div>
  return <div className={`demand-card demand-card-face ${card.row}`} style={style}><div className="demand-card-heading"><span>{card.row === 'ordinary' ? '\u666e\u901a\u9700\u6c42' : '\u9032\u968e\u9700\u6c42'}</span></div><div className="demand-card-icons">{card.icons.map((icon) => <span key={icon} title={iconOf(icon).label}><GameIcon name={icon} fontSize="inherit" /></span>)}</div><div className="demand-card-labels">{card.icons.map((icon) => <small key={icon}>{iconOf(icon).label}</small>)}</div></div>
}

export function DemandMarket({ period, round, reveal = false, embedded = false, seed = 1 }: { period: number; round: number; reveal?: boolean; embedded?: boolean; seed?: number }) {
  const revealedColumn = reveal ? 5 : -1
  const { ordinary, advanced } = arrangeDemandCards(seed)
  return <section className={`demand-market ${embedded ? 'demand-market-embedded' : 'panel'}`} aria-label="\u9700\u6c42\u5361"><div className="demand-market-heading"><div><p className="eyebrow">DEMAND CARDS</p><h2>\u9700\u6c42\u5361</h2></div><span>\u7b2c {period} \u671f\u30fb\u7b2c {round} \u56de\u5408</span></div><div className="demand-card-board"><div className="demand-card-board-labels"><span>\u5361\u80cc</span><span>\u666e\u901a\u9700\u6c42\u5361</span><span>\u9032\u968e\u9700\u6c42\u5361</span></div><div className="demand-card-board-grid">{DEMAND_TYPES.map((type, column) => <div className="demand-card-column" key={type.id}><div className="demand-card-type-label" style={{ color: type.color }}><span><GameIcon name={type.id} fontSize="small" /></span><small>{type.label}</small></div><DemandCard card={ordinary[column]} revealed={column <= revealedColumn} /><DemandCard card={advanced[column]} revealed={column <= revealedColumn} /></div>)}</div></div></section>
}
