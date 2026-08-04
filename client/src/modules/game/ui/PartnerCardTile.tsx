// @ts-nocheck
import { Coffee, Groups, Psychology, Settings } from '@mui/icons-material'
import { Box, Card as MuiCard, Stack, Typography } from '@mui/material'
import type { PlayerCard } from './gameDashboardCardTypes'
import { dashboardCardColors } from './gameDashboardCardTheme'

const partnerMeta: Record<string, { icon: typeof Coffee; role: string }> = {
  'partner-barista': { icon: Coffee, role: '咖啡師・營運專業' },
  'partner-roaster': { icon: Coffee, role: '烘豆師・產品專業' },
  'partner-marketer': { icon: Psychology, role: '行銷師・品牌專業' },
  'partner-service': { icon: Groups, role: '服務設計・體驗專業' },
}

export function PartnerCardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const colorTheme = dashboardCardColors.partner[card.id]
  const meta = { ...(partnerMeta[card.id] ?? { icon: Groups, role: '創業夥伴' }), ...colorTheme }
  if (!colorTheme) return null
  const Icon = meta.icon
  return <MuiCard onClick={onClick} variant="outlined" sx={{ height: '100%', minHeight: 290, overflow: 'hidden', cursor: onClick ? 'pointer' : 'default', bgcolor: meta.pale, borderColor: selected ? meta.color : `${meta.color}66`, borderWidth: selected ? 3 : 1 }}><Box sx={{ px: 1.5, py: .9, bgcolor: meta.color, color: 'white', textAlign: 'center' }}><Typography variant="caption" fontWeight={900}>{card.name}</Typography></Box><Box sx={{ px: 1.5, py: .85, display: 'flex', alignItems: 'center', gap: 1, bgcolor: `${meta.color}dd`, color: 'white' }}><Settings sx={{ fontSize: 19 }} /><Box><Typography variant="caption" sx={{ display: 'block', opacity: .8 }}>卡片功能</Typography><Typography variant="body2" fontWeight={900}>{meta.role}</Typography></Box></Box><Box sx={{ display: 'grid', placeItems: 'center', minHeight: 105, color: meta.color }}><Icon sx={{ fontSize: 64 }} /></Box><Box sx={{ px: 1.5, py: 1, bgcolor: '#ffffffcc' }}><Typography variant="caption" color="text.secondary">卡片說明</Typography><Typography variant="body2" sx={{ minHeight: 34 }}>{card.description ?? '創業合夥人提供的特殊能力。'}</Typography><Typography variant="caption" fontWeight={900} color={meta.color}>{card.effect ?? '—'}</Typography></Box><Stack direction="row" justifyContent="space-between" sx={{ px: 1.5, py: .8, bgcolor: meta.pale }}><Typography variant="caption">成本</Typography><Typography variant="body2" fontWeight={900} color={meta.color}>${card.cost?.cash ?? 0} 萬</Typography></Stack></MuiCard>
}
