import App from './App.jsx'
import ImageManage from './pages/ImageManage.jsx'
import VideoManage from './pages/VideoManage.jsx'
import FileManage from './pages/FileManage.jsx'
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
        element: <VideoManage />
      },
      {
        path: '/file',
        element: <FileManage />
      }
    ]
  }
])

export default router
