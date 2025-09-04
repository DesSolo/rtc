import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import Index from './pages/Index'

createRoot(document.getElementById('root')).render(
    <BrowserRouter>
        <Index />
    </BrowserRouter>
)
