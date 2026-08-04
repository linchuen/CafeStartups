// MUI 9 exposes several system props through `sx`; this screen is kept isolated
// while the project migrates its remaining legacy JSX to the same API.
// @ts-nocheck
import { useMemo } from 'react'
import { Coffee, Groups, Storefront, Paid, Psychology } from '@mui/icons-material'
import { Alert, Box, Button, Card as MuiCard, CardContent, Chip, Container, Divider, Grid, LinearProgress, Paper, Stack, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Toolbar, Typography } from '@mui/material'
import type { CardFaceData } from '../../../CardFace'

type PlayerCard = CardFaceData
type DashboardPlayer = { id: string; displayName: string; bot?: boolean; handCount: number; cash: number; loans: number; revenue?: number; score?: number; brandAwareness?: number; products?: number; values?: number; resources?: number }
type CashFlowStatement = { period: number; beginningCash: number; operatingRevenue: number; otherIncome: number; operatingExpenses: number; interestPaid: number; principalRepayment: number; newLoans: number; endingCash: number }
type DashboardRoom = { period: number; round: number; phase: string; players: DashboardPlayer[]; demandBoard?: Record<string, number>; me?: { id: string; hand: PlayerCard[] | null; tableau: PlayerCard[] | null; partner?: PlayerCard; starterShop?: PlayerCard; cash: number; loans: number; revenue?: number; score?: number; selectedKPIs?: string[]; cashFlow?: CashFlowStatement[] | null; cashFlowRounds?: CashFlowStatement[] | null; brandAwareness?: number; products?: number; values?: number; resources?: number } }

const kindMeta: Record<string, { label: string; color: string; icon: typeof Coffee }> = {
  partner: { label: '合夥人卡', color: '#714b38', icon: Groups },
  starter_shop: { label: '創始店卡', color: '#287477', icon: Storefront },
  resource: { label: '資源', color: '#3976a6', icon: Coffee },
  product: { label: '產品', color: '#b98a25', icon: Coffee },
  value: { label: '價值', color: '#b44f52', icon: Psychology },
  channel: { label: '通路', color: '#3f7d66', icon: Storefront },
  marketing: { label: '行銷', color: '#7a5ba5', icon: Psychology },
}

function CardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const meta = kindMeta[card.kind] ?? { label: card.kind, color: '#765341', icon: Coffee }
  const Icon = meta.icon
  return <MuiCard onClick={onClick} variant="outlined" sx={{ height: '100%', minHeight: 172, cursor: onClick ? 'pointer' : 'default', borderColor: selected ? meta.color : 'divider', borderWidth: selected ? 2 : 1, transition: 'transform .15s, box-shadow .15s', '&:hover': onClick ? { transform: 'translateY(-3px)', boxShadow: 3 } : undefined }}>
    <Box sx={{ height: 6, bgcolor: meta.color }} />
    <CardContent sx={{ display: 'flex', height: 'calc(100% - 6px)', flexDirection: 'column', gap: 1.1 }}>
      <Stack direction="row" justifyContent="space-between" alignItems="center"><Chip size="small" label={meta.label} sx={{ color: meta.color, bgcolor: `${meta.color}18`, fontWeight: 700 }} /><Typography variant="caption" color="text.secondary">第 {card.period} 期</Typography></Stack>
      <Box sx={{ display: 'grid', placeItems: 'center', height: 44, borderRadius: 2, color: meta.color, bgcolor: `${meta.color}12` }}><Icon /></Box>
      <Typography variant="subtitle1" fontWeight={800}>{card.name}</Typography>
      <Typography variant="caption" color="text.secondary" sx={{ minHeight: 32 }}>{card.description ?? '可用於咖啡館經營的卡牌。'}</Typography>
      <Stack direction="row" justifyContent="space-between" sx={{ mt: 'auto' }}><Typography variant="caption" color="text.secondary">成本</Typography><Typography variant="body2" fontWeight={800} color={meta.color}>${card.cost.cash}</Typography></Stack>
    </CardContent>
  </MuiCard>
}

