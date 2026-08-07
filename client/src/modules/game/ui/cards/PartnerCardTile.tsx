// @ts-nocheck
import { Box, Card as MuiCard, Typography } from '@mui/material'
import type { PlayerCard } from '../../model/cardTypes'
import { dashboardCardColors } from './gameDashboardCardTheme'
import { CardArtwork } from './CardArtwork'
import { CustomerTypeTokens } from './CustomerTypeTokens'
import { GameIcon } from './GameIcon'
import { CardCost } from './CardCost'

export function PartnerCardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const functionKey = card.colorKey ?? card.function
  const colorTheme = dashboardCardColors.management[functionKey ?? '']
  if (!colorTheme) return null

  const customerEntries = Object.entries(card.customerCount ?? {}).filter(([key, value]) => (key === 'gourmet' || key === 'regular') && value > 0)
  const icon = card.icons[0] ?? 'people'

  return <MuiCard onClick={onClick} variant="outlined" sx={{ display: 'flex', height: '100%', minHeight: 290, flexDirection: 'column', overflow: 'hidden', cursor: onClick ? 'pointer' : 'default', bgcolor: colorTheme.pale, borderColor: selected ? colorTheme.color : `${colorTheme.color}66`, borderWidth: selected ? 3 : 1 }}>
    <Box sx={{ minHeight: 39, px: 1.5, py: .9, boxSizing: 'border-box', bgcolor: colorTheme.color, color: 'white', textAlign: 'center' }}><Typography variant="caption" fontWeight={900}>{card.name}</Typography></Box>
    <Box sx={{ display: 'flex', height: 88, flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: .2, px: 1.5, py: .85, boxSizing: 'border-box', bgcolor: `${colorTheme.color}dd`, color: 'white' }}>
      <Typography variant="caption" sx={{ opacity: .8 }}>卡片功能</Typography>
      {customerEntries.length > 0
        ? <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: .8 }}>
            {customerEntries.map(([key, value]) => <CustomerTypeTokens key={key} type={key} count={value} size={24} />)}
          </Box>
        : <GameIcon name={icon} sx={{ fontSize: 30 }} />}
    </Box>
    <Box sx={{ display: 'grid', height: 180, flex: '0 0 180px', placeItems: 'center', overflow: 'hidden' }}><CardArtwork card={card} /></Box>
    <Box sx={{ minHeight: 88, flex: 1, px: 1.5, py: 1, boxSizing: 'border-box', bgcolor: '#ffffffcc' }}><Typography variant="caption" color="text.secondary">卡片說明</Typography><Typography variant="body2" sx={{ minHeight: 34 }}>{card.description ?? '合夥人卡'}</Typography></Box>
    <Box sx={{ minHeight: 42, display: 'flex', alignItems: 'center', px: 1.5, py: .8, boxSizing: 'border-box', bgcolor: colorTheme.pale }}><CardCost card={card} color={colorTheme.color} /></Box>
  </MuiCard>
}
