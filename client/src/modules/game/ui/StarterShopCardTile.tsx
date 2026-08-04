// @ts-nocheck
import { Storefront } from '@mui/icons-material'
import { Box, Card as MuiCard, Stack, Typography } from '@mui/material'
import type { PlayerCard } from './gameDashboardCardTypes'

const shopMeta: Record<string, { color: string; pale: string; role: string }> = {
  'starter-songshan': { color: '#2f7f82', pale: '#dff1ef', role: '饕客聚集' },
  'starter-minsheng': { color: '#3e7896', pale: '#e1eff5', role: '饕客與一般客' },
  'starter-xinyi': { color: '#8b6a38', pale: '#f3ecd8', role: '分店客群拓展' },
  'starter-station': { color: '#735b91', pale: '#eee7f5', role: '一般顧客聚集' },
}

export function StarterShopCardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  const meta = shopMeta[card.id] ?? { color: '#287477', pale: '#dff0ed', role: '店面顧客來源' }
  const gourmet = card.demand?.gourmet ?? 0
  const regular = card.demand?.regular ?? 0
  return <MuiCard onClick={onClick} variant="outlined" sx={{ height: '100%', minHeight: 290, overflow: 'hidden', cursor: onClick ? 'pointer' : 'default', bgcolor: meta.pale, borderColor: selected ? meta.color : `${meta.color}66`, borderWidth: selected ? 3 : 1 }}><Box sx={{ px: 1.5, py: .9, bgcolor: meta.color, color: 'white', textAlign: 'center' }}><Typography variant="caption" fontWeight={900}>{card.name}</Typography></Box><Box sx={{ px: 1.5, py: .85, display: 'flex', alignItems: 'center', gap: 1, bgcolor: `${meta.color}dd`, color: 'white' }}><Storefront sx={{ fontSize: 20 }} /><Box><Typography variant="caption" sx={{ display: 'block', opacity: .8 }}>店面功能</Typography><Typography variant="body2" fontWeight={900}>{meta.role}</Typography></Box></Box><Box sx={{ display: 'grid', placeItems: 'center', minHeight: 105, color: meta.color }}><Storefront sx={{ fontSize: 64 }} /></Box><Box sx={{ px: 1.5, py: 1, bgcolor: '#ffffffcc' }}><Typography variant="caption" color="text.secondary">顧客效果</Typography><Stack direction="row" spacing={1.5} sx={{ minHeight: 34 }}>{gourmet > 0 && <Typography fontWeight={900} sx={{ color: '#ff7900' }}>{'★'.repeat(gourmet)} <small>饕客</small></Typography>}{regular > 0 && <Typography fontWeight={900} sx={{ color: '#e5b832' }}>{'★'.repeat(regular)} <small>一般客</small></Typography>}{gourmet === 0 && regular === 0 && <Typography variant="body2">{card.effect ?? '開店後可獲得對應顧客。'}</Typography>}</Stack><Typography variant="caption" color="text.secondary">{card.description}</Typography></Box><Stack direction="row" justifyContent="space-between" sx={{ px: 1.5, py: .8, bgcolor: meta.pale }}><Typography variant="caption">店面成本</Typography><Typography variant="body2" fontWeight={900} color={meta.color}>${card.cost?.cash ?? 0} 萬</Typography></Stack></MuiCard>
}
