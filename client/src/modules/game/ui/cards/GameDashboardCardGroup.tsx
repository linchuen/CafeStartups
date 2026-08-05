// @ts-nocheck
import { Box, Grid, Stack, Typography } from '@mui/material'
import { ManagementCardTile } from './ManagementCardTile'
import { PartnerCardTile } from './PartnerCardTile'
import { StarterShopCardTile } from './StarterShopCardTile'
import type { PlayerCard } from '../../model/cardTypes'

function CardTile({ card, selected, onClick }: { card: PlayerCard; selected?: boolean; onClick?: () => void }) {
  if (card.kind === 'partner') return <PartnerCardTile card={card} selected={selected} onClick={onClick} />
  if (card.kind === 'starter_shop') return <StarterShopCardTile card={card} selected={selected} onClick={onClick} />
  return <ManagementCardTile card={card} selected={selected} onClick={onClick} />
}

export function GameDashboardCardGroup({ title, subtitle, cards, selectedId, onSelect, accentColor }: { title: string; subtitle: string; cards: PlayerCard[]; selectedId?: string; onSelect?: (id: string) => void; accentColor?: string }) {
  if (!cards.length) return null
  return <Box sx={{ mt: 2.5 }}><Stack direction="row" alignItems="baseline" spacing={1} sx={{ mb: 1.2, pl: 1, borderLeft: accentColor ? `4px solid ${accentColor}` : undefined }}><Typography variant="subtitle1" fontWeight={800} sx={{ color: accentColor }}>{title}</Typography><Typography variant="caption" color="text.secondary">{subtitle}</Typography></Stack><Grid container spacing={1.5}>{cards.map((card) => <Grid key={card.id} size={{ xs: 12, sm: 6, md: 3, lg: 2 }}><CardTile card={card} selected={selectedId === card.id} onClick={onSelect ? () => onSelect(card.id) : undefined} /></Grid>)}</Grid></Box>
}
