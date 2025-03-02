import { useState } from 'react'
import { Outlet } from 'react-router'
import { useNavigate } from "react-router"
import { Button, ConfigProvider, Layout, Menu, } from 'antd'
import { DeliveredProcedureOutlined, PictureOutlined, PlaySquareOutlined, FileOutlined } from '@ant-design/icons'
import zh_CN from 'antd/locale/zh_CN'
import 'antd/dist/reset.css'

const Sider = Layout.Sider
const Header = Layout.Header
const Content = Layout.Content

const menuConfig = [
  {
    label: '文件管理',
    title: '文件管理',
    icon: <DeliveredProcedureOutlined />,
    children: [
      {
        label: '图片文件',
        title: '图片文件',
        icon: <PictureOutlined />,
        link: "/image"
      },
      {
        label: '视频文件',
        title: '视频文件',
        icon: <PlaySquareOutlined />,
        link: '/video'
      },
      {
        label: '普通文件',
        title: '普通文件',
        icon: <FileOutlined />,
        link: "/file"
      },

    ]
  }
]

function getItem(idx, keyPrefix, item) {
  if (!(item instanceof Object)) {
    return [{}, {}]
  }
  let key = keyPrefix + "key" + idx
  let children = []
  let keyLink = {}
  if (item?.children) {
    children = item?.children.map((subItem, subIdx) => {
      let [tmpItem, tmpKeyLink] = getItem(subIdx, key, subItem)
      if (tmpItem) {
        keyLink = { ...keyLink, ...tmpKeyLink }
        return tmpItem
      }
      return null
    })
  }


  if (children.length == 0) {
    children = null
  }

  keyLink[key] = item.link

  return [{
    key: key,
    icon: item?.icon || <QuestionOutlined />,
    label: item.label,
    title: item?.title || item.label,
    children: children,
    link: item.link,
  }, keyLink]
}

function initMenu() {
  let keyLink = {}
  let res = menuConfig.map((item, idx) => {
    let [tmpItem, tmpKeyLink] = getItem(idx, "", item)
    if (tmpItem) {
      keyLink = { ...keyLink, ...tmpKeyLink }
      return tmpItem
    }
    return null
  })

  return [res, keyLink]
}



function App() {
  const [items, keyLink] = initMenu()
  let navigate = useNavigate()

  let onClick = ({ key }) => {
    let to = keyLink[key]
    if (to != "") {
      navigate(to)
    }
    console.log("menu not link attr")
  }


  return (
    <>
      <ConfigProvider
        locale={zh_CN}
        componentSize="large"
      >
        <Layout style={{ width: "100vw", height: "100vh" }}>
          <Header>头</Header>
          <Layout>
            <Sider>
              <Menu
                mode="inline"
                triggerSubMenuAction="click"
                items={items}
                onClick={onClick}
              />
            </Sider>
            <Content>
              <Outlet />
            </Content>
          </Layout>
        </Layout >
      </ConfigProvider >

    </>
  )
}

export default App
