// @ts-nocheck
import { Box, Stack, Typography } from '@mui/material'
import type { PlayerCard } from '../../model/cardTypes'
import { GameIcon } from './GameIcon'
import { cardColorTokens } from './gameDashboardCardTheme'

const costIconColors: Record<string, { color: string; background: string }> = {
  data: { color: '#2d6897', background: '#d5e7ef' },
  marketing_resource: { color: '#2d6897', background: '#d5e7ef' },
  operations: { color: '#2d6897', background: '#d5e7ef' },
  beans: { color: '#c88d28', background: '#f3dfb7' },
  coffee: { color: '#c88d28', background: '#f3dfb7' },
  dessert: { color: '#c88d28', background: '#f3dfb7' },
  service: { color: cardColorTokens.value.color, background: cardColorTokens.value.pale },
  taste: { color: cardColorTokens.value.color, background: cardColorTokens.value.pale },
  value: { color: cardColorTokens.value.color, background: cardColorTokens.value.pale },
  channel: { color: '#3f7d66', background: '#e2f0e8' },
  marketing: { color: '#7a5ba5', background: '#eee7f7' },
  procurement: { color: '#2d6897', background: '#d5e7ef' },
}

export function CardCost({ card, color }: { card: PlayerCard; color?: string }) {
  const icons = card.cost?.icons ?? []
  const accent = color ?? '#b36f42'
  const isPartnerInitialCashBonus = card.kind === 'partner' && card.function === 'channel' && (card.cost?.cash ?? 0) > 0
  const cashLabel = `${isPartnerInitialCashBonus ? '+' : ''}$${card.cost?.cash ?? 0} 萬`

  return <Stack direction="row" spacing={.55} alignItems="center" sx={{ minWidth: 0 }}>
    <Typography variant="caption" color="text.secondary">成本</Typography>
    <Typography variant="body2" fontWeight={900} sx={{ color: accent }}>{cashLabel}</Typography>
    {icons.length > 0 && <Box component="span" role="img" aria-label={`資源需求 ${icons.length} 個`} sx={{ display: 'inline-flex', alignItems: 'center', gap: .3, ml: .25 }}>
      {icons.map((icon, index) => {
        const iconColor = costIconColors[icon] ?? { color: accent, background: `${accent}22` }
        return <Box key={`${icon}-${index}`} component="span" sx={{ display: 'grid', placeItems: 'center', width: 18, height: 18, border: `1px solid ${iconColor.color}55`, borderRadius: 1, color: iconColor.color, backgroundColor: iconColor.background }}><GameIcon name={icon} sx={{ fontSize: 12 }} /></Box>
      })}
    </Box>}
  </Stack>
}
