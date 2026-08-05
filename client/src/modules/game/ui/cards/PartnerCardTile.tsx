// @ts-nocheck
import { Box, Card as MuiCard, Stack, Typography } from '@mui/material'
import type { PlayerCard } from '../../model/cardTypes'
import { dashboardCardColors, partnerFunctionById } from './gameDashboardCardTheme'
import { CardArtwork } from './CardArtwork'
import { GameIcon } from './GameIcon'

const partnerMeta: Record<string, { icon: string; role: string }> = {
  'partner-barista': { icon: 'operations', role: '營運功能' },
  'partner-roaster': { icon: 'coffee', role: '產品功能' },
  'partner-marketer': { icon: 'marketing', role: '行銷功能' },
  'partner-service': { icon: 'people', role: '價值功能' },
  'partner-finance': { icon: 'people', role: '資源功能' },
  'partner-pastry': { icon: 'coffee', role: '產品功能' },
  'partner-supply': { icon: 'people', role: '資源功能' },
  'partner-community': { icon: 'people', role: '通路功能' },
  'partner-hr': { icon: 'people', role: '價值功能' },
  'partner-analytics': { icon: 'marketing', role: '通路功能' },
}

export function PartnerCardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const functionKey = card.function ?? partnerFunctionById[card.id]
  const colorTheme = dashboardCardColors.management[functionKey ?? '']
  const meta = { ...(partnerMeta[card.id] ?? { icon: 'people', role: '合夥人功能' }), ...colorTheme }
  if (!colorTheme) return null
  return <MuiCard onClick={onClick} variant="outlined" sx={{ display: 'flex', height: '100%', minHeight: 290, flexDirection: 'column', overflow: 'hidden', cursor: onClick ? 'pointer' : 'default', bgcolor: meta.pale, borderColor: selected ? meta.color : `${meta.color}66`, borderWidth: selected ? 3 : 1 }}><Box sx={{ minHeight: 39, px: 1.5, py: .9, boxSizing: 'border-box', bgcolor: meta.color, color: 'white', textAlign: 'center' }}><Typography variant="caption" fontWeight={900}>{card.name}</Typography></Box><Box sx={{ display: 'flex', height: 88, flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: .2, px: 1.5, py: .85, boxSizing: 'border-box', bgcolor: `${meta.color}dd`, color: 'white' }}><Typography variant="caption" sx={{ opacity: .8 }}>{meta.role}</Typography><GameIcon name={meta.icon} sx={{ fontSize: 26 }} /></Box><Box sx={{ display: 'grid', height: 105, flex: '0 0 105px', placeItems: 'center', overflow: 'hidden' }}><CardArtwork card={card} /></Box><Box sx={{ minHeight: 88, flex: 1, px: 1.5, py: 1, boxSizing: 'border-box', bgcolor: '#ffffffcc' }}><Typography variant="caption" color="text.secondary">卡片說明</Typography><Typography variant="body2" sx={{ minHeight: 34 }}>{card.description ?? '合夥人卡'}</Typography></Box><Stack sx={{ minHeight: 42, bgcolor: meta.pale }} /></MuiCard>
}
