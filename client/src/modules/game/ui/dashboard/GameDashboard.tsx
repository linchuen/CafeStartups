// MUI 9 exposes several system props through `sx`; this screen is kept isolated
// while the project migrates its remaining legacy JSX to the same API.
// @ts-nocheck
import { useMemo, useState } from 'react'
import { Coffee, Groups, Paid } from '@mui/icons-material'
import { Alert, Box, Button, Chip, Container, Divider, Grid, MenuItem, Paper, Select, Stack, Toolbar, Typography } from '@mui/material'
import type { CardFaceData } from '../../model/cardTypes'
import type { CashFlowStatement, GameState } from '../../model/gameTypes'
import { GameDashboardCardGroup } from '../cards/GameDashboardCardGroup'
import { GameDashboardSetupCards } from '../cards/GameDashboardSetupCards'
import { CustomerTypeTokens } from '../cards/CustomerTypeTokens'
import { GameIcon } from '../cards/GameIcon'

type PlayerCard = CardFaceData
type DashboardRoom = GameState

function MarketRankingPanel({ room, command, busy }: { room: DashboardRoom; command: (type: string, extra?: Record<string, unknown>) => void; busy: boolean }) {
  if (room.phase !== 'learning') return null
  const ranking = room.marketRanking ?? []
  const draws = [...(room.marketDraws ?? [])].sort((left, right) => right.count - left.count)
  const typeLabels = { gourmet: '饕客', regular: '一般客', difficult: '困難客' }
  return <Paper sx={{ p: 2.2, borderRadius: 3, border: '2px solid #b36f42', bgcolor: '#fffaf4', mb: 2 }}><Stack direction={{ xs: 'column', sm: 'row' }} justifyContent="space-between" alignItems={{ xs: 'start', sm: 'center' }} spacing={1}><Box><Typography variant="overline" color="primary">MARKET RANKING</Typography><Typography variant="h6" fontWeight={900}>市場排名與抽取結果</Typography><Typography variant="body2" color="text.secondary">排名會標示玩家，並顯示抽到的客群與顧客數量。</Typography></Box><Button variant="contained" onClick={() => command(ranking.length ? 'RESOLVE_LEARNING' : 'DRAW_MARKET')} disabled={busy}>{ranking.length ? '確認排名並結算' : '抽取市場袋顧客數'}</Button></Stack>{draws.length > 0 ? <Grid container spacing={1.2} sx={{ mt: 1.5 }}>{draws.map((draw, index) => <Grid key={draw.playerId} size={{ xs: 12, sm: 6, md: 3 }}><Paper variant="outlined" sx={{ p: 1.2, bgcolor: '#fff' }}><Stack direction="row" justifyContent="space-between" alignItems="center"><Typography variant="caption" color="text.secondary">第 {index + 1} 名</Typography><Typography variant="subtitle2" fontWeight={900}>{draw.playerName}</Typography></Stack><Stack direction="row" spacing={1} alignItems="center" sx={{ mt: 1 }}><CustomerTypeTokens type={draw.customerType} count={draw.count} size={16} /><Box><Typography variant="body2" fontWeight={900}>{typeLabels[draw.customerType]}・{draw.count} 位</Typography><Typography variant="caption" color="text.secondary">本期抽取結果</Typography></Box></Stack></Paper></Grid>)}</Grid> : ranking.length > 0 && <Grid container spacing={1.2} sx={{ mt: 1.5 }}>{ranking.map((count, index) => <Grid key={index} size={{ xs: 6, sm: 3 }}><Paper variant="outlined" sx={{ p: 1.4, textAlign: 'center', bgcolor: '#fff' }}><Typography variant="caption" color="text.secondary">第 {index + 1} 名</Typography><Typography variant="h4" fontWeight={900} color="primary">{count}</Typography><Typography variant="caption">位顧客</Typography></Paper></Grid>)}</Grid>}</Paper>
}

