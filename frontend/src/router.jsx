import App from './App.jsx'
import ImageManage from './pages/ImageManage.jsx'
import Sub from './Sub.jsx'
import { createHashRouter } from 'react-router'

const router = createHashRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {
        path: '/image',
        element: <ImageManage />
      },
      {
        path: '/video',
        element: <Sub />
      },
      {
        path: '/file',
        element: <Sub />
      }
    ]
  }
])

export default router
