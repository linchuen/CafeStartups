import placeholder from '../../../../assets/card-placeholder.svg'
import barista from '../../../../assets/partners/partner-barista.png'
import roaster from '../../../../assets/partners/partner-roaster.png'
import marketer from '../../../../assets/partners/partner-marketer.png'
import taste from '../../../../assets/partners/partner-taste.png'
import finance from '../../../../assets/partners/partner-finance.png'
import pastry from '../../../../assets/partners/partner-pastry.png'
import supply from '../../../../assets/partners/partner-supply.png'
import community from '../../../../assets/partners/partner-community.png'
import hr from '../../../../assets/partners/partner-hr.png'
import analytics from '../../../../assets/partners/partner-analytics.png'
import beans from '../../../../assets/partners/partner-beans.png'
import value from '../../../../assets/partners/partner-value.png'
import marketingResource from '../../../../assets/partners/partner-marketing-resource.png'
import type { PlayerCard } from '../../model/cardTypes'

const partnerArtwork: Record<string, string> = {
  'partner-barista': barista,
  'partner-roaster': roaster,
  'partner-marketer': marketer,
  'partner-service': taste,
  'partner-finance': finance,
  'partner-pastry': pastry,
  'partner-supply': supply,
  'partner-community': community,
  'partner-hr': hr,
  'partner-analytics': analytics,
  'partner-beans': beans,
  'partner-value': value,
  'partner-marketing-resource': marketingResource,
}

export function CardArtwork({ card }: { card: PlayerCard }) {
  return <img src={partnerArtwork[card.id] ?? placeholder} alt="" style={{ display: 'block', width: '100%', height: '100%', boxSizing: 'border-box', padding: '8px 16px', objectFit: 'contain', objectPosition: 'center' }} />
}