function CardGroup({ title, subtitle, cards, selectedId, onSelect }: { title: string; subtitle: string; cards: PlayerCard[]; selectedId?: string; onSelect?: (id: string) => void }) {
  if (!cards.length) return null
  return <Box sx={{ mt: 2.5 }}><Stack direction="row" alignItems="baseline" spacing={1} sx={{ mb: 1.2 }}><Typography variant="subtitle1" fontWeight={800}>{title}</Typography><Typography variant="caption" color="text.secondary">{subtitle}</Typography></Stack><Grid container spacing={1.5}>{cards.map((card) => <Grid key={card.id} size={{ xs: 12, sm: 6, md: 3, lg: 2 }}><CardTile card={card} selected={selectedId === card.id} onClick={onSelect ? () => onSelect(card.id) : undefined} /></Grid>)}</Grid></Box>
}

function GamePanel({ room, market }: { room: DashboardRoom; market: [string, number][] }) {
  const metrics = [['品牌知名度', room.me?.brandAwareness ?? 0], ['產品能力', room.me?.products ?? 0], ['價值主張', room.me?.values ?? 0], ['資源能力', room.me?.resources ?? 0]]
  return <Paper sx={{ p: { xs: 2, md: 2.5 }, mb: 2.5, borderRadius: 3, border: '1px solid', borderColor: 'divider' }}><Stack direction={{ xs: 'column', md: 'row' }} justifyContent="space-between" spacing={2}><Box><Typography variant="overline" color="primary">GAME BOARD</Typography><Typography variant="h5" fontWeight={900}>遊戲面板</Typography><Typography variant="body2" color="text.secondary">第 {room.period} 期・第 {room.round} 回合・{room.phase}</Typography></Box><Stack direction="row" spacing={1} alignItems="center" sx={{ minWidth: { md: 280 } }}>{[1, 2, 3].map((period) => <Box key={period} sx={{ flex: 1 }}><Stack direction="row" justifyContent="space-between"><Typography variant="caption">第 {period} 期</Typography><Typography variant="caption" color={room.period >= period ? 'primary' : 'text.disabled'}>{room.period >= period ? '進行中' : '未開始'}</Typography></Stack><LinearProgress variant="determinate" value={room.period >= period ? 100 : 0} color={room.period === period ? 'primary' : 'secondary'} /></Box>)}</Stack></Stack><Divider sx={{ my: 2 }} /><Grid container spacing={1.5}>{metrics.map(([label, value]) => <Grid key={label} size={{ xs: 6, sm: 3 }}><Paper variant="outlined" sx={{ p: 1.4, bgcolor: '#faf7f2' }}><Typography variant="caption" color="text.secondary">{label}</Typography><Typography variant="h6" fontWeight={900}>{value}</Typography></Paper></Grid>)}{market.map(([key, value]) => <Grid key={key} size={{ xs: 6, sm: 3 }}><Paper variant="outlined" sx={{ p: 1.4, bgcolor: '#f1f8f7' }}><Typography variant="caption" color="text.secondary">市場・{key}</Typography><Typography variant="h6" fontWeight={900}>{value}</Typography></Paper></Grid>)}</Grid></Paper>
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
      <GamePanel room={room} market={market} />
      <Paper elevation={0} sx={{ p: { xs: 2, md: 2.5 }, border: '1px solid', borderColor: 'divider', borderRadius: 3, mb: 2.5 }}><Stack direction={{ xs: 'column', md: 'row' }} justifyContent="space-between" spacing={1}><Box><Typography variant="overline" color="primary">YOUR COLLECTION</Typography><Typography variant="h5" fontWeight={900}>我的卡牌</Typography></Box><Typography variant="body2" color="text.secondary" sx={{ maxWidth: 420 }}>開局取得的合夥人卡、創始店卡，以及目前手牌與已打出的卡牌都會顯示在這裡。</Typography></Stack><CardGroup title="開局卡牌" subtitle="合夥人卡・創始店卡" cards={startingCards} /><CardGroup title="目前手牌" subtitle={`${hand.length} 張・點選卡牌後執行操作`} cards={hand} selectedId={selectedCard} onSelect={(id) => { setSelectedCard(id); command('SELECT_CARD', { cardId: id }) }} /><CardGroup title="已打出卡牌" subtitle={`${tableau.length} 張`} cards={tableau} /></Paper>
      <Grid container spacing={2.5}>
        <Grid size={{ xs: 12, md: 3 }}><Stack spacing={2}><Paper sx={{ p: 2.2, borderRadius: 3 }}><Typography variant="overline" color="text.secondary">我的咖啡館</Typography><Typography variant="h3" fontWeight={900} sx={{ my: 1 }}>${me?.cash ?? 0}<Typography component="span" variant="body2" color="text.secondary"> 萬</Typography></Typography><Divider sx={{ mb: 1 }} />{[['貸款', `${me?.loans ?? 0} / 6`], ['營收', `$${me?.revenue ?? 0} 萬`], ['分數', `${me?.score ?? '—'}`], ['品牌 / 產品', `${me?.brandAwareness ?? 0} / ${me?.products ?? 0}`], ['價值 / 資源', `${me?.values ?? 0} / ${me?.resources ?? 0}`]].map(([label, value]) => <Stack key={label} direction="row" justifyContent="space-between" sx={{ py: .8 }}><Typography variant="body2" color="text.secondary">{label}</Typography><Typography variant="body2" fontWeight={800}>{value}</Typography></Stack>)}</Paper><Button fullWidth variant="outlined" startIcon={<Paid />} onClick={() => command('TAKE_LOAN')} disabled={busy}>取得一筆貸款</Button></Stack></Grid>
        <Grid size={{ xs: 12, md: 9 }}><Stack spacing={2}><Paper sx={{ p: 2.2, borderRadius: 3, bgcolor: '#714b38', color: 'white' }}><Stack direction={{ xs: 'column', sm: 'row' }} justifyContent="space-between" spacing={1}><Box><Typography variant="overline" sx={{ opacity: .7 }}>目前階段</Typography><Typography variant="h6" fontWeight={900}>{phaseLabel[room.phase] ?? room.phase}</Typography></Box><Typography variant="body2" sx={{ opacity: .75 }}>{room.players.filter((player) => !player.handCount).length ? '等待玩家行動' : '系統同步中'}</Typography></Stack></Paper><Paper sx={{ p: 2.2, borderRadius: 3 }}><Typography variant="h6" fontWeight={900}>需求市場</Typography><Grid container spacing={1.5} sx={{ mt: .5 }}>{market.length ? market.map(([key, value]) => <Grid key={key} size={{ xs: 6, sm: 3 }}><Paper variant="outlined" sx={{ p: 1.5, bgcolor: '#faf7f2' }}><Typography variant="caption" color="text.secondary">{key}</Typography><Typography variant="h6" fontWeight={900}>{value}</Typography></Paper></Grid>) : <Grid size={12}><Typography variant="body2" color="text.secondary">市場資料將在本期開始後顯示。</Typography></Grid>}</Grid></Paper>{room.phase === 'experiment' && <Stack direction={{ xs: 'column', sm: 'row' }} justifyContent="space-between" spacing={1}><Button variant="outlined" onClick={() => selected && command('DISCARD_SELECTED_CARD')} disabled={busy || !selected}>棄牌 +20</Button><Button variant="contained" onClick={() => selected && command('PLAY_SELECTED_CARD')} disabled={busy || !selected}>打出 {selected ? `$${selected.cost.cash} 萬` : '選擇一張牌'}</Button><Button variant="contained" color="secondary" onClick={() => command('PASS_HAND')} disabled={busy}>完成選擇並傳牌 →</Button></Stack>}{room.phase === 'learning' && <Button variant="contained" onClick={() => command('RESOLVE_LEARNING')} disabled={busy}>結算本期並繼續</Button>}{error && <Alert severity="error">{error}</Alert>}</Stack></Grid>
      </Grid>
      <Paper sx={{ mt: 2.5, p: 2, borderRadius: 3 }}><Typography variant="subtitle1" fontWeight={900} sx={{ mb: 1 }}>玩家</Typography><Stack direction="row" flexWrap="wrap" gap={1}>{room.players.map((player) => <Chip key={player.id} icon={<Groups />} label={`${player.displayName}${player.id === me?.id ? '（你）' : player.bot ? '（電腦）' : ''} · ${player.handCount} 張`} variant={player.id === me?.id ? 'filled' : 'outlined'} color={player.id === me?.id ? 'primary' : 'default'} />)}</Stack></Paper>
      <CashFlowPanel room={room} />
    </Container>
  </Box>
}
