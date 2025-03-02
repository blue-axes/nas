
import { Layout, Image } from 'antd'

const Header = Layout.Header
const Content = Layout.Content

function ImageManage() {
    
    return (<>
        <Layout >
            <Header style={{ backgroundColor: "white" }}>菜单</Header>
            <Content>
                <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'space-evenly', gap: "10px" }}>
                    <Image height={"300px"} width={"300px"} src="react.svg" />
                    <Image height={"300px"} width={"300px"} />
                    <Image height={"300px"} width={"300px"} />
                    <Image height={"300px"} width={"300px"} />
                    <Image height={"300px"} width={"300px"} />
                    <Image height={"300px"} width={"300px"} />
                    <Image height={"300px"} width={"300px"} />
                    <Image height={"300px"} width={"300px"} />
                </div>
            </Content>
        </Layout >
    </>)
}

export default ImageManage;