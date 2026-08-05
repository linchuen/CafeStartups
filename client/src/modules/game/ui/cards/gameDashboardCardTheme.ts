import { Assignment, Campaign, LocalShipping, Settings } from '@mui/icons-material'

export type CardColorTheme = { color: string; pale: string }

export const cardColorTokens = {
  resource: { color: '#2d6897', pale: '#d5e7ef' },
  product: { color: '#c88d28', pale: '#f3dfb7' },
  value: { color: '#b44f52', pale: '#f3d9d3' },
  channel: { color: '#3f7d66', pale: '#e2f0e8' },
  marketing: { color: '#7a5ba5', pale: '#eee7f7' },
} as const

export const costStructureIcons = [
  { id: 'data', label: '資料資源', icon: Assignment },
  { id: 'marketing', label: '行銷資源', icon: Campaign },
  { id: 'operations', label: '營運資源', icon: Settings },
  { id: 'procurement', label: '採購資源', icon: LocalShipping },
] as const

export const dashboardCardColors = {
  starterShop: {
    'starter-songshan': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-minsheng': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-xinyi': { color: '#3f7d66', pale: '#e2f0e8' },
  'starter-station': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-daan': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-beitou': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-neihu': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-banqiao': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-ximen': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-gongguan': { color: '#3f7d66', pale: '#e2f0e8' },
  } satisfies Record<string, CardColorTheme>,
  management: cardColorTokens satisfies Record<string, CardColorTheme>,
} as const
