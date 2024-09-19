import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import axios from 'axios'

interface Coffee {
  id: number
  title: string
  price: number
}

function App() {
  const [coffees, setCoffees] = useState<Coffee[]>([])

  useEffect(() => {
    const fetchCoffees = async () => {
      try {
        const response = await axios.get(`${import.meta.env.VITE_INVENTORY_SERVICE_URL}/coffees`)
        setCoffees(response.data?.data)
      } catch (error) {
        console.error('Error fetching coffees:', error)
      }
    }

    fetchCoffees()
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-br from-amber-100 to-amber-200 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <h1 className="text-5xl font-extrabold text-center text-brown-800 mb-2">Brew Haven</h1>
        <p className="text-xl text-center text-brown-600 mb-12">Discover Your Perfect Cup</p>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
          {coffees.map((coffee) => (
            <div key={coffee.id} className="bg-white rounded-lg shadow-md overflow-hidden transition-transform hover:scale-105">
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-800 mb-2">{coffee.title}</h2>
                <p className="text-2xl font-bold text-amber-600 mb-4">${coffee.price.toFixed(2)}</p>
                <Link 
                  to={`/purchase/${coffee.id}`}
                  className="block w-full bg-amber-500 hover:bg-amber-600 text-white font-bold py-2 px-4 rounded transition-colors duration-200"
                >
                  Purchase
                </Link>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default App
