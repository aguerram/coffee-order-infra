import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom'
import CoffeeIndex from './CoffeeIndex.tsx'
import PurchaseIndex from './PurchaseIndex.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Router>
      <Routes>
        <Route path="/" element={<CoffeeIndex />} />
        <Route path="/purchase/:coffeeId" element={<PurchaseIndex />} />
      </Routes>
    </Router>
  </StrictMode>,
)
