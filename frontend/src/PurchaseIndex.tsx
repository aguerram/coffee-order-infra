import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import axios from 'axios'

interface Coffee {
  id: number
  title: string
  price: number
}

function PurchaseIndex() {
  const { coffeeId } = useParams<{ coffeeId: string }>()
  const [coffee, setCoffee] = useState<Coffee | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [quantity, setQuantity] = useState(1)
  const navigate = useNavigate()

  useEffect(() => {
    fetchCoffee()
  }, [coffeeId])

  const fetchCoffee = async () => {
    try {
      const response = await axios.get(`${import.meta.env.VITE_INVENTORY_SERVICE_URL}/coffees/${coffeeId}`)
      setCoffee(response.data.data)
      setError(null)
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 404) {
        setError('Coffee not found')
      } else {
        setError('An error occurred while fetching the coffee')
      }
      setCoffee(null)
    }
  }

  const handlePurchase = async () => {
    try {
      const response = await axios.post(`${import.meta.env.VITE_ORDER_SERVICE_URL}/orders`, {
        coffee_id: coffee?.id,
        quantity: quantity
      })
      if (response.data.status === 'pending') {
        alert('Purchase successful!')
        navigate(`/`)
      }
    } catch (error) {
      setError('An error occurred while processing your purchase')
    }
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-amber-100 to-amber-200 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-3xl mx-auto text-center">
          <h1 className="text-3xl font-bold text-red-600 mb-4">{error}</h1>
          <Link to="/" className="text-amber-600 hover:text-amber-700">
            Return to Coffee List
          </Link>
        </div>
      </div>
    )
  }

  if (!coffee) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-amber-100 to-amber-200 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-3xl mx-auto text-center">
          <h1 className="text-3xl font-bold text-brown-800 mb-4">Loading...</h1>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-amber-100 to-amber-200 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <h1 className="text-4xl font-extrabold text-center text-brown-800 mb-8">{coffee?.title}</h1>
        <div className="bg-white rounded-lg shadow-md overflow-hidden p-6">
          <p className="text-2xl font-bold text-amber-600 mb-4">${coffee?.price.toFixed(2)}</p>
          <div className="flex items-center justify-between mb-4">
            <label htmlFor="quantity" className="text-lg font-medium text-gray-700">Quantity:</label>
            <input
              type="number"
              id="quantity"
              min="1"
              value={quantity}
              onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value)))}
              className="w-20 px-2 py-1 border rounded-md"
            />
          </div>
          <button 
            onClick={handlePurchase}
            className="w-full bg-amber-500 hover:bg-amber-600 text-white font-bold py-2 px-4 rounded transition-colors duration-200"
          >
            Confirm Purchase
          </button>
        </div>
        <div className="mt-8 text-center">
          <Link to="/" className="text-amber-600 hover:text-amber-700">
            Return to Coffee List
          </Link>
        </div>
      </div>
    </div>
  )
}

export default PurchaseIndex