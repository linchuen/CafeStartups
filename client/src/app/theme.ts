import { createTheme } from '@mui/material/styles'

export const cafeTheme = createTheme({
  palette: {
    mode: 'light',
    primary: { main: '#714b38' },
    secondary: { main: '#287477' },
    background: { default: '#f8f5f1', paper: '#ffffff' },
  },
  typography: {
    fontFamily: 'Inter, "Noto Sans TC", system-ui, sans-serif',
    h4: { letterSpacing: '-0.03em' },
    h5: { letterSpacing: '-0.02em' },
  },
  shape: { borderRadius: 12 },
  components: {
    MuiButton: { defaultProps: { disableElevation: true } },
    MuiPaper: { defaultProps: { elevation: 0 } },
  },
})
