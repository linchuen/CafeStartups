// @ts-nocheck
import { Coffee, Groups, Psychology } from '@mui/icons-material'
import { Box, Card as MuiCard, Stack, Typography } from '@mui/material'
import type { PlayerCard } from './gameDashboardCardTypes'
import { dashboardCardColors, partnerFunctionById } from './gameDashboardCardTheme'
import { CardArtwork } from './CardArtwork'

const partnerMeta: Record<string, { icon: typeof Coffee; role: string }> = {
  'partner-barista': { icon: Coffee, role: '咖啡師・營運專業' },
  'partner-roaster': { icon: Coffee, role: '烘豆師・產品專業' },
  'partner-marketer': { icon: Psychology, role: '行銷師・品牌專業' },
  'partner-service': { icon: Groups, role: '服務設計・體驗專業' },
  'partner-finance': { icon: Groups, role: '財務管理・資金專業' },
  'partner-pastry': { icon: Coffee, role: '甜點開發・產品專業' },
  'partner-supply': { icon: Groups, role: '供應鏈管理・採購專業' },
  'partner-community': { icon: Groups, role: '社區合作・通路專業' },
  'partner-hr': { icon: Groups, role: '人才培訓・服務專業' },
  'partner-analytics': { icon: Psychology, role: '數據營運・通路專業' },
}

export function PartnerCardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const functionKey = card.colorKey ?? card.function ?? partnerFunctionById[card.id]
  const colorTheme = dashboardCardColors.management[functionKey ?? '']
  const meta = { ...(partnerMeta[card.id] ?? { icon: Groups, role: '創業夥伴' }), ...colorTheme }
  if (!colorTheme) return null
  const Icon = meta.icon
  return <MuiCard onClick={onClick} variant="outlined" sx={{ display: 'flex', height: '100%', minHeight: 290, flexDirection: 'column', overflow: 'hidden', cursor: onClick ? 'pointer' : 'default', bgcolor: meta.pale, borderColor: selected ? meta.color : `${meta.color}66`, borderWidth: selected ? 3 : 1 }}><Box sx={{ minHeight: 39, px: 1.5, py: .9, boxSizing: 'border-box', bgcolor: meta.color, color: 'white', textAlign: 'center' }}><Typography variant="caption" fontWeight={900}>{card.name}</Typography></Box><Box sx={{ display: 'flex', height: 88, flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: .2, px: 1.5, py: .85, boxSizing: 'border-box', bgcolor: `${meta.color}dd`, color: 'white' }}><Typography variant="caption" sx={{ opacity: .8 }}>卡片功能</Typography><Icon sx={{ fontSize: 26 }} /></Box><Box sx={{ display: 'grid', height: 105, flex: '0 0 105px', placeItems: 'center', overflow: 'hidden' }}><CardArtwork card={card} /></Box><Box sx={{ minHeight: 88, flex: 1, px: 1.5, py: 1, boxSizing: 'border-box', bgcolor: '#ffffffcc' }}><Typography variant="caption" color="text.secondary">卡片說明</Typography><Typography variant="body2" sx={{ minHeight: 34 }}>{card.description ?? '創業合夥人提供的特殊能力。'}</Typography><Typography variant="caption" fontWeight={900} color={meta.color}>{card.effect ?? '—'}</Typography></Box><Stack direction="row" justifyContent="space-between" sx={{ minHeight: 42, alignItems: 'center', px: 1.5, py: .8, boxSizing: 'border-box', bgcolor: meta.pale }}><Typography variant="caption">成本</Typography><Typography variant="body2" fontWeight={900} color={meta.color}>${card.cost?.cash ?? 0} 萬</Typography></Stack></MuiCard>
}
