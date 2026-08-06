import type { CSSProperties } from 'react'
import { GameIcon } from '../cards/GameIcon'

export type DemandCustomerType = 'gourmet' | 'regular'
export type DemandCardData = { id: string; kind: 'demand'; icons: string[]; source?: string }

export const DEMAND_QUANTITIES: Record<DemandCustomerType, readonly number[]> = {
  gourmet: [1, 2, 2, 2],
  regular: [1, 1, 1, 2],
}

export function demandQuantityForPosition(customerType: DemandCustomerType, round: number) {
  return DEMAND_QUANTITIES[customerType][round] ?? DEMAND_QUANTITIES[customerType][DEMAND_QUANTITIES[customerType].length - 1]
}

const DEMAND_TYPES = [
  { id: 'coffee', label: '咖啡', color: '#b87932', pale: '#f4e4ce' },
  { id: 'dessert', label: '甜點', color: '#d18b3d', pale: '#f7e6c8' },
  { id: 'beans', label: '咖啡豆', color: '#8b6845', pale: '#ece1d2' },
  { id: 'taste', label: '可口美味', color: '#bd5c50', pale: '#f3dcd7' },
  { id: 'service', label: '貼心服務', color: '#b94e55', pale: '#f3d8db' },
  { id: 'price', label: '合理價格', color: '#ad4e4b', pale: '#f1d9d4' },
] as const

export const DEMAND_CARDS: DemandCardData[] = DEMAND_TYPES.flatMap((type, index) => [
  { id: `demand-${String(index + 1).padStart(2, '0')}`, kind: 'demand' as const, icons: [type.id], source: 'mvp-fixture' },
  { id: `demand-${String(index + 7).padStart(2, '0')}`, kind: 'demand' as const, icons: [type.id, DEMAND_TYPES[(index + 1) % DEMAND_TYPES.length].id], source: 'mvp-fixture' },
])

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
  first: shuffle(DEMAND_CARDS.slice(0, DEMAND_TYPES.length), seed),
  second: shuffle(DEMAND_CARDS.slice(DEMAND_TYPES.length), seed + 97),
})

export function DemandCard({ card, revealed, quantity = card.icons.length }: { card: DemandCardData; revealed: boolean; quantity?: number }) {
  const firstType = iconOf(card.icons[0])
  const style = { '--demand-accent': firstType.color, '--demand-pale': firstType.pale } as CSSProperties
  const visibleIcons = Array.from({ length: quantity }, (_, index) => card.icons[index % card.icons.length])
  if (!revealed) return <div className="demand-card demand-card-back" style={style} aria-label="\u5c1a\u672a\u63ed\u793a\u9700\u6c42\u5361"><span>?</span></div>
  return <div className="demand-card demand-card-face" style={style}><div className="demand-card-heading"><span>需求卡</span></div><div className="demand-card-icons">{visibleIcons.map((icon, index) => <span key={`${icon}-${index}`} title={iconOf(icon).label}><GameIcon name={icon} fontSize="inherit" /></span>)}</div><div className="demand-card-labels">{visibleIcons.map((icon, index) => <small key={`${icon}-${index}`}>{iconOf(icon).label}</small>)}</div></div>
}

export function DemandMarket({ period, round, reveal = false, embedded = false, seed = 1 }: { period: number; round: number; reveal?: boolean; embedded?: boolean; seed?: number }) {
  const revealedColumn = reveal ? Math.min(5, Math.max(-1, round)) : -1
  const { first, second } = arrangeDemandCards(seed)
  return <section className={`demand-market ${embedded ? 'demand-market-embedded' : 'panel'}`} aria-label="\u9700\u6c42\u5361"><div className="demand-market-heading"><div><p className="eyebrow">DEMAND CARDS</p><h2>\u9700\u6c42\u5361</h2></div><span>\u7b2c {period} \u671f\u30fb\u7b2c {round} \u56de\u5408</span></div><div className="demand-card-board"><div className="demand-card-board-labels"><span>\u5361\u80cc</span><span>\u9700\u6c42\u4f4d\u7f6e 1</span><span>\u9700\u6c42\u4f4d\u7f6e 2</span></div><div className="demand-card-board-grid">{DEMAND_TYPES.map((type, column) => <div className="demand-card-column" key={type.id}><div className="demand-card-type-label" style={{ color: type.color }}><span><GameIcon name={type.id} fontSize="small" /></span><small>{type.label}</small></div><DemandCard card={first[column]} revealed={column <= revealedColumn} /><DemandCard card={second[column]} revealed={column <= revealedColumn} /></div>)}</div></div></section>
}