const KPI_OPTIONS = [
  ['gourmet_satisfaction', '饕客滿意度'],
  ['regular_satisfaction', '一般客滿意度'],
  ['total_satisfaction', '總滿意度'],
  ['channel', '通路'],
  ['awareness', '品牌知名度'],
  ['products', '特色產品'],
  ['quality', '品質'],
  ['cost', '成本'],
  ['surplus', '盈餘'],
]

function KPISelectionPanel({ room, command, busy }: { room: DashboardRoom; command: (type: string, extra?: Record<string, unknown>) => void; busy: boolean }) {
  const saved = room.me?.selectedKPIs?.length === 2 ? room.me.selectedKPIs : ['awareness', 'products']
  const [first, setFirst] = useState(saved[0])
  const [second, setSecond] = useState(saved[1])
  if (room.phase !== 'hypothesis' || room.period <= 1) return null
  return <Paper sx={{ p: 2.2, borderRadius: 3, border: '2px solid #7656a5', bgcolor: '#f8f3fc', mb: 2 }}><Stack direction={{ xs: 'column', sm: 'row' }} justifyContent="space-between" alignItems={{ xs: 'start', sm: 'center' }} spacing={1}><Box><Typography variant="overline" sx={{ color: '#7656a5' }}>KEY METRICS</Typography><Typography variant="h6" fontWeight={900}>第 {room.period} 期設定關鍵指標</Typography><Typography variant="body2" color="text.secondary">選擇兩個不同指標，完成後會發放本期新卡。</Typography></Box><Stack direction="row" spacing={1}><Select size="small" value={first} onChange={(event) => setFirst(event.target.value)} disabled={busy}>{KPI_OPTIONS.map(([id, label]) => <MenuItem key={id} value={id}>{label}</MenuItem>)}</Select><Select size="small" value={second} onChange={(event) => setSecond(event.target.value)} disabled={busy}>{KPI_OPTIONS.map(([id, label]) => <MenuItem key={id} value={id}>{label}</MenuItem>)}</Select><Button variant="contained" onClick={() => command('SET_KPI', { kpis: [first, second] })} disabled={busy || first === second}>確認並發牌</Button></Stack></Stack></Paper>
}

function PlayerSummaryBar({ room }: { room: DashboardRoom }) {
  const me = room.me
  const metrics = [['餘額', `$${me?.cash ?? 0} 萬`], ['營收', `$${me?.revenue ?? 0} 萬`], ['貸款', `${me?.loans ?? 0} / 6`], ['品牌', `${me?.brandAwareness ?? 0}`], ['產品', `${me?.products ?? 0}`], ['價值', `${me?.values ?? 0}`], ['資源', `${me?.resources ?? 0}`]]
  return <Paper elevation={0} sx={{ p: 1.1, mb: 1.5, border: '1px solid', borderColor: 'divider', borderRadius: 2.5, bgcolor: '#fffaf4' }}><Stack direction="row" spacing={1} sx={{ overflowX: 'auto' }}>{metrics.map(([label, value]) => <Box key={label} sx={{ minWidth: { xs: 74, sm: 88 }, px: 1, py: .55, borderRight: '1px solid', borderColor: 'divider', '&:last-child': { borderRight: 0 } }}><Typography variant="caption" color="text.secondary" noWrap>{label}</Typography><Typography variant="body2" fontWeight={900} noWrap>{value}</Typography></Box>)}</Stack></Paper>
}

function CustomerCountsPanel({ counts }: { counts: Record<string, number> }) {
  return <Stack direction="row" spacing={1} sx={{ mb: 1.5 }}><Paper variant="outlined" sx={{ flex: 1, p: 1, bgcolor: '#fff7ed', textAlign: 'center' }}><CustomerTypeTokens type="gourmet" count={counts.gourmet ?? 0} size={18} /><Typography variant="caption" display="block">卡牌饕客數</Typography><Typography variant="h6" fontWeight={900}>{counts.gourmet ?? 0}</Typography></Paper><Paper variant="outlined" sx={{ flex: 1, p: 1, bgcolor: '#fffbea', textAlign: 'center' }}><CustomerTypeTokens type="regular" count={counts.regular ?? 0} size={18} /><Typography variant="caption" display="block">卡牌一般客數</Typography><Typography variant="h6" fontWeight={900}>{counts.regular ?? 0}</Typography></Paper></Stack>
}

