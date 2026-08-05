export type CardFaceData = {
  id: string
  name: string
  kind: string
  period: number
  description?: string
  effect?: string
  function?: string
  colorKey?: string
  startingCashBonus?: number
  starterShopId?: string
  demand?: Record<string, number>
  minPlayers?: number
  cost: { cash: number; icons: string[] }
  icons: string[]
  marketChange?: Record<string, number>
  customerCount?: Record<string, number>
}

export type PlayerCard = CardFaceData
