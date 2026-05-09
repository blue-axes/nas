import {
  Layout,
  Button,
  Upload,
  Drawer,
  message,
  Empty,
  Table,
  Tooltip,
} from "antd";
import {
  UploadOutlined,
  FolderOutlined,
  HomeOutlined,
  InboxOutlined,
  DeleteOutlined,
  FileOutlined,
  DownloadOutlined,
} from "@ant-design/icons";
import { useImmerReducer } from "use-immer";
import { useEffect, useState } from "react";
import path from "path-browserify";

import { ImageList, FileList } from "./../reducers/ImageManageReducer";
import { ReadDir, UploadFile, DeleteFile } from "../apis/SimpleUpload";
import PathTravel from "../components/PathTravel";
import useScreenWidth from "../hooks/useScreenWidth";

const pathPrefix = "/simple_upload/object";
const { Dragger } = Upload;

function FileManage() {
  const sw = useScreenWidth();
  const drawerWidth = sw < 768 ? "100%" : "50%";
  const [messageApi, contextHolder] = message.useMessage();

  const [currentDir, setCurrentDir] = useState("/other/");
  const [showUploadDrawer, setShowUploadDrawer] = useState(false);
  const [pathItems, setPathItems] = useState([
    { title: <HomeOutlined />, path: "/" },
    { title: "other", path: "/other/" },
  ]);

  const [imageList, dispatch] = useImmerReducer(ImageList, []);
  const [fileList, dispatchFileList] = useImmerReducer(FileList, []);

  useEffect(() => {
    ReadDir(currentDir).then((data) => {
      dispatch({ type: "list", payload: data.List ? data.List : [] });
    });
  }, []);

  const changeDir = (dirName) => {
    const nextDir = path.normalize(path.join(currentDir, dirName));
    setCurrentDir(nextDir);
    setPathItems([...pathItems, { title: dirName, path: nextDir }]);
    ReadDir(nextDir).then((data) => {
      dispatch({ type: "list", payload: data?.List });
    });
  };

  const changeDirAbsolute = ({ absolutePath }) => {
    if (absolutePath === currentDir) return;
    setCurrentDir(absolutePath);
    let absPath = "";
    const newPathItems = [];
    for (const item of pathItems) {
      absPath = path.join(absPath, item.path);
      newPathItems.push(item);
      if (absPath === absolutePath) break;
    }
    setPathItems(newPathItems);
    ReadDir(absolutePath).then((data) => {
      dispatch({ type: "list", payload: data?.List });
    });
  };

  const deleteFile = (filepath) => {
    DeleteFile(filepath)
      .then(() => dispatch({ type: "remove", payload: path.basename(filepath) }))
      .catch((err) => {
        messageApi.open({
          type: "error",
          content: "Delete failed: " + err.message,
          duration: 5,
        });
      });
  };

  const uploadFile = () => {
    if (fileList.length == 0) {
      messageApi.open({ type: "error", content: "Please select files" });
      return;
    }
    fileList.map((file) => {
      const stream = new ReadableStream({
        start(controller) {
          const reader = new FileReader();
          reader.onload = () => {
            const chunkSize = 1024 * 1024;
            let offset = 0;
            const readNextChunk = () => {
              const chunk = reader.result.slice(offset, offset + chunkSize);
              if (chunk.byteLength > 0) {
                controller.enqueue(new Uint8Array(chunk));
                offset += chunkSize;
                let percent = (offset / file.size) * 100;
                if (percent > 100) percent = 100;
                dispatchFileList({
                  type: "process",
                  payload: { file, process: percent },
                });
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
      const formData = new FormData();
      formData.append("File", file);
      dispatchFileList({ type: "process", payload: { file, process: 0 } });
      UploadFile(path.join(currentDir, file.name), formData)
        .then(() => dispatchFileList({ type: "done", payload: { file } }))
        .catch((err) => {
          console.log(err);
          dispatchFileList({ type: "error", payload: { file } });
          messageApi.open({
            type: "error",
            content: "Upload failed: " + err.message,
            duration: 5,
          });
        });
    });
  };

  const formatSize = (bytes) => {
    if (!bytes) return "-";
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
    if (bytes < 1024 * 1024 * 1024)
      return (bytes / (1024 * 1024)).toFixed(1) + " MB";
    return (bytes / (1024 * 1024 * 1024)).toFixed(1) + " GB";
  };

  const getExt = (name) => {
    const idx = name.lastIndexOf(".");
    if (idx < 0) return "";
    return name.slice(idx + 1).toUpperCase();
  };

  const columns = [
    {
      title: "",
      width: 40,
      render: (_, record) => {
        if (record.FileType === "dir") {
          return (
            <FolderOutlined
              style={{
                color: "#00d4ff",
                fontSize: "18px",
                filter: "drop-shadow(0 0 6px rgba(0,212,255,0.3))",
              }}
            />
          );
        }
        return (
          <FileOutlined
            style={{ color: "#94a3b8", fontSize: "18px" }}
          />
        );
      },
    },
    {
      title: "Name",
      dataIndex: "Name",
      key: "Name",
      render: (text, record) => {
        if (record.FileType === "dir") {
          return (
            <a
              onClick={() => changeDir(record.Name)}
              style={{ fontWeight: 500, color: "#00d4ff" }}
            >
              {text}
            </a>
          );
        }
        return (
          <a
            href={path.join(pathPrefix, currentDir, record.Name)}
            target="_blank"
            rel="noreferrer"
            style={{ color: "#e2e8f0" }}
          >
            {text}
          </a>
        );
      },
    },
    {
      title: "Type",
      width: 100,
      render: (_, record) => {
        if (record.FileType === "dir") return "Folder";
        return getExt(record.Name) || "File";
      },
    },
    {
      title: "Size",
      width: 120,
      render: (_, record) => {
        if (record.FileType === "dir") return "-";
        return formatSize(record.Size);
      },
    },
    {
      title: "",
      width: 80,
      render: (_, record) => {
        if (record.FileType === "dir") return null;
        return (
          <div style={{ display: "flex", gap: "4px" }}>
            <Tooltip title="Download">
              <Button
                type="text"
                size="small"
                icon={<DownloadOutlined />}
                href={path.join(pathPrefix, currentDir, record.Name)}
                target="_blank"
              />
            </Tooltip>
            <Tooltip title="Delete">
              <Button
                type="text"
                size="small"
                danger
                icon={<DeleteOutlined />}
                onClick={() => deleteFile(path.join(currentDir, record.Name))}
              />
            </Tooltip>
          </div>
        );
      },
    },
  ];

  return (
    <>
      {contextHolder}
      <div className="tech-page">
        <div className="tech-header">
          <PathTravel items={pathItems} onClick={changeDirAbsolute} />
          <Button
            type="primary"
            icon={<UploadOutlined />}
            onClick={() => setShowUploadDrawer(true)}
          >
            Upload
          </Button>
        </div>

        <div className="tech-content">
          <Drawer
            open={showUploadDrawer}
            width={drawerWidth}
            maskClosable={false}
            onClose={() => {
              setShowUploadDrawer(false);
              dispatchFileList({ type: "clear" });
            }}
            extra={
              <Button type="primary" onClick={uploadFile}>
                Start Upload
              </Button>
            }
          >
            <div style={{ height: "15%" }}>
              <Dragger
                beforeUpload={(file) => {
                  dispatchFileList({ type: "add", payload: file });
                  return false;
                }}
                onRemove={(file) =>
                  dispatchFileList({ type: "remove", payload: file })
                }
                fileList={fileList}
                multiple={true}
                listType="picture"
              >
                <p className="ant-upload-drag-icon">
                  <InboxOutlined />
                </p>
                <p style={{ color: "#94a3b8" }}>
                  Drag files here or click to select
                </p>
              </Dragger>
            </div>
          </Drawer>

          {imageList.length === 0 ? (
            <div className="tech-empty">
              <Empty description="Empty directory" />
            </div>
          ) : (
            <Table
              dataSource={imageList}
              columns={columns}
              rowKey="Name"
              pagination={false}
              size="middle"
              onRow={(record) => {
                if (record.FileType === "dir") {
                  return {
                    onDoubleClick: () => changeDir(record.Name),
                    style: { cursor: "pointer" },
                  };
                }
                return {};
              }}
            />
          )}
        </div>
      </div>
    </>
  );
}

export default FileManage;