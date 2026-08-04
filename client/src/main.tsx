import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './styles.css'
import { App } from './App'
import { CssBaseline, ThemeProvider } from '@mui/material'
import { cafeTheme } from './app/theme'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider theme={cafeTheme}>
      <CssBaseline />
      <App />
    </ThemeProvider>
  </StrictMode>,
)
