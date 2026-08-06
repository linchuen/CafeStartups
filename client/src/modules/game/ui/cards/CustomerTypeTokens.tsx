import { Box } from '@mui/material'

export type CustomerType = 'gourmet' | 'regular' | 'difficult'

const customerTypeMeta: Record<CustomerType, { label: string; color: string }> = {
  gourmet: { label: '饕客', color: '#ff7900' },
  regular: { label: '一般客', color: '#e5b832' },
  difficult: { label: '困難客', color: '#1976d2' },
}

export function CustomerTypeTokens({ type, count, size = 18 }: { type: CustomerType; count: number; size?: number }) {
  const meta = customerTypeMeta[type]
  const tokenCount = Math.max(1, Math.abs(count))

  return <Box component="span" role="img" aria-label={`${meta.label} ${count}`} sx={{ display: 'inline-flex', alignItems: 'center', gap: .4, color: meta.color, lineHeight: 1 }}>
    {count < 0 && <Box component="span" sx={{ fontSize: size * .8 }}>−</Box>}
    {Array.from({ length: tokenCount }, (_, index) => <Box key={index} component="span" sx={{ display: 'inline-block', width: size, height: size, border: '1px solid rgba(80, 50, 20, .2)', borderRadius: '3px', backgroundColor: 'currentColor', boxShadow: 'inset 0 0 0 3px rgba(255,255,255,.2)' }} />)}
  </Box>
}
