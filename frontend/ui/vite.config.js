import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
    plugins: [react()],
    base: "/ui/",
    server: {
        proxy: {
                '/api/v1': 'http://localhost:8080'
            }
        }
})
