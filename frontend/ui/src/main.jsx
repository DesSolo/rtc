import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import Index from './pages/Index'
import './index.css'

createRoot(document.getElementById('root')).render(
    <BrowserRouter basename="/ui">
        <Index />
    </BrowserRouter>
)
