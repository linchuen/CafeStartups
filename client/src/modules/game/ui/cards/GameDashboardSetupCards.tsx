// @ts-nocheck
import { useState } from 'react'
import { Box, Button, Grid, Paper, Stack, Typography } from '@mui/material'
import { PartnerCardTile } from './PartnerCardTile'
import { StarterShopCardTile } from './StarterShopCardTile'
import type { PlayerCard } from '../../model/cardTypes'

export function GameDashboardSetupCards({ partners, shops, selectedPartnerId, selectedShopId, busy, onSelect, onBegin }: { partners: PlayerCard[]; shops: PlayerCard[]; selectedPartnerId?: string; selectedShopId?: string; busy: boolean; onSelect: (partnerId: string, shopId: string) => void; onBegin: () => void }) {
  const [hasChosen, setHasChosen] = useState(false)
  if (!partners.length || !shops.length) return null
  const choose = (partnerId: string, shopId: string) => { setHasChosen(true); onSelect(partnerId, shopId) }
  return <Paper sx={{ p: { xs: 2, md: 2.5 }, mb: 2.5, border: '2px solid #b36f42', borderRadius: 3, bgcolor: '#fffaf4' }}><Stack spacing={.6} sx={{ mb: 2 }}><Typography variant="overline" color="primary">ROUND 0 · START YOUR CAFÉ</Typography><Typography variant="h5" fontWeight={900}>選擇你的合夥人與初始店面</Typography><Typography variant="body2" color="text.secondary">選擇完成後，確認才會進入下一回合。</Typography></Stack><Typography variant="subtitle1" fontWeight={900} sx={{ mb: 1 }}>合夥人卡</Typography><Grid container spacing={1.5} sx={{ mb: 2.5 }}>{partners.map((card) => <Grid key={card.id} size={{ xs: 12, sm: 6, md: 3 }}><PartnerCardTile card={card} selected={selectedPartnerId === card.id} onClick={busy ? undefined : () => choose(card.id, selectedShopId ?? shops[0].id)} /></Grid>)}</Grid><Typography variant="subtitle1" fontWeight={900} sx={{ mb: 1 }}>初始店面</Typography><Grid container spacing={1.5}>{shops.map((card) => <Grid key={card.id} size={{ xs: 12, sm: 6, md: 3 }}><StarterShopCardTile card={card} selected={selectedShopId === card.id} onClick={busy ? undefined : () => choose(selectedPartnerId ?? partners[0].id, card.id)} /></Grid>)}</Grid><Stack direction="row" justifyContent="flex-end" sx={{ mt: 2 }}><Button variant="contained" onClick={onBegin} disabled={busy || !hasChosen}>確認選擇並進入下一回合</Button></Stack></Paper>
}