function paymentBreakdown(card: PlayerCard, ownedIcons: string[]) {
  const available = [...ownedIcons]
  let missingIcons = 0
  for (const icon of card.cost?.icons ?? []) {
    const index = available.indexOf(icon)
    if (index === -1) missingIcons += 1
    else available.splice(index, 1)
  }
  return [`$${card.cost?.cash ?? 0}`, ...(missingIcons > 0 ? [`$${missingIcons * 20}`] : [])].join(' + ')
}

function HandActionBar({ room, selected, ownedIcons, command, busy, showLoan = true, showCardActions = true }: { room: DashboardRoom; selected?: PlayerCard; ownedIcons: string[]; command: (type: string, extra?: Record<string, unknown>) => void; busy: boolean; showLoan?: boolean; showCardActions?: boolean }) {
  return <Stack direction={{ xs: 'column', sm: 'row' }} spacing={1} sx={{ mb: 1.5 }}>{showLoan && <Button size="small" variant="outlined" onClick={() => command('TAKE_LOAN')} disabled={busy}>＋ 取得貸款</Button>}{showCardActions && room.phase === 'experiment' && <><Button size="small" variant="outlined" onClick={() => command('DISCARD_SELECTED_CARD')} disabled={busy || !selected}>棄牌並傳牌 +20</Button><Button size="small" variant="contained" onClick={() => command('PLAY_SELECTED_CARD')} disabled={busy || !selected}>打出並傳牌 {selected ? paymentBreakdown(selected, ownedIcons) : '選擇一張牌'}</Button></>}</Stack>
}

