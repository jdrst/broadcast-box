import React from 'react'
import { Routes, Route } from 'react-router-dom'

import Header from './components/header'
import Selection from './components/selection'
import PlayerPage from './components/player'
import Publish from './components/publish'
import LoginPage from './components/login'
import OverviewPage from './components/overview'

function App() {
  return (
    <Routes>
      <Route path='/' element={<Header />}>
        {/* <Route index element={<Selection />} /> */}
        <Route index element={<OverviewPage />} />
        <Route path='/publish/*' element={<Publish />} />
        <Route path='/login' element={<LoginPage />} />
        <Route path='/*' element={<PlayerPage />} />
      </Route>
    </Routes>
  )
}

export default App
