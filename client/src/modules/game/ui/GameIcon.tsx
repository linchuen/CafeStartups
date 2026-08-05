import type { ComponentType } from 'react'
import type { SvgIconProps } from '@mui/material/SvgIcon'
import { Assignment, AttachMoney, Cake, Campaign, Coffee, Favorite, Groups, HelpOutlined, LocalShipping, Restaurant, Settings, Spa, Storefront, VolunteerActivism } from '@mui/icons-material'

const iconMap: Record<string, ComponentType<SvgIconProps>> = {
  assignment: Assignment,
  beans: Spa,
  channel: Storefront,
  coffee: Coffee,
  data: Assignment,
  dessert: Cake,
  marketing: Campaign,
  marketing_resource: Campaign,
  operations: Settings,
  people: Groups,
  price: AttachMoney,
  procurement: LocalShipping,
  service: VolunteerActivism,
  taste: Restaurant,
  value: Favorite,
}

export function GameIcon({ name, ...props }: SvgIconProps & { name: string }) {
  const Icon = iconMap[name] ?? HelpOutlined
  return <Icon {...props} />
}
