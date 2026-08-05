import placeholder from '../../../assets/card-placeholder.svg'
import type { PlayerCard } from './gameDashboardCardTypes'

// Add artwork here by card name when the final illustrations are available.
const artworkByName: Record<string, string> = {}

export function CardArtwork({ card }: { card: PlayerCard }) {
  return <img src={artworkByName[card.name] ?? placeholder} alt="" style={{ display: 'block', width: '100%', height: '100%', objectFit: 'cover' }} />
}
