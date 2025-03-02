import App from './App.jsx'
import ImageManage from './pages/ImageManage.jsx'
import Sub from './Sub.jsx'
import { createBrowserRouter } from 'react-router'

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        path: "/image",
        element: <ImageManage />
      },
      {
        path: "/video",
        element: <Sub />
      },
      {
        path: "/file",
        element: <Sub />
      }
    ]
  }
])

export default router
