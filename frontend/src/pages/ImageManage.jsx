import { Layout, Image, Button, Upload, Drawer, message } from "antd";
import {
  UploadOutlined,
  FolderOutlined,
  HomeOutlined,
  InboxOutlined,
  DeleteOutlined,
} from "@ant-design/icons";
import { useImmerReducer } from "use-immer";
import { useEffect, useState } from "react";
import path from "path-browserify";

import { ImageList, FileList } from "./../reducers/ImageManageReducer";
import { ReadDir, UploadFile, DeleteFile } from "../apis/SimpleUpload";
import PathTravel from "../components/PathTravel";

const Header = Layout.Header;
const Content = Layout.Content;
const pathPrefix = "/simple_upload/object";
const { Dragger } = Upload;

function ImageManage() {
  const [messageApi, contextHolder] = message.useMessage();

  const [currentDir, setCurrentDir] = useState("/");
  const [showUploadDrawer, setShowUploadDrawer] = useState(false);
  const [pathItems, setPathItems] = useState([
    {
      title: <HomeOutlined />,
      path: "/",
    },
  ]);

  const [imageList, dispatch] = useImmerReducer(ImageList, []);
  const [fileList, dispatchFileList] = useImmerReducer(FileList, []);

  useEffect(() => {
    ReadDir(currentDir).then((data) => {
      dispatch({
        type: "list",
        payload: data.List ? data.List : [],
      });
    });
  }, []);

  const changeDir = (dirName) => {
    let nextDir = path.normalize(path.join(currentDir, dirName));
    setCurrentDir(nextDir);
    setPathItems([
      ...pathItems,
      {
        title: dirName,
        path: nextDir,
      },
    ]);
    ReadDir(nextDir).then((data) => {
      dispatch({
        type: "list",
        payload: data?.List,
      });
    });
  };

  const changeDirAbsolute = ({ absolutePath }) => {
    if (absolutePath === currentDir) {
      return;
    }
    setCurrentDir(absolutePath);
    // 重新生成
    let absPath = "";
    let newPathItems = [];

    for (let item of pathItems) {
      absPath = path.join(absPath, item.path);
      newPathItems.push(item);
      if (absPath === absolutePath) {
        break;
      }
    }
    setPathItems(newPathItems);

    ReadDir(absolutePath).then((data) => {
      dispatch({
        type: "list",
        payload: data?.List,
      });
    });
  };

  const items = imageList.map((item) => {
    if (item.FileType == "dir") {
      return (
        <div
          key={item.Name}
          style={{
            height: "300px",
            width: "300px",
            border: "1px solid #ccc",
          }}
          onClick={() => changeDir(item.Name)}
        >
          <FolderOutlined
            style={{
              fontSize: "150px",
              margin: "auto",
              display: "block",
              lineHeight: "250px",
            }}
          />
          <div
            style={{
              fontSize: "3em",
              width: "100%",
              textAlign: "center",
              marginTop: "-90px",
            }}
          >
            {item.Name}
          </div>
        </div>
      );
    } else {
      return (
        <div
          style={{
            height: "300px",
            width: "300px",
            position: "relative",
          }}
        >
          <Image
            height={"100%"}
            width={"100%"}
            src={path.join(pathPrefix, currentDir, item.Name)}
          />
          <Button
            style={{
              position: "absolute",
              top: 0,
              right: 0,
            }}
            type="dashed"
            shape="circle"
            ghost
            icon={<DeleteOutlined />}
            onClick={() => {
              deleteFile(path.join(currentDir, item.Name));
            }}
          ></Button>
        </div>
      );
    }
  });

  const uploadFile = () => {
    if (fileList.length == 0) {
      messageApi.open({
        type: "error",
        content: "请选择文件",
      });
    }

    fileList.map((file) => {
      const stream = new ReadableStream({
        start(controller) {
          const reader = new FileReader();

          reader.onload = () => {
            const chunkSize = 1024 * 1024; // 每次读取 1MB
            let offset = 0;

            const readNextChunk = () => {
              const chunk = reader.result.slice(offset, offset + chunkSize);
              if (chunk.byteLength > 0) {
                controller.enqueue(new Uint8Array(chunk));
                offset += chunkSize;

                // 计算上传进度
                let percent = (offset / file.size) * 100;
                if (percent > 100) {
                  percent = 100;
                }
                dispatchFileList({
                  type: "process",
                  payload: {
                    file: file,
                    process: percent,
                  },
                });

                // 继续读取下一块
                readNextChunk();
              } else {
                controller.close();
              }
            };

            readNextChunk();
          };

          reader.readAsArrayBuffer(file);
        },
      });
      file.stream = stream;

      let formData = new FormData();
      formData.append("File", file);
      dispatchFileList({
        type: "process",
        payload: {
          file: file,
          process: 0,
        },
      });
      UploadFile("/img/" + file.name + "", formData)
        .then((data) => {
          dispatchFileList({
            type: "done",
            payload: {
              file: file,
            },
          });
        })
        .catch((err) => {
          console.log(err);
          dispatchFileList({
            type: "error",
            payload: {
              file: file,
            },
          });
          messageApi.open({
            type: "error",
            content: "上传失败:" + err.message,
            duration: 5,
          });
        });
    });
  };

  const deleteFile = (filepath) => {
    DeleteFile(filepath)
      .then((data) => {
        // 删除成功
        dispatch({
          type: "remove",
          payload: path.basename(filepath),
        });
      })
      .catch((err) => {
        messageApi.open({
          type: "error",
          content: "删除失败:" + err.message,
          duration: 5,
        });
      });
  };

  return (
    <>
      {contextHolder}
      <Layout>
        <Header
          style={{
            backgroundColor: "white",
            display: "flex",
            justifyContent: "space-between",
          }}
        >
          <PathTravel
            items={pathItems}
            onClick={changeDirAbsolute}
          ></PathTravel>
          <div
            style={{
              float: "right",
            }}
          >
            <Button
              icon={<UploadOutlined />}
              onClick={() => {
                setShowUploadDrawer(true);
              }}
            >
              上传
            </Button>
          </div>
        </Header>
        <Content>
          <Drawer
            open={showUploadDrawer}
            width="50%"
            maskClosable={false}
            onClose={() => {
              setShowUploadDrawer(false);
              dispatchFileList({
                type: "clear",
              });
            }}
            extra={
              <Button type="primary" onClick={uploadFile}>
                开始上传
              </Button>
            }
          >
            <div style={{ height: "15%" }}>
              <Dragger
                beforeUpload={(file) => {
                  dispatchFileList({
                    type: "add",
                    payload: file,
                  });
                  return false;
                }}
                onRemove={(file) => {
                  dispatchFileList({
                    type: "remove",
                    payload: file,
                  });
                }}
                fileList={fileList}
                multiple={true}
                listType="picture"
              >
                <p className="ant-upload-drag-icon">
                  <InboxOutlined />
                </p>
              </Dragger>
            </div>
          </Drawer>

          <div
            style={{
              padding: "0px 10px",
              display: "flex",
              flexWrap: "wrap",
              // justifyContent: "space-evenly",
              gap: "10px",
            }}
          >
            {items}
          </div>
        </Content>
      </Layout>
    </>
  );
}

export default ImageManage;
