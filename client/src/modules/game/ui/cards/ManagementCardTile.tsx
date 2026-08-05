// @ts-nocheck
import { Box, Card as MuiCard, Stack, Typography } from '@mui/material'
import type { PlayerCard } from '../../model/cardTypes'
import { dashboardCardColors } from './gameDashboardCardTheme'
import { CardArtwork } from './CardArtwork'
import { GameIcon } from './GameIcon'
import { CustomerTypeTokens } from './CustomerTypeTokens'
import { CardCost } from './CardCost'

export function ManagementCardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const meta = dashboardCardColors.management[card.function ?? card.kind]
  if (!meta) return null
  const marketEntries = Object.entries(card.marketChange ?? {}).filter(([key, value]) => (key === 'gourmet' || key === 'regular') && value !== 0)
  const customerEntries = Object.entries(card.customerCount ?? {}).filter(([key, value]) => (key === 'gourmet' || key === 'regular') && value !== 0)
  const visibleCustomerEntries = customerEntries.length > 0 ? customerEntries : marketEntries.length > 0 ? marketEntries : [['regular', 1]]
  const visibleMarketEntries = card.kind === 'channel' ? visibleCustomerEntries : marketEntries

  return <MuiCard onClick={onClick} variant="outlined" sx={{ display: 'flex', height: '100%', minHeight: 290, flexDirection: 'column', overflow: 'hidden', cursor: onClick ? 'pointer' : 'default', bgcolor: meta.pale, borderColor: selected ? meta.color : `${meta.color}66`, borderWidth: selected ? 3 : 1 }}>
    <Box sx={{ minHeight: 39, px: 1.5, py: .9, boxSizing: 'border-box', bgcolor: meta.color, color: 'white', textAlign: 'center' }}><Typography variant="caption" fontWeight={900}>{card.name}</Typography></Box>
    <Box sx={{ display: 'flex', height: 88, flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: .2, px: 1.5, py: .85, boxSizing: 'border-box', bgcolor: `${meta.color}dd`, color: 'white' }}>
      <Typography variant="caption" sx={{ opacity: .8 }}>卡片功能</Typography>
      <Box sx={{ display: 'flex', minWidth: 0, maxWidth: '100%', alignItems: 'center', justifyContent: 'center', gap: .8, overflow: 'hidden' }}>
        {card.kind === 'channel' ? visibleMarketEntries.map(([key, value]) => <CustomerTypeTokens key={key} type={key} count={value} size={18} />) : card.kind === 'resource' ? ['data', 'procurement', 'operations', 'marketing_resource'].map((icon) => <GameIcon key={icon} name={icon} sx={{ fontSize: 24 }} />) : card.icons.slice(0, 4).map((icon, index) => <GameIcon key={`${icon}-${index}`} name={icon} sx={{ fontSize: 24 }} />)}
      </Box>
    </Box>
    <Box sx={{ display: 'grid', height: 105, flex: '0 0 105px', placeItems: 'center', overflow: 'hidden' }}><CardArtwork card={card} /></Box>
    <Box sx={{ minHeight: 88, flex: 1, px: 1.5, py: 1, boxSizing: 'border-box', bgcolor: '#ffffffcc' }}><Typography variant="caption" color="text.secondary">卡片說明</Typography><Typography variant="body2" sx={{ minHeight: 34 }}>{card.description ?? '可執行的經營管理行動。'}</Typography></Box>
    <Stack direction="column" sx={{ minHeight: 62, bgcolor: meta.pale }}>
      <Box sx={{ display: 'flex', width: '100%', minHeight: 34, alignItems: 'center', px: 1.5, py: .65, borderBottom: '1px solid rgba(111,82,65,.16)' }}><CardCost card={card} color={meta.color} /></Box>
      <Box sx={{ display: 'flex', minHeight: 28, alignItems: 'center', justifyContent: 'flex-end', gap: .25, px: 1.5, py: .45 }}>{visibleMarketEntries.map(([key, value]) => <CustomerTypeTokens key={key} type={key} count={value} size={14} />)}</Box>
    </Stack>
  </MuiCard>
}
