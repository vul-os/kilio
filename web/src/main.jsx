import React from 'react'
import ReactDOM from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import './styles/base.css'
import App from './App.jsx'
import ReporterApp from './reporter/ReporterApp.jsx'
import HandlerApp from './handler/HandlerApp.jsx'

const router = createBrowserRouter([
  { path: '/', element: <App /> },
  { path: '/report/*', element: <ReporterApp /> },
  { path: '/handler/*', element: <HandlerApp /> },
])

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
)
