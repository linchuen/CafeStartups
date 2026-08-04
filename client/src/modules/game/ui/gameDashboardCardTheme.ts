export type CardColorTheme = { color: string; pale: string }

export const dashboardCardColors = {
  partner: {
    'partner-barista': { color: '#2f87a8', pale: '#dceff5' },
    'partner-roaster': { color: '#9a623b', pale: '#f3e3d4' },
    'partner-marketer': { color: '#7656a5', pale: '#eee5f7' },
    'partner-service': { color: '#39826c', pale: '#dff0e8' },
  } satisfies Record<string, CardColorTheme>,
  starterShop: {
    'starter-songshan': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-minsheng': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-xinyi': { color: '#3f7d66', pale: '#e2f0e8' },
    'starter-station': { color: '#3f7d66', pale: '#e2f0e8' },
  } satisfies Record<string, CardColorTheme>,
  management: {
    resource: { color: '#2d6897', pale: '#d5e7ef' },
    product: { color: '#c88d28', pale: '#f3dfb7' },
    value: { color: '#bd584f', pale: '#f3d9d3' },
    channel: { color: '#3f7d66', pale: '#e2f0e8' },
    marketing: { color: '#7a5ba5', pale: '#eee7f7' },
  } satisfies Record<string, CardColorTheme>,
} as const
