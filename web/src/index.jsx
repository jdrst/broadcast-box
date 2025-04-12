import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import App from './App'
import { CinemaModeProvider } from './components/player'
import { BrowserRouter } from 'react-router-dom'
import { AuthContextProvider } from './components/login'

const root = ReactDOM.createRoot(document.getElementById('root'))
root.render(
  <React.StrictMode>
    <BrowserRouter basename={import.meta.env.PUBLIC_URL}>
      <AuthContextProvider>
        <CinemaModeProvider>
          <App />
        </CinemaModeProvider>
      </AuthContextProvider>
    </BrowserRouter>
  </React.StrictMode>
)
