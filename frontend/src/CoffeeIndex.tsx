import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'

interface Coffee {
  id: number
  title: string
  price: number
}

function App() {
  const [coffees, setCoffees] = useState<Coffee[]>([])

  useEffect(() => {
    fetchCoffees()
  }, [])

  const fetchCoffees = async () => {
    try {
      const response = await fetch(`${import.meta.env.VITE_INVENTORY_SERVICE_URL}/coffees`)
      const data = await response.json()      
      setCoffees(data?.data)
    } catch (error) {
      console.error('Error fetching coffees:', error)
    }
  }

  return (
    <div>
      <h1>Coffee Shop</h1>
      <ul>
        {coffees.map((coffee) => (
          <li key={coffee.id}>
            {coffee.title} - ${coffee.price.toFixed(2)}
            <Link to={`/purchase/${coffee.id}`}>
              <button>Purchase</button>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  )
}

export default App
