import {
  Button,
  Upload,
  Drawer,
  message,
  Empty,
  Table,
  Tooltip,
  Modal,
  Input,
  Tag,
  Popover,
} from "antd";
import {
  UploadOutlined,
  FolderOutlined,
  HomeOutlined,
  InboxOutlined,
  DeleteOutlined,
  FileOutlined,
  DownloadOutlined,
  TagsOutlined,
  FolderAddOutlined,
} from "@ant-design/icons";
import { useImmerReducer } from "use-immer";
import { useEffect, useState, useCallback } from "react";
import path from "path-browserify";

import { ImageList, FileList } from "./../reducers/ImageManageReducer";
import { ReadDir, UploadFile, DeleteFile, SearchFiles, Mkdir } from "../apis/SimpleUpload";
import PathTravel from "../components/PathTravel";
import SearchBar from "../components/SearchBar";
import TagEditor from "../components/TagEditor";
import useScreenWidth from "../hooks/useScreenWidth";

const pathPrefix = "/simple_upload/object";
const { Dragger } = Upload;

function AllFileManage() {
  const sw = useScreenWidth();
  const [messageApi, contextHolder] = message.useMessage();

  const [currentDir, setCurrentDir] = useState("");
  const [showUploadDrawer, setShowUploadDrawer] = useState(false);
  const [pathItems, setPathItems] = useState([
    { title: <HomeOutlined />, path: "/" },
  ]);

  const [imageList, dispatch] = useImmerReducer(ImageList, []);
  const [fileList, dispatchFileList] = useImmerReducer(FileList, []);
  const [searching, setSearching] = useState(false);
  const [showNewFolder, setShowNewFolder] = useState(false);
  const [newFolderName, setNewFolderName] = useState("");
  const [creatingFolder, setCreatingFolder] = useState(false);

  useEffect(() => {
    ReadDir(currentDir).then((data) => {
      dispatch({ type: "list", payload: data.List ? data.List : [] });
    });
  }, []);

  const changeDir = (dirName) => {
    const nextDir = currentDir ? path.join(currentDir, dirName) : dirName;
    setCurrentDir(nextDir);
    setPathItems([...pathItems, { title: dirName, path: dirName }]);
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
        messageApi.open({ type: "error", content: "Delete failed: " + err.message, duration: 5 });
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
                dispatchFileList({ type: "process", payload: { file, process: percent } });
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
          messageApi.open({ type: "error", content: "Upload failed: " + err.message, duration: 5 });
        });
    });
  };

  const handleSearch = useCallback(
    (keyword, tag) => {
      if (keyword || tag) {
        setSearching(true);
        SearchFiles(keyword, tag).then((data) => {
          dispatch({ type: "list", payload: data?.List || [] });
        });
      } else {
        setSearching(false);
        ReadDir(currentDir).then((data) => {
          dispatch({ type: "list", payload: data?.List || [] });
        });
      }
    },
    [currentDir, dispatch],
  );

  const handleCreateFolder = () => {
    const name = newFolderName.trim();
    if (!name) return;
    setCreatingFolder(true);
    Mkdir(path.join(currentDir, name))
      .then(() => {
        setShowNewFolder(false);
        setNewFolderName("");
        ReadDir(currentDir).then((data) => {
          dispatch({ type: "list", payload: data?.List || [] });
        });
      })
      .catch((err) =>
        messageApi.open({ type: "error", content: "Create folder failed: " + err.message }),
      )
      .finally(() => setCreatingFolder(false));
  };

  const handleTagsUpdated = (itemName, newTags) => {
    dispatch({ type: "updateTags", payload: { name: itemName, tags: newTags } });
  };

  const formatSize = (bytes) => {
    if (!bytes) return "-";
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
    if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(1) + " MB";
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
          return <FolderOutlined style={{ color: "#00d4ff", fontSize: sw < 480 ? 16 : 18, filter: "drop-shadow(0 0 6px rgba(0,212,255,0.3))" }} />;
        }
        return <FileOutlined style={{ color: "#94a3b8", fontSize: sw < 480 ? 16 : 18 }} />;
      },
    },
    {
      title: "Name",
      dataIndex: "Name",
      key: "Name",
      render: (text, record) => {
        if (record.FileType === "dir") {
          return <a onClick={() => changeDir(record.Name)} style={{ fontWeight: 500, color: "#00d4ff" }}>{text}</a>;
        }
        return <a href={path.join(pathPrefix, currentDir, record.Name)} target="_blank" rel="noreferrer" style={{ color: "#e2e8f0" }}>{text}</a>;
      },
    },
    ...(sw >= 480 ? [
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
        width: 100,
        render: (_, record) => {
          if (record.FileType === "dir") return "-";
          return formatSize(record.Size);
        },
      },
    ] : []),
    {
      title: "Tags",
      width: sw < 480 ? 120 : 200,
      render: (_, record) => {
        if (record.FileType === "dir") return null;
        const filepath = currentDir ? path.join(currentDir, record.Name) : record.Name;
        return (
          <div style={{ display: "flex", alignItems: "center", gap: 4 }}>
            <div style={{ display: "flex", flexWrap: "wrap", gap: 2, flex: 1 }}>
              {(record.Tags || []).slice(0, sw < 480 ? 1 : 3).map((t) => (
                <Tag key={t} color="cyan" style={{ margin: 0, fontSize: 10, lineHeight: "16px", padding: "0 4px", borderRadius: 4 }}>{t}</Tag>
              ))}
            </div>
            <Popover
              trigger="click"
              placement={sw < 480 ? "top" : "left"}
              content={<TagEditor filepath={filepath} currentTags={record.Tags} onUpdated={(newTags) => handleTagsUpdated(record.Name, newTags)} />}
            >
              <Button type="text" size="small" icon={<TagsOutlined style={{ color: "#00d4ff", fontSize: 14 }} />} />
            </Popover>
          </div>
        );
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
              <Button type="text" size="small" icon={<DownloadOutlined />} href={path.join(pathPrefix, currentDir, record.Name)} target="_blank" />
            </Tooltip>
            <Tooltip title="Delete">
              <Button type="text" size="small" danger icon={<DeleteOutlined />} onClick={() => deleteFile(currentDir ? path.join(currentDir, record.Name) : record.Name)} />
            </Tooltip>
          </div>
        );
      },
    },
  ];

  const drawerWidth = sw < 768 ? "100%" : "50%";

  return (
    <>
      {contextHolder}
      <div className="tech-page">
        <div className="tech-header" style={{ gap: 12 }}>
          <PathTravel items={pathItems} onClick={changeDirAbsolute} />
          <div style={{ flex: 1 }} />
          {!searching && (
            <Button
              icon={<FolderAddOutlined />}
              size="small"
              style={{ background: "rgba(0,212,255,0.08)", borderColor: "rgba(0,212,255,0.2)", color: "#00d4ff" }}
              onClick={() => setShowNewFolder(true)}
            >
              {sw > 480 && "New Folder"}
            </Button>
          )}
          <Button type="primary" icon={<UploadOutlined />} onClick={() => setShowUploadDrawer(true)}>
            {sw > 480 && "Upload"}
          </Button>
        </div>

        <div className="tech-content">
          <div style={{ marginBottom: 16 }}>
            <SearchBar onSearch={handleSearch} />
          </div>

          <Drawer
            open={showUploadDrawer}
            width={drawerWidth}
            maskClosable={false}
            onClose={() => {
              setShowUploadDrawer(false);
              dispatchFileList({ type: "clear" });
            }}
            extra={<Button type="primary" onClick={uploadFile}>Start Upload</Button>}
          >
            <div style={{ height: "15%" }}>
              <Dragger
                beforeUpload={(file) => {
                  dispatchFileList({ type: "add", payload: file });
                  return false;
                }}
                onRemove={(file) => dispatchFileList({ type: "remove", payload: file })}
                fileList={fileList}
                multiple={true}
                listType="picture"
              >
                <p className="ant-upload-drag-icon"><InboxOutlined /></p>
                <p style={{ color: "#94a3b8" }}>Drag files here or click to select</p>
              </Dragger>
            </div>
          </Drawer>

          <Modal
            open={showNewFolder}
            onCancel={() => { setShowNewFolder(false); setNewFolderName(""); }}
            onOk={handleCreateFolder}
            confirmLoading={creatingFolder}
            title="New Folder"
          >
            <Input placeholder="Folder name" value={newFolderName} onChange={(e) => setNewFolderName(e.target.value)} onPressEnter={handleCreateFolder} />
          </Modal>

          {imageList.length === 0 ? (
            <div className="tech-empty">
              <Empty description={searching ? "No results found" : "Empty directory"} />
            </div>
          ) : (
            <Table
              dataSource={imageList}
              columns={columns}
              rowKey="Name"
              pagination={false}
              size={sw < 768 ? "small" : "middle"}
              scroll={{ x: sw < 480 ? 400 : 650 }}
              onRow={(record) => {
                if (record.FileType === "dir") {
                  return { onClick: () => changeDir(record.Name), style: { cursor: "pointer" } };
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

export default AllFileManage;