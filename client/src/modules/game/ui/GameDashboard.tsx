// MUI 9 exposes several system props through `sx`; this screen is kept isolated
// while the project migrates its remaining legacy JSX to the same API.
// @ts-nocheck
import { useMemo, type CSSProperties } from 'react'
import { Coffee, Groups, Storefront, Paid, Psychology, Settings } from '@mui/icons-material'
import { Alert, Box, Button, Card as MuiCard, CardContent, Chip, Container, Divider, Grid, LinearProgress, Paper, Stack, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Toolbar, Typography } from '@mui/material'
import type { CardFaceData } from '../../../CardFace'
import { GameDashboardCardGroup } from './GameDashboardCardGroup'

type PlayerCard = CardFaceData
type DashboardPlayer = { id: string; displayName: string; bot?: boolean; handCount: number; cash: number; loans: number; revenue?: number; score?: number; brandAwareness?: number; products?: number; values?: number; resources?: number }
type CashFlowStatement = { period: number; beginningCash: number; operatingRevenue: number; otherIncome: number; operatingExpenses: number; interestPaid: number; principalRepayment: number; newLoans: number; endingCash: number }
type DashboardRoom = { period: number; round: number; phase: string; players: DashboardPlayer[]; demandBoard?: Record<string, number>; me?: { id: string; hand: PlayerCard[] | null; tableau: PlayerCard[] | null; partner?: PlayerCard; starterShop?: PlayerCard; cash: number; loans: number; revenue?: number; score?: number; selectedKPIs?: string[]; cashFlow?: CashFlowStatement[] | null; cashFlowRounds?: CashFlowStatement[] | null; brandAwareness?: number; products?: number; values?: number; resources?: number } }

const kindMeta: Record<string, { label: string; color: string; background: string; icon: typeof Coffee }> = {
  partner: { label: '合夥人卡', color: '#714b38', background: '#f3e5d8', icon: Groups },
  starter_shop: { label: '創始店卡', color: '#287477', background: '#dff0ed', icon: Storefront },
  resource: { label: '資源', color: '#3976a6', background: '#e4eff7', icon: Coffee },
  product: { label: '產品', color: '#b98a25', background: '#f8efd1', icon: Coffee },
  value: { label: '價值', color: '#b44f52', background: '#f6e3e4', icon: Psychology },
  channel: { label: '通路', color: '#3f7d66', background: '#e2f0e8', icon: Storefront },
  marketing: { label: '行銷', color: '#7a5ba5', background: '#eee7f7', icon: Psychology },
}

const partnerMeta: Record<string, { color: string; pale: string; icon: typeof Coffee; role: string }> = {
  'partner-barista': { color: '#2f87a8', pale: '#dceff5', icon: Coffee, role: '咖啡師・營運專業' },
  'partner-roaster': { color: '#9a623b', pale: '#f3e3d4', icon: Coffee, role: '烘豆師・產品專業' },
  'partner-marketer': { color: '#7656a5', pale: '#eee5f7', icon: Psychology, role: '行銷師・品牌專業' },
  'partner-service': { color: '#39826c', pale: '#dff0e8', icon: Groups, role: '服務設計・體驗專業' },
}

function PartnerCardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const meta = partnerMeta[card.id] ?? { color: '#714b38', pale: '#f3e5d8', icon: Groups, role: '創業夥伴' }
  const Icon = meta.icon
  return <MuiCard onClick={onClick} variant="outlined" sx={{ height: '100%', minHeight: 290, overflow: 'hidden', cursor: onClick ? 'pointer' : 'default', bgcolor: meta.pale, borderColor: selected ? meta.color : `${meta.color}66`, borderWidth: selected ? 3 : 1, transition: 'transform .15s, box-shadow .15s', '&:hover': onClick ? { transform: 'translateY(-3px)', boxShadow: 4 } : undefined }}>
    <Box sx={{ px: 1.5, py: .9, bgcolor: meta.color, color: 'white', textAlign: 'center' }}><Typography variant="caption" fontWeight={900} sx={{ letterSpacing: '.08em' }}>{card.name}</Typography></Box>
    <Box sx={{ px: 1.5, py: .85, display: 'flex', alignItems: 'center', gap: 1, bgcolor: `${meta.color}dd`, color: 'white' }}><Settings sx={{ fontSize: 19 }} /><Box><Typography variant="caption" sx={{ display: 'block', opacity: .8 }}>卡片功能</Typography><Typography variant="body2" fontWeight={900}>{meta.role}</Typography></Box></Box>
    <Box sx={{ display: 'grid', placeItems: 'center', minHeight: 105, color: meta.color, background: `linear-gradient(135deg, ${meta.pale} 0%, #ffffff 100%)` }}><Icon sx={{ fontSize: 64, opacity: .9 }} /></Box>
    <Box sx={{ px: 1.5, py: 1, bgcolor: '#ffffffcc' }}><Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: .35 }}>卡片說明</Typography><Typography variant="body2" sx={{ minHeight: 34 }}>{card.description ?? '創業合夥人提供的特殊能力。'}</Typography><Typography variant="caption" fontWeight={900} color={meta.color}>{card.effect ?? '—'}</Typography></Box>
    <Stack direction="row" justifyContent="space-between" sx={{ px: 1.5, py: .8, bgcolor: meta.pale }}><Typography variant="caption" color="text.secondary">成本</Typography><Typography variant="body2" fontWeight={900} color={meta.color}>${card.cost?.cash ?? 0} 萬</Typography></Stack>
  </MuiCard>
}

const shopMeta: Record<string, { color: string; pale: string; role: string }> = {
  'starter-songshan': { color: '#2f7f82', pale: '#dff1ef', role: '饕客聚集' },
  'starter-minsheng': { color: '#3e7896', pale: '#e1eff5', role: '饕客與一般客' },
  'starter-xinyi': { color: '#8b6a38', pale: '#f3ecd8', role: '分店客群拓展' },
  'starter-station': { color: '#735b91', pale: '#eee7f5', role: '一般顧客聚集' },
}

function StarterShopCardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const meta = shopMeta[card.id] ?? { color: '#287477', pale: '#dff0ed', role: '店面顧客來源' }
  const gourmet = card.demand?.gourmet ?? 0
  const regular = card.demand?.regular ?? 0
  return <MuiCard onClick={onClick} variant="outlined" sx={{ height: '100%', minHeight: 290, overflow: 'hidden', cursor: onClick ? 'pointer' : 'default', bgcolor: meta.pale, borderColor: selected ? meta.color : `${meta.color}66`, borderWidth: selected ? 3 : 1, transition: 'transform .15s, box-shadow .15s', '&:hover': onClick ? { transform: 'translateY(-3px)', boxShadow: 4 } : undefined }}>
    <Box sx={{ px: 1.5, py: .9, bgcolor: meta.color, color: 'white', textAlign: 'center' }}><Typography variant="caption" fontWeight={900} sx={{ letterSpacing: '.08em' }}>{card.name}</Typography></Box>
    <Box sx={{ px: 1.5, py: .85, display: 'flex', alignItems: 'center', gap: 1, bgcolor: `${meta.color}dd`, color: 'white' }}><Storefront sx={{ fontSize: 20 }} /><Box><Typography variant="caption" sx={{ display: 'block', opacity: .8 }}>店面功能</Typography><Typography variant="body2" fontWeight={900}>{meta.role}</Typography></Box></Box>
    <Box sx={{ display: 'grid', placeItems: 'center', minHeight: 105, color: meta.color, background: `linear-gradient(135deg, ${meta.pale} 0%, #ffffff 100%)` }}><Storefront sx={{ fontSize: 64, opacity: .9 }} /></Box>
    <Box sx={{ px: 1.5, py: 1, bgcolor: '#ffffffcc' }}><Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: .35 }}>顧客效果</Typography><Stack direction="row" spacing={1.5} alignItems="center" sx={{ minHeight: 34 }}>{gourmet > 0 && <Typography variant="body2" fontWeight={900} sx={{ color: '#ff7900', textShadow: '0 1px 1px rgba(160, 70, 0, .22)' }}>{'★'.repeat(gourmet)} <Box component="span" sx={{ color: '#6f625a', fontSize: 12 }}>饕客</Box></Typography>}{regular > 0 && <Typography variant="body2" fontWeight={900} sx={{ color: '#e5b832' }}>{'★'.repeat(regular)} <Box component="span" sx={{ color: '#6f625a', fontSize: 12 }}>一般客</Box></Typography>}{gourmet === 0 && regular === 0 && <Typography variant="body2">{card.effect ?? '開店後可獲得對應顧客。'}</Typography>}</Stack><Typography variant="caption" color="text.secondary">{card.description}</Typography></Box>
    <Stack direction="row" justifyContent="space-between" sx={{ px: 1.5, py: .8, bgcolor: meta.pale }}><Typography variant="caption" color="text.secondary">店面成本</Typography><Typography variant="body2" fontWeight={900} color={meta.color}>${card.cost?.cash ?? 0} 萬</Typography></Stack>
  </MuiCard>
}

const managementMeta: Record<string, { color: string; pale: string; label: string; function: string }> = {
  resource: { color: '#2d6897', pale: '#d5e7ef', label: '資源／營運', function: '取得營運資源' },
  product: { color: '#c88d28', pale: '#f3dfb7', label: '產品／商品', function: '推出特色產品' },
  value: { color: '#bd584f', pale: '#f3d9d3', label: '顧客／價值', function: '提升顧客價值' },
  channel: { color: '#2d6897', pale: '#d5e7ef', label: '資源／營運', function: '拓展銷售通路' },
  marketing: { color: '#bd584f', pale: '#f3d9d3', label: '顧客／價值', function: '推廣品牌與服務' },
}

const managementIcons: Record<string, string> = { operations: '⚙', coffee: '☕', value: '♥', channel: '⚑', marketing: '⚑', people: '♟' }

function MarketStars({ marketChange }: { marketChange?: Record<string, number> }) {
  const entries = Object.entries(marketChange ?? {}).filter(([key, value]) => (key === 'gourmet' || key === 'regular') && value !== 0)
  if (!entries.length) return <span className="management-market-empty">無市場變動</span>
  return <span className="management-market-stars">{entries.map(([key, value]) => <span className={key === 'gourmet' ? 'management-stars-gourmet' : 'management-stars-regular'} key={key}>{value > 0 ? '★'.repeat(value) : `−${Math.abs(value)}`}<small>{key === 'gourmet' ? '饕客' : key === 'regular' ? '一般客' : '奧客'}</small></span>)}</span>
}

function LegacyManagementCardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const meta = managementMeta[card.kind] ?? managementMeta.resource
  const title = card.name
  return <article className={`management-card ${selected ? 'is-selected' : ''}`} style={{ '--management-color': meta.color, '--management-pale': meta.pale } as CSSProperties} onClick={onClick} role={onClick ? 'button' : undefined} tabIndex={onClick ? 0 : undefined}>
    <header className="management-card-title"><strong>{title}</strong><small>第 {card.period} 期</small></header>
    <div className="management-card-function"><div className="management-icon-row">{card.icons.slice(0, 4).map((icon, index) => <span key={`${icon}-${index}`}>{managementIcons[icon] ?? '◆'}</span>)}</div><strong>{meta.function}</strong></div>
    <div className="management-card-description"><b>卡片說明</b><p>{card.description ?? '可執行的經營管理行動。'}</p></div>
    <div className="management-card-art"><span>{managementIcons[card.icons[0]] ?? '◆'}</span></div>
    <footer className="management-card-footer"><div><b>市場變動</b><MarketStars marketChange={card.marketChange} /></div><div className="management-card-cost"><b>成本</b><strong>{card.cost?.cash ?? 0}</strong></div></footer>
  </article>
}

function ManagementCardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const meta = managementMeta[card.kind] ?? managementMeta.resource
  const marketChange = Object.entries(card.marketChange ?? {}).filter(([, value]) => value !== 0).map(([key, value]) => `${key === 'gourmet' ? '饕客' : key === 'regular' ? '一般客' : '奧客'} ${value > 0 ? '+' : ''}${value}`).join('／') || '無市場變動'
  return <MuiCard onClick={onClick} variant="outlined" sx={{ height: '100%', minHeight: 290, overflow: 'hidden', cursor: onClick ? 'pointer' : 'default', bgcolor: meta.pale, borderColor: selected ? meta.color : `${meta.color}66`, borderWidth: selected ? 3 : 1, transition: 'transform .15s, box-shadow .15s', '&:hover': onClick ? { transform: 'translateY(-3px)', boxShadow: 4 } : undefined }}>
    <Box sx={{ px: 1.5, py: .9, bgcolor: meta.color, color: 'white', textAlign: 'center' }}><Typography variant="caption" fontWeight={900}>{card.name}</Typography></Box>
    <Box sx={{ px: 1.5, py: .85, display: 'flex', alignItems: 'center', gap: 1, bgcolor: `${meta.color}dd`, color: 'white' }}><Box sx={{ display: 'flex', gap: .4 }}>{card.icons.slice(0, 4).map((icon, index) => <Typography key={`${icon}-${index}`} sx={{ fontSize: 18 }}>{managementIcons[icon] ?? '◆'}</Typography>)}</Box><Box><Typography variant="caption" sx={{ display: 'block', opacity: .8 }}>卡片功能</Typography><Typography variant="body2" fontWeight={900}>{meta.function}</Typography></Box></Box>
    <Box sx={{ display: 'grid', placeItems: 'center', minHeight: 88, color: meta.color, background: 'linear-gradient(145deg, #f0f0ed, #dfe1df)' }}><Typography sx={{ fontSize: 46, opacity: .82 }}>{managementIcons[card.icons[0]] ?? '◆'}</Typography></Box>
    <Box sx={{ px: 1.5, py: 1, bgcolor: '#ffffffcc' }}><Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: .35 }}>卡片說明</Typography><Typography variant="body2" sx={{ minHeight: 34 }}>{card.description ?? '可執行的經營管理行動。'}</Typography></Box>
    <Stack direction="row" justifyContent="space-between" sx={{ px: 1.5, py: .8, bgcolor: meta.pale }}><Box><Typography variant="caption" fontWeight={900} color={meta.color}>{marketChange}</Typography><Typography display="block" variant="caption" color="text.secondary">市場變動</Typography></Box><Box><Typography variant="body2" fontWeight={900} color={meta.color}>${card.cost?.cash ?? 0}</Typography><Typography display="block" variant="caption" color="text.secondary">成本・至少 {card.minPlayers ?? 1} 人</Typography></Box></Stack>
  </MuiCard>
}

function CardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  if (card.kind === 'partner') return <PartnerCardTile card={card} selected={selected} onClick={onClick} />
  if (card.kind === 'starter_shop') return <StarterShopCardTile card={card} selected={selected} onClick={onClick} />
  return <LegacyManagementCardTile card={card} selected={selected} onClick={onClick} />
  const meta = kindMeta[card.kind] ?? { label: card.kind, color: '#765341', background: '#f1e8df', icon: Coffee }
  const Icon = meta.icon
  return <MuiCard onClick={onClick} variant="outlined" sx={{ height: '100%', minHeight: 172, cursor: onClick ? 'pointer' : 'default', bgcolor: meta.background, borderColor: selected ? meta.color : `${meta.color}55`, borderWidth: selected ? 2 : 1, transition: 'transform .15s, box-shadow .15s', '&:hover': onClick ? { transform: 'translateY(-3px)', boxShadow: 3 } : undefined }}>
    <Box sx={{ height: 6, bgcolor: meta.color }} />
    <CardContent sx={{ display: 'flex', height: 'calc(100% - 6px)', flexDirection: 'column', gap: 1.1 }}>
      <Stack direction="row" justifyContent="space-between" alignItems="center"><Chip size="small" label={meta.label} sx={{ color: meta.color, bgcolor: `${meta.color}18`, fontWeight: 700 }} /><Typography variant="caption" color="text.secondary">第 {card.period} 期</Typography></Stack>
      <Box sx={{ display: 'grid', placeItems: 'center', height: 44, borderRadius: 2, color: meta.color, bgcolor: `${meta.color}12` }}><Icon /></Box>
      <Typography variant="subtitle1" fontWeight={800}>{card.name}</Typography>
      <Typography variant="caption" color="text.secondary" sx={{ minHeight: 30 }}>{card.description ?? '可用於咖啡館經營的卡牌。'}</Typography>
      <Box sx={{ px: 1, py: .7, borderRadius: 1, bgcolor: '#ffffffaa' }}><Typography variant="caption" fontWeight={800} color={meta.color}>{card.kind === 'starter_shop' ? '顧客效果' : card.kind === 'partner' ? '合夥人功能' : '卡牌效果'}</Typography><Typography variant="body2" fontWeight={800}>{card.effect ?? '—'}</Typography></Box>
      <Stack direction="row" justifyContent="space-between" sx={{ mt: 'auto' }}><Typography variant="caption" color="text.secondary">成本</Typography><Typography variant="body2" fontWeight={800} color={meta.color}>${card.cost?.cash ?? 0} 萬</Typography></Stack>
    </CardContent>
  </MuiCard>
}

function CardGroup({ title, subtitle, cards, selectedId, onSelect }: { title: string; subtitle: string; cards: PlayerCard[]; selectedId?: string; onSelect?: (id: string) => void }) {
  if (!cards.length) return null
  return <Box sx={{ mt: 2.5 }}><Stack direction="row" alignItems="baseline" spacing={1} sx={{ mb: 1.2 }}><Typography variant="subtitle1" fontWeight={800}>{title}</Typography><Typography variant="caption" color="text.secondary">{subtitle}</Typography></Stack><Grid container spacing={1.5}>{cards.map((card) => <Grid key={card.id} size={{ xs: 12, sm: 6, md: 3, lg: 2 }}><CardTile card={card} selected={selectedId === card.id} onClick={onSelect ? () => onSelect(card.id) : undefined} /></Grid>)}</Grid></Box>
}

function CompactGamePanel({ room, market }: { room: DashboardRoom; market: [string, number][] }) {
  const metrics = [['品牌知名度', room.me?.brandAwareness ?? 0], ['產品能力', room.me?.products ?? 0], ['價值主張', room.me?.values ?? 0], ['資源能力', room.me?.resources ?? 0]]
  return <Paper sx={{ p: { xs: 2, md: 2.5 }, mb: 2.5, borderRadius: 3, border: '1px solid', borderColor: 'divider' }}><Stack direction={{ xs: 'column', md: 'row' }} justifyContent="space-between" spacing={2}><Box><Typography variant="overline" color="primary">GAME BOARD</Typography><Typography variant="h5" fontWeight={900}>遊戲面板</Typography><Typography variant="body2" color="text.secondary">第 {room.period} 期・第 {room.round} 回合・{room.phase}</Typography></Box><Stack direction="row" spacing={1} alignItems="center" sx={{ minWidth: { md: 280 } }}>{[1, 2, 3].map((period) => <Box key={period} sx={{ flex: 1 }}><Stack direction="row" justifyContent="space-between"><Typography variant="caption">第 {period} 期</Typography><Typography variant="caption" color={room.period >= period ? 'primary' : 'text.disabled'}>{room.period >= period ? '進行中' : '未開始'}</Typography></Stack><LinearProgress variant="determinate" value={room.period >= period ? 100 : 0} color={room.period === period ? 'primary' : 'secondary'} /></Box>)}</Stack></Stack><Divider sx={{ my: 2 }} /><Grid container spacing={1.5}>{metrics.map(([label, value]) => <Grid key={label} size={{ xs: 6, sm: 3 }}><Paper variant="outlined" sx={{ p: 1.4, bgcolor: '#faf7f2' }}><Typography variant="caption" color="text.secondary">{label}</Typography><Typography variant="h6" fontWeight={900}>{value}</Typography></Paper></Grid>)}{market.map(([key, value]) => <Grid key={key} size={{ xs: 6, sm: 3 }}><Paper variant="outlined" sx={{ p: 1.4, bgcolor: '#f1f8f7' }}><Typography variant="caption" color="text.secondary">市場・{key}</Typography><Typography variant="h6" fontWeight={900}>{value}</Typography></Paper></Grid>)}</Grid></Paper>
}

const referenceKpis = [
  ['品牌知名度', '品牌聲量'], ['現金', '現金餘額'], ['通路', '銷售通路'],
  ['顧客關係', '回訪顧客'], ['產品品質', '產品品質'], ['服務體驗', '服務體驗'],
  ['產品組合', '產品組合'], ['營運效率', '營運效率'], ['成本控制', '成本控制'],
]