export function GameDashboard({ room, selectedCard, setSelectedCard, command, busy, error, leave, setupInitialCards }: { room: DashboardRoom; selectedCard: string; setSelectedCard: (id: string) => void; command: (type: string, extra?: Record<string, unknown>) => void; busy: boolean; error: string; leave: () => void; setupInitialCards: (partnerId: string, starterShopId: string) => void }) {
  const me = room.me
  const startingCards = (room.period === 0 && !me?.initialCardsSelected ? [] : [me?.partner, me?.starterShop]).filter((card): card is PlayerCard => Boolean(card))
  const hand = me?.hand ?? []
  const tableau = me?.tableau ?? []
  const acquiredCards = [...startingCards, ...tableau]
  const ownedIcons = acquiredCards.flatMap((card) => card.icons ?? [])
  const orderedAcquiredCards = [...acquiredCards].sort((left, right) => {
    const leftKey = left.kind === 'starter_shop' ? 'starter_shop' : left.colorKey ?? left.function ?? left.kind
    const rightKey = right.kind === 'starter_shop' ? 'starter_shop' : right.colorKey ?? right.function ?? right.kind
    return leftKey.localeCompare(rightKey) || left.id.localeCompare(right.id)
  })
  const iconCounts = acquiredCards.flatMap((card) => [...(card.kind === 'marketing' ? Array.from({ length: card.brandAwareness ?? 1 }, () => 'marketing') : (card.icons ?? [])), ...(card.kind === 'starter_shop' ? ['channel'] : [])]).reduce<Record<string, number>>((counts, icon) => ({ ...counts, [icon]: (counts[icon] ?? 0) + 1 }), {})
  const customerCounts = acquiredCards.reduce<Record<string, number>>((counts, card) => {
    const source = card.customerCount
    for (const type of ['gourmet', 'regular']) counts[type] = (counts[type] ?? 0) + (source?.[type] ?? 0)
    return counts
  }, {})
  const selected = hand.find((card) => card.id === selectedCard)
  const phaseLabel: Record<string, string> = { hypothesis: '設定關鍵指標', experiment: '選擇實驗策略', learning: '結算本期市場', finished: '遊戲完成' }
  const totalCards = startingCards.length + hand.length + tableau.length
  const market = useMemo(() => Object.entries(room.demandBoard ?? {}), [room.demandBoard])

  return <Box sx={{ minHeight: '100vh', bgcolor: '#f8f5f1', color: '#33251f' }}>
    <Box sx={{ bgcolor: '#4f3428', color: 'white' }}><Container maxWidth="xl"><Toolbar disableGutters sx={{ justifyContent: 'space-between' }}><Stack direction="row" spacing={1.5} alignItems="center"><Coffee /><Typography variant="h6" fontWeight={900}>Cafe Startups</Typography></Stack><Stack direction="row" spacing={1}><Chip label={`第 ${room.period} 期`} sx={{ color: 'white', bgcolor: '#ffffff20' }} /><Button color="inherit" size="small" onClick={leave}>離開</Button></Stack></Toolbar></Container></Box>
    <Container maxWidth="xl" sx={{ py: 3 }}>
      <Stack direction={{ xs: 'column', md: 'row' }} justifyContent="space-between" alignItems={{ xs: 'start', md: 'center' }} spacing={1} sx={{ mb: 2.5 }}><Box><Typography variant="overline" color="text.secondary">PERIOD {room.period} · ROUND {room.round}/6</Typography><Typography variant="h4" fontWeight={900}>{phaseLabel[room.phase] ?? room.phase}</Typography></Box><Chip icon={<Coffee />} label={`${totalCards} 張卡牌`} color="primary" variant="outlined" /></Stack>
      {room.phase === 'hypothesis' && room.period === 0 && <GameDashboardSetupCards partners={room.partnerOptions ?? []} shops={room.starterShopOptions ?? []} selectedPartnerId={me?.initialCardsSelected ? me?.partner?.id : undefined} selectedShopId={me?.initialCardsSelected ? me?.starterShop?.id : undefined} busy={busy} onSelect={setupInitialCards} onBegin={() => command('BEGIN_EXPERIMENT')} />}
      <KPISelectionPanel key={room.period} room={room} command={command} busy={busy} />
      <CustomerCountsPanel counts={customerCounts} />
      <PlayerSummaryBar room={room} />
      <MarketRankingPanel room={room} command={command} busy={busy} />
      <Paper elevation={0} sx={{ p: { xs: 2, md: 2.5 }, border: '1px solid', borderColor: 'divider', borderRadius: 3, mb: 2.5 }}><Stack direction={{ xs: 'column', md: 'row' }} justifyContent="space-between" spacing={1}><Box><Typography variant="overline" color="primary">YOUR COLLECTION</Typography><Typography variant="h5" fontWeight={900}>我的卡牌</Typography></Box><Typography variant="body2" color="text.secondary" sx={{ maxWidth: 420 }}>開局取得的卡牌、目前手牌與已打出的卡牌都會顯示在這裡。</Typography></Stack><Grid container spacing={1.2} sx={{ mt: 1.5 }}>{Object.entries(iconCounts).map(([icon, count]) => <Grid key={icon} size={{ xs: 4, sm: 2, md: 1.5 }}><Paper variant="outlined" sx={{ p: 1.1, bgcolor: '#faf7f2', textAlign: 'center' }}><GameIcon name={icon} sx={{ fontSize: 24 }} /><Typography display="block" variant="h6" fontWeight={900}>{count}</Typography></Paper></Grid>)}</Grid><GameDashboardCardGroup title="已取得卡牌" subtitle={`${acquiredCards.length} 張・開局卡與已打出卡牌`} cards={orderedAcquiredCards} collapsible /><GameDashboardCardGroup title="目前手牌" subtitle={`${hand.length} 張・點選卡牌後執行操作`} cards={hand} selectedId={selectedCard} onSelect={(id) => { setSelectedCard(id); command('SELECT_CARD', { cardId: id }) }} headerContent={<HandActionBar room={room} selected={selected} ownedIcons={ownedIcons} command={command} busy={busy} />} /></Paper>
      <Grid container spacing={2.5}>
        <Grid size={{ xs: 12, md: 3 }}><Stack spacing={2}><Paper sx={{ p: 2.2, borderRadius: 3 }}><Typography variant="overline" color="text.secondary">我的咖啡館</Typography><Typography variant="h3" fontWeight={900} sx={{ my: 1 }}>${me?.cash ?? 0}<Typography component="span" variant="body2" color="text.secondary"> 萬</Typography></Typography><Divider sx={{ mb: 1 }} />{[['貸款', `${me?.loans ?? 0} / 6`], ['營收', `$${me?.revenue ?? 0} 萬`], ['分數', `${me?.score ?? '—'}`], ['品牌 / 產品', `${me?.brandAwareness ?? 0} / ${me?.products ?? 0}`], ['價值 / 資源', `${me?.values ?? 0} / ${me?.resources ?? 0}`]].map(([label, value]) => <Stack key={label} direction="row" justifyContent="space-between" sx={{ py: .8 }}><Typography variant="body2" color="text.secondary">{label}</Typography><Typography variant="body2" fontWeight={800}>{value}</Typography></Stack>)}</Paper><Button fullWidth variant="outlined" startIcon={<Paid />} onClick={() => command('TAKE_LOAN')} disabled={busy}>取得一筆貸款</Button></Stack></Grid>
        <Grid size={{ xs: 12, md: 9 }}><Stack spacing={2}><Paper sx={{ p: 2.2, borderRadius: 3, bgcolor: '#714b38', color: 'white' }}><Stack direction={{ xs: 'column', sm: 'row' }} justifyContent="space-between" spacing={1}><Box><Typography variant="overline" sx={{ opacity: .7 }}>目前階段</Typography><Typography variant="h6" fontWeight={900}>{phaseLabel[room.phase] ?? room.phase}</Typography></Box><Typography variant="body2" sx={{ opacity: .75 }}>{room.players.filter((player) => !player.handCount).length ? '等待玩家行動' : '系統同步中'}</Typography></Stack></Paper><Paper sx={{ p: 2.2, borderRadius: 3 }}><Typography variant="h6" fontWeight={900}>需求市場</Typography><Grid container spacing={1.5} sx={{ mt: .5 }}>{market.length ? market.map(([key, value]) => <Grid key={key} size={{ xs: 6, sm: 3 }}><Paper variant="outlined" sx={{ p: 1.5, bgcolor: '#faf7f2' }}><Typography variant="caption" color="text.secondary">{key}</Typography><Typography variant="h6" fontWeight={900}>{value}</Typography></Paper></Grid>) : <Grid size={12}><Typography variant="body2" color="text.secondary">市場資料將在本期開始後顯示。</Typography></Grid>}</Grid></Paper>{room.phase === 'experiment' && <Stack direction={{ xs: 'column', sm: 'row' }} justifyContent="space-between" spacing={1}><Button variant="outlined" onClick={() => selected && command('DISCARD_SELECTED_CARD')} disabled={busy || !selected}>棄牌 +20</Button><Button variant="contained" onClick={() => selected && command('PLAY_SELECTED_CARD')} disabled={busy || !selected}>打出 {selected ? `$${selected.cost.cash} 萬` : '選擇一張牌'}</Button><Button variant="contained" color="secondary" onClick={() => command('PASS_HAND')} disabled={busy}>完成選擇並傳牌 →</Button></Stack>}{room.phase === 'learning' && <Button variant="contained" onClick={() => command('RESOLVE_LEARNING')} disabled={busy}>結算本期並繼續</Button>}{error && <Alert severity="error">{error}</Alert>}</Stack></Grid>
      </Grid>
      <Paper sx={{ mt: 2.5, p: 2, borderRadius: 3 }}><Typography variant="subtitle1" fontWeight={900} sx={{ mb: 1 }}>玩家</Typography><Stack direction="row" flexWrap="wrap" gap={1}>{room.players.map((player) => <Chip key={player.id} icon={<Groups />} label={`${player.displayName}${player.id === me?.id ? '（你）' : player.bot ? '（電腦）' : ''} · ${player.handCount} 張`} variant={player.id === me?.id ? 'filled' : 'outlined'} color={player.id === me?.id ? 'primary' : 'default'} />)}</Stack></Paper>
    </Container>
  </Box>
}
