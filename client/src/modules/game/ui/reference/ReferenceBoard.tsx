import { arrangeDemandCards, demandQuantityForPosition, DemandCard, DemandCardData } from '../market/DemandMarket'
import { AttachMoney, Coffee, Star, Verified } from '@mui/icons-material'

const REFERENCE_KPIS = [
  { label: '饕客滿意度', value: '每張 +5 分', tone: 'gourmet', symbol: '' },
  { label: '盈餘', value: '每 30 → +1', tone: 'cash', symbol: '30' },
  { label: '通路', value: '每張 +4 分', tone: 'channel', symbol: '▰' },
  { label: '一般客滿意度', value: '每張 +4 分', tone: 'regular', symbol: '' },
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

export function ReferenceBoard(props: { period: number; phase?: string; round?: number; seed?: string; demandCards?: Record<'gourmet' | 'regular', { id: string; position: number; icons: string[]; revealed: boolean }[]> }) {
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
        <ReferenceDemandRow type="gourmet" label="饕客" base="$10" additions={REFERENCE_PERIODS.map((item) => item.gourmet)} activePeriod={props.period} revealed={props.phase === 'experiment' || props.phase === 'learning' || props.phase === 'finished'} round={props.round} seed={props.seed} demandCards={props.demandCards?.gourmet} />
        <ReferenceDemandRow type="regular" label="一般客" base="$10" additions={REFERENCE_PERIODS.map((item) => item.regular)} activePeriod={props.period} revealed={props.phase === 'experiment' || props.phase === 'learning' || props.phase === 'finished'} round={props.round} seed={props.seed} demandCards={props.demandCards?.regular} />
        </div></div><div className="reference-bottom">
          <div className="reference-summary"><span className="reference-cube blue">◆</span><span>關鍵資源</span><strong>$0</strong><div className="reference-metric-order" aria-label="知名度 → 品質 → 產品 → 成本"><b className="metric-awareness" title="知名度"><Star fontSize="inherit" /></b><b>→</b><b className="metric-quality" title="品質"><Verified fontSize="inherit" /></b><b>→</b><b className="metric-product" title="產品"><Coffee fontSize="inherit" /></b><b>→</b><b className="metric-cost" title="成本"><AttachMoney fontSize="inherit" /></b></div></div>
          <div className="reference-ranking"><div className="reference-ranking-title"><span>市場排名</span><small>抽取市場袋顧客數</small></div><div className="reference-ranking-grid"><div className="reference-rank-labels"><span>1st</span><span>2nd</span><span>3rd</span><span>4th</span></div>{REFERENCE_PERIODS.map((item) => <div className={props.period === item.id ? 'reference-rank-column is-current' : 'reference-rank-column'} key={item.id}>{item.customers.map((count, index) => <span key={`${item.id}-${index}`}><b>{count}</b><i>●</i></span>)}</div>)}</div></div>
        </div>
      </div>
    </div>
  </details>
}

function ReferenceDemandRow(props: { type: 'gourmet' | 'regular'; label: string; base: string; additions: string[]; activePeriod: number; revealed?: boolean; round?: number; seed?: string; demandCards?: { id: string; position: number; icons: string[]; revealed: boolean }[] }) {
  const questionMarks = (period: number) => '?'.repeat(demandQuantityForPosition(props.type, period))
  const numericSeed = [...(props.seed ?? 'local-game')].reduce((value, character) => ((value * 31) + character.charCodeAt(0)) >>> 0, 1)
  const cards = arrangeDemandCards(numericSeed)[props.type === 'gourmet' ? 'first' : 'second']
  const slot = (column: number, fallback: string) => {
    const savedCard = props.demandCards?.find((item) => item.position === column)
    const card: DemandCardData | undefined = savedCard ? { id: savedCard.id, kind: 'demand', icons: savedCard.icons, source: 'mvp-fixture' } : cards[column]
    const isRevealed = savedCard
      ? savedCard.revealed
      : Boolean(card && props.revealed && props.round !== undefined && props.round >= column)
    return isRevealed ? <DemandCard card={card!} revealed quantity={demandQuantityForPosition(props.type, column)} /> : <><b>{questionMarks(column)}</b><small>{fallback}</small></>
  }
  return <div className={`reference-demand-row reference-demand-${props.type}`}>
    <div className="reference-customer-label"><span className="reference-customer-figures"><i /><i /></span><strong>{props.label}</strong><small>{props.type === 'gourmet' ? '饕客需求 +10' : '一般客需求 +10'}</small></div>
    <div className="reference-demand-slot reference-demand-base">{slot(0, props.base)}</div>
    {props.additions.map((addition, index) => <div className={`reference-demand-slot ${props.activePeriod === index + 1 ? 'is-current' : ''}`} key={`${props.type}-${index}`}>{slot(index + 1, addition)}</div>)}
  </div>
}