function GamePanel({ room, market }: { room: DashboardRoom; market: [string, number][] }) {
  const periodNames = ['試營運', '正式營運', '擴大營運']
  return <Paper sx={{ mb: 2.5, overflow: 'hidden', borderRadius: 2, bgcolor: '#4d382b', color: '#fff8ef', boxShadow: '0 12px 28px rgba(73,47,33,.2)' }}>
    <Box sx={{ px: { xs: 2, md: 3 }, py: 1.8, display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid rgba(255,255,255,.16)' }}><Stack direction="row" spacing={1.2} alignItems="center"><Coffee sx={{ color: '#e4b27b' }} /><Box><Typography variant="overline" sx={{ color: '#dfc6b2', letterSpacing: '.12em' }}>CAFE STARTUPS</Typography><Typography variant="h6" fontWeight={900}>遊戲面板</Typography></Box></Stack><Chip label={`第 ${room.period} 期・第 ${room.round} 回合`} sx={{ color: '#f7dfc9', bgcolor: 'rgba(255,255,255,.12)', border: '1px solid rgba(255,255,255,.18)' }} /></Box>
    <Box sx={{ p: { xs: 1.5, md: 2.5 } }}>
      <Typography variant="caption" sx={{ color: '#dfc6b2' }}>關鍵指標區・選擇兩項觀察你的創業假設</Typography>
      <Grid container spacing={1} sx={{ mt: 1, mb: 2.5 }}>{referenceKpis.map(([label, detail], index) => <Grid key={label} size={{ xs: 6, sm: 4, md: 2.4 }}><Paper sx={{ p: 1.2, minHeight: 68, bgcolor: index < 2 ? '#805b43' : '#604536', color: '#fff8ef', border: '1px solid rgba(255,255,255,.12)', borderRadius: 1.5 }}><Typography variant="caption" sx={{ color: '#f1d8c1', display: 'block' }}>{label}</Typography><Typography variant="body2" fontWeight={800}>{detail}</Typography></Paper></Grid>)}</Grid>
      <Box sx={{ borderRadius: 1.5, overflow: 'hidden', border: '1px solid rgba(255,255,255,.15)' }}>
        <Box sx={{ display: 'grid', gridTemplateColumns: '1.25fr repeat(3, 1fr)', bgcolor: '#765341', color: '#f6dfca' }}><Typography sx={{ p: 1.2 }} variant="caption">需求市場</Typography>{periodNames.map((name, index) => <Box key={name} sx={{ p: 1.2, textAlign: 'center', bgcolor: room.period === index + 1 ? '#9a6847' : 'transparent' }}><Typography variant="caption" fontWeight={800}>第 {index + 1} 期</Typography><Typography display="block" variant="caption" sx={{ opacity: .75 }}>{name}</Typography></Box>)}</Box>
        {['gourmet', 'regular', 'difficult'].map((kind) => <Box key={kind} sx={{ display: 'grid', gridTemplateColumns: '1.25fr repeat(3, 1fr)', bgcolor: '#5b4335', borderTop: '1px solid rgba(255,255,255,.1)' }}><Typography sx={{ p: 1.2 }} variant="body2">{kind === 'gourmet' ? '饕客' : kind === 'regular' ? '一般客' : '奧客'}</Typography>{[1, 2, 3].map((period) => <Typography key={period} sx={{ p: 1.2, textAlign: 'center', bgcolor: room.period === period ? 'rgba(228,178,123,.18)' : 'transparent' }} variant="body2" fontWeight={700}>{market.find(([key]) => key === kind)?.[1] ?? 0}</Typography>)}</Box>)}
      </Box>
      <Stack direction={{ xs: 'column', md: 'row' }} spacing={1.5} sx={{ mt: 2 }}><Paper sx={{ p: 1.5, flex: 1, bgcolor: '#5b4335', color: '#f6dfca' }}><Typography variant="caption">目前現金</Typography><Typography variant="h6" fontWeight={900}>${room.me?.cash ?? 0} 萬</Typography></Paper><Paper sx={{ p: 1.5, flex: 1, bgcolor: '#5b4335', color: '#f6dfca' }}><Typography variant="caption">玩家排名</Typography><Typography variant="h6" fontWeight={900}>{room.players.slice().sort((a, b) => (b.score ?? 0) - (a.score ?? 0)).findIndex((player) => player.id === room.me?.id) + 1 || '—'} / {room.players.length}</Typography></Paper><Paper sx={{ p: 1.5, flex: 1, bgcolor: '#5b4335', color: '#f6dfca' }}><Typography variant="caption">遊戲狀態</Typography><Typography variant="h6" fontWeight={900}>{room.phase === 'finished' ? '已完成' : '進行中'}</Typography></Paper></Stack>
    </Box>
  </Paper>
}

function CashFlowPanel({ room }: { room: DashboardRoom }) {
  const statements = new Map((room.me?.cashFlow ?? []).map((statement) => [statement.period, statement]))
  const rounds = new Map((room.me?.cashFlowRounds ?? []).map((statement) => [statement.period, statement]))
  const current = (period: number) => period === room.period ? rounds.get(period) ?? statements.get(period) : statements.get(period)
  const money = (value: number | undefined) => value === undefined ? '—' : `$${value}`
  const rows: [string, keyof CashFlowStatement][] = [['期初現金', 'beginningCash'], ['營業收入', 'operatingRevenue'], ['其他收入', 'otherIncome'], ['營運支出', 'operatingExpenses'], ['支付利息', 'interestPaid'], ['償還本金', 'principalRepayment'], ['新增貸款', 'newLoans'], ['期末現金', 'endingCash']]
  return <Paper sx={{ mt: 2.5, borderRadius: 3, overflow: 'hidden' }}><Box sx={{ p: 2.2 }}><Typography variant="overline" color="primary">CASH FLOW STATEMENT</Typography><Typography variant="h5" fontWeight={900}>現金流量表</Typography><Typography variant="body2" color="text.secondary">依期別整理目前的現金流入、支出與期末結餘。</Typography></Box><TableContainer><Table size="small"><TableHead><TableRow><TableCell>項目</TableCell>{[1, 2, 3].map((period) => <TableCell align="right" key={period}>第 {period} 期</TableCell>)}</TableRow></TableHead><TableBody>{rows.map(([label, field]) => <TableRow key={field} sx={field === 'endingCash' ? { bgcolor: '#f5eee8' } : undefined}><TableCell sx={{ fontWeight: field === 'endingCash' ? 800 : 400 }}>{label}</TableCell>{[1, 2, 3].map((period) => <TableCell align="right" key={period} sx={{ fontWeight: field === 'endingCash' ? 800 : 400 }}>{money(current(period)?.[field])}</TableCell>)}</TableRow>)}</TableBody></Table></TableContainer></Paper>
}

export function GameDashboard({ room, selectedCard, setSelectedCard, command, busy, error, leave }: { room: DashboardRoom; selectedCard: string; setSelectedCard: (id: string) => void; command: (type: string, extra?: Record<string, unknown>) => void; busy: boolean; error: string; leave: () => void }) {
  const me = room.me
  const startingCards = [me?.partner, me?.starterShop].filter((card): card is PlayerCard => Boolean(card))
  const hand = me?.hand ?? []
  const tableau = me?.tableau ?? []
  const selected = hand.find((card) => card.id === selectedCard)
  const phaseLabel: Record<string, string> = { hypothesis: '設定關鍵指標', experiment: '選擇實驗策略', learning: '結算本期市場', finished: '遊戲完成' }
  const totalCards = startingCards.length + hand.length + tableau.length
  const market = useMemo(() => Object.entries(room.demandBoard ?? {}), [room.demandBoard])

  return <Box sx={{ minHeight: '100vh', bgcolor: '#f8f5f1', color: '#33251f' }}>
    <Box sx={{ bgcolor: '#4f3428', color: 'white' }}><Container maxWidth="xl"><Toolbar disableGutters sx={{ justifyContent: 'space-between' }}><Stack direction="row" spacing={1.5} alignItems="center"><Coffee /><Typography variant="h6" fontWeight={900}>Cafe Startups</Typography></Stack><Stack direction="row" spacing={1}><Chip label={`第 ${room.period} 期`} sx={{ color: 'white', bgcolor: '#ffffff20' }} /><Button color="inherit" size="small" onClick={leave}>離開</Button></Stack></Toolbar></Container></Box>
    <Container maxWidth="xl" sx={{ py: 3 }}>
      <Stack direction={{ xs: 'column', md: 'row' }} justifyContent="space-between" alignItems={{ xs: 'start', md: 'center' }} spacing={1} sx={{ mb: 2.5 }}><Box><Typography variant="overline" color="text.secondary">PERIOD {room.period} · ROUND {room.round}/6</Typography><Typography variant="h4" fontWeight={900}>{phaseLabel[room.phase] ?? room.phase}</Typography></Box><Chip icon={<Coffee />} label={`${totalCards} 張卡牌`} color="primary" variant="outlined" /></Stack>
      <Paper elevation={0} sx={{ p: { xs: 2, md: 2.5 }, border: '1px solid', borderColor: 'divider', borderRadius: 3, mb: 2.5 }}><Stack direction={{ xs: 'column', md: 'row' }} justifyContent="space-between" spacing={1}><Box><Typography variant="overline" color="primary">YOUR COLLECTION</Typography><Typography variant="h5" fontWeight={900}>我的卡牌</Typography></Box><Typography variant="body2" color="text.secondary" sx={{ maxWidth: 420 }}>開局取得的合夥人卡、創始店卡，以及目前手牌與已打出的卡牌都會顯示在這裡。</Typography></Stack><GameDashboardCardGroup title="開局卡牌" subtitle="合夥人卡・創始店卡" cards={startingCards} /><GameDashboardCardGroup title="目前手牌" subtitle={`${hand.length} 張・點選卡牌後執行操作`} cards={hand} selectedId={selectedCard} onSelect={(id) => { setSelectedCard(id); command('SELECT_CARD', { cardId: id }) }} /><GameDashboardCardGroup title="已打出卡牌" subtitle={`${tableau.length} 張`} cards={tableau} /></Paper>
      <Grid container spacing={2.5}>
        <Grid size={{ xs: 12, md: 3 }}><Stack spacing={2}><Paper sx={{ p: 2.2, borderRadius: 3 }}><Typography variant="overline" color="text.secondary">我的咖啡館</Typography><Typography variant="h3" fontWeight={900} sx={{ my: 1 }}>${me?.cash ?? 0}<Typography component="span" variant="body2" color="text.secondary"> 萬</Typography></Typography><Divider sx={{ mb: 1 }} />{[['貸款', `${me?.loans ?? 0} / 6`], ['營收', `$${me?.revenue ?? 0} 萬`], ['分數', `${me?.score ?? '—'}`], ['品牌 / 產品', `${me?.brandAwareness ?? 0} / ${me?.products ?? 0}`], ['價值 / 資源', `${me?.values ?? 0} / ${me?.resources ?? 0}`]].map(([label, value]) => <Stack key={label} direction="row" justifyContent="space-between" sx={{ py: .8 }}><Typography variant="body2" color="text.secondary">{label}</Typography><Typography variant="body2" fontWeight={800}>{value}</Typography></Stack>)}</Paper><Button fullWidth variant="outlined" startIcon={<Paid />} onClick={() => command('TAKE_LOAN')} disabled={busy}>取得一筆貸款</Button></Stack></Grid>
        <Grid size={{ xs: 12, md: 9 }}><Stack spacing={2}><Paper sx={{ p: 2.2, borderRadius: 3, bgcolor: '#714b38', color: 'white' }}><Stack direction={{ xs: 'column', sm: 'row' }} justifyContent="space-between" spacing={1}><Box><Typography variant="overline" sx={{ opacity: .7 }}>目前階段</Typography><Typography variant="h6" fontWeight={900}>{phaseLabel[room.phase] ?? room.phase}</Typography></Box><Typography variant="body2" sx={{ opacity: .75 }}>{room.players.filter((player) => !player.handCount).length ? '等待玩家行動' : '系統同步中'}</Typography></Stack></Paper><Paper sx={{ p: 2.2, borderRadius: 3 }}><Typography variant="h6" fontWeight={900}>需求市場</Typography><Grid container spacing={1.5} sx={{ mt: .5 }}>{market.length ? market.map(([key, value]) => <Grid key={key} size={{ xs: 6, sm: 3 }}><Paper variant="outlined" sx={{ p: 1.5, bgcolor: '#faf7f2' }}><Typography variant="caption" color="text.secondary">{key}</Typography><Typography variant="h6" fontWeight={900}>{value}</Typography></Paper></Grid>) : <Grid size={12}><Typography variant="body2" color="text.secondary">市場資料將在本期開始後顯示。</Typography></Grid>}</Grid></Paper>{room.phase === 'experiment' && <Stack direction={{ xs: 'column', sm: 'row' }} justifyContent="space-between" spacing={1}><Button variant="outlined" onClick={() => selected && command('DISCARD_SELECTED_CARD')} disabled={busy || !selected}>棄牌 +20</Button><Button variant="contained" onClick={() => selected && command('PLAY_SELECTED_CARD')} disabled={busy || !selected}>打出 {selected ? `$${selected.cost.cash} 萬` : '選擇一張牌'}</Button><Button variant="contained" color="secondary" onClick={() => command('PASS_HAND')} disabled={busy}>完成選擇並傳牌 →</Button></Stack>}{room.phase === 'learning' && <Button variant="contained" onClick={() => command('RESOLVE_LEARNING')} disabled={busy}>結算本期並繼續</Button>}{error && <Alert severity="error">{error}</Alert>}</Stack></Grid>
      </Grid>
      <Paper sx={{ mt: 2.5, p: 2, borderRadius: 3 }}><Typography variant="subtitle1" fontWeight={900} sx={{ mb: 1 }}>玩家</Typography><Stack direction="row" flexWrap="wrap" gap={1}>{room.players.map((player) => <Chip key={player.id} icon={<Groups />} label={`${player.displayName}${player.id === me?.id ? '（你）' : player.bot ? '（電腦）' : ''} · ${player.handCount} 張`} variant={player.id === me?.id ? 'filled' : 'outlined'} color={player.id === me?.id ? 'primary' : 'default'} />)}</Stack></Paper>
    </Container>
  </Box>
}
