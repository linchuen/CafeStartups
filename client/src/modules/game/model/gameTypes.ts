import type { CardFaceData } from './cardTypes'

export type Card = CardFaceData

export type Player = {
  id: string
  displayName: string
  bot?: boolean
  ready?: boolean
  cash: number
  loans: number
  revenue?: number
  score?: number
  selectedKPIs?: string[]
  brandAwareness?: number
  products?: number
  values?: number
  resources?: number
  gourmetSatisfaction?: number
  regularSatisfaction?: number
  iconValues?: Record<string, number>
  handCount: number
}

export type CashFlowStatement = {
  period: number
  beginningCash: number
  operatingRevenue: number
  gourmetRevenue: number
  regularRevenue: number
  gourmetCustomerCount: number
  regularCustomerCount: number
  otherIncome: number
  operatingExpenses: number
  interestPaid: number
  principalRepayment: number
  newLoans: number
  endingCash: number
}

export type GameState = {
  id: string
  status: string
  seed: string
  gameVersion: number
  period: number
  phase: string
  round: number
  demandBoard?: Record<string, number>
  demandCards?: Record<'gourmet' | 'regular', { id: string; position: number; icons: string[]; revealed: boolean }[]>
  marketRanking?: number[]
  marketDraws?: { rank: number; playerId: string; playerName: string; customerCounts: Record<'gourmet' | 'regular' | 'difficult', number>; total: number }[]
  marketBag?: Record<'gourmet' | 'regular' | 'difficult', number>
  partnerOptions?: Card[]
  starterShopOptions?: Card[]
  players: Player[]
  me?: {
    id: string
    hand: Card[]
    tableau: Card[]
    discardCount: number
    partner?: Card
    starterShop?: Card
    initialCardsSelected?: boolean
    cash: number
    loans: number
    customers?: { kind: string; demand: string; unitPrice: number; count: number }[]
    revenue?: number
    score?: number
    selectedKPIs?: string[]
    kpiSelectionPeriod?: number
    cashFlow?: CashFlowStatement[]
    cashFlowRounds?: CashFlowStatement[]
    brandAwareness?: number
    products?: number
    values?: number
    resources?: number
    gourmetSatisfaction?: number
    regularSatisfaction?: number
    iconValues?: Record<string, number>
  }
}

export type ApiError = Error & { code?: string }
