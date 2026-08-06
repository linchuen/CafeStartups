export type CardFaceData = {
  id: string
  name: string
  kind: string
  period: number
  description?: string
  function?: string
  colorKey?: string
  startingCashBonus?: number
  starterShopId?: string
  minPlayers?: number
  cost: { cash: number; icons: string[] }
  icons: string[]
  marketChange?: Record<string, number>
  customerCount?: Record<string, number>
  brandAwareness?: number
}

export type PlayerCard = CardFaceData
