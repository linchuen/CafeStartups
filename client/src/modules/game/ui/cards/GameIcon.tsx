import type { ComponentType } from 'react'
import type { SvgIconProps } from '@mui/material/SvgIcon'
import { Assignment, AttachMoney, Cake, Campaign, Coffee, Favorite, Groups, HelpOutlined, LocalShipping, Restaurant, Settings, Spa, Storefront, VolunteerActivism } from '@mui/icons-material'
import iconDefinitions from './game-icon-map.json'

const iconComponents: Record<string, ComponentType<SvgIconProps>> = {
  Assignment,
  AttachMoney,
  Cake,
  Campaign,
  Coffee,
  Favorite,
  Groups,
  LocalShipping,
  Restaurant,
  Settings,
  Spa,
  Storefront,
  VolunteerActivism,
}

const iconMap: Record<string, ComponentType<SvgIconProps>> = Object.fromEntries(
  Object.entries(iconDefinitions).map(([name, definition]) => [name, iconComponents[definition.component] ?? HelpOutlined]),
)

export function GameIcon({ name, ...props }: SvgIconProps & { name: string }) {
  const Icon = iconMap[name] ?? HelpOutlined
  return <Icon {...props} titleAccess={props.titleAccess ?? iconDefinitions[name as keyof typeof iconDefinitions]?.label} />
}
