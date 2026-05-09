import {
  Layout,
  Image,
  Button,
  Upload,
  Drawer,
  message,
  Spin,
  Tooltip,
  Empty,
} from "antd";
import {
  UploadOutlined,
  FolderOutlined,
  HomeOutlined,
  InboxOutlined,
  DeleteOutlined,
} from "@ant-design/icons";
import { useImmerReducer } from "use-immer";
import { useEffect, useState, useRef, useCallback } from "react";
import path from "path-browserify";

import { ImageList, FileList } from "./../reducers/ImageManageReducer";
import { ReadDir, UploadFile, DeleteFile } from "../apis/SimpleUpload";
import PathTravel from "../components/PathTravel";
import useScreenWidth from "../hooks/useScreenWidth";

const Header = Layout.Header;
const Content = Layout.Content;
const pathPrefix = "/simple_upload/object";
const { Dragger } = Upload;
const PAGE_SIZE = 20;

const styles = {
  imageWrap: {
    position: "relative",
    aspectRatio: "1 / 1",
    overflow: "hidden",
    background: "rgba(10, 14, 39, 0.95)",
  },
  imageOverlay: {
    position: "absolute",
    inset: 0,
    background:
      "linear-gradient(transparent 55%, rgba(0,0,0,0.7))",
    opacity: 0,
    transition: "opacity 0.25s ease",
    display: "flex",
    alignItems: "flex-end",
    justifyContent: "space-between",
    padding: "10px",
    pointerEvents: "none",
  },
  imageOverlayVisible: {
    opacity: 1,
  },
  imageName: {
    color: "#e2e8f0",
    fontSize: "13px",
    fontWeight: 500,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
    textShadow: "0 1px 4px rgba(0,0,0,0.6)",
    maxWidth: "75%",
  },
  deleteBtn: {
    pointerEvents: "auto",
  },
};

function LazyImage({ src, name }) {
  const containerRef = useRef(null);
  const [inView, setInView] = useState(false);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setInView(true);
          observer.unobserve(el);
        }
      },
      { rootMargin: "200px" },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  return (
    <div ref={containerRef} style={styles.imageWrap}>
      {inView ? (
        <Image
          src={src}
          alt={name}
          width="100%"
          height="100%"
          style={{ objectFit: "cover" }}
          preview={{ mask: "Preview" }}
        />
      ) : (
        <div
          style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            height: "100%",
          }}
        >
          <Spin size="small" />
        </div>
      )}
    </div>
  );
}

function ImageCard({ item, currentDir, onDelete }) {
  const [hover, setHover] = useState(false);
  const src = path.join(pathPrefix, currentDir, item.Name);

  return (
    <div
      className="tech-card"
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
    >
      <LazyImage src={src} name={item.Name} />
      <div
        style={{
          ...styles.imageOverlay,
          ...(hover ? styles.imageOverlayVisible : {}),
        }}
      >
        <span style={styles.imageName} title={item.Name}>
          {item.Name}
        </span>
        <Tooltip title="Delete">
          <Button
            style={styles.deleteBtn}
            type="text"
            size="small"
            danger
            icon={<DeleteOutlined style={{ color: "#f43f5e" }} />}
            onClick={(e) => {
              e.stopPropagation();
              onDelete(path.join(currentDir, item.Name));
            }}
          />
        </Tooltip>
      </div>
    </div>
  );
}

function FolderCard({ item, onEnter }) {
  return (
    <div className="tech-card" onClick={() => onEnter(item.Name)}>
      <FolderOutlined className="tech-folder-icon" />
      <div className="tech-name" title={item.Name}>
        {item.Name}
      </div>
    </div>
  );
}

function ImageManage() {
  const sw = useScreenWidth();
  const [messageApi, contextHolder] = message.useMessage();

  const [currentDir, setCurrentDir] = useState("/img/");
  const [showUploadDrawer, setShowUploadDrawer] = useState(false);
  const [pathItems, setPathItems] = useState([
    { title: <HomeOutlined />, path: "/" },
    { title: "img", path: "/img/" },
  ]);
  const [showCount, setShowCount] = useState(PAGE_SIZE);

  const [imageList, dispatch] = useImmerReducer(ImageList, []);
  const [fileList, dispatchFileList] = useImmerReducer(FileList, []);
  const sentinelRef = useRef(null);

  const loadMore = useCallback(() => {
    setShowCount((prev) => prev + PAGE_SIZE);
  }, []);

  useEffect(() => {
    const el = sentinelRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) loadMore();
      },
      { rootMargin: "300px" },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, [loadMore, imageList.length]);

  useEffect(() => {
    ReadDir(currentDir).then((data) => {
      setShowCount(PAGE_SIZE);
      dispatch({ type: "list", payload: data.List ? data.List : [] });
    });
  }, []);

  const changeDir = (dirName) => {
    const nextDir = path.normalize(path.join(currentDir, dirName));
    setCurrentDir(nextDir);
    setPathItems([...pathItems, { title: dirName, path: nextDir }]);
    setShowCount(PAGE_SIZE);
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
    setShowCount(PAGE_SIZE);
    ReadDir(absolutePath).then((data) => {
      dispatch({ type: "list", payload: data?.List });
    });
  };

  const deleteFile = (filepath) => {
    DeleteFile(filepath)
      .then(() => {
        dispatch({ type: "remove", payload: path.basename(filepath) });
      })
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

  const allItems = imageList.map((item) => {
    if (item.FileType == "dir") {
      return <FolderCard key={item.Name} item={item} onEnter={changeDir} />;
    }
    return (
      <ImageCard
        key={item.Name}
        item={item}
        currentDir={currentDir}
        onDelete={deleteFile}
      />
    );
  });

  const imgGridCols =
    sw < 480 ? "repeat(2, 1fr)" : sw < 768 ? "repeat(auto-fill, minmax(150px, 1fr))" : "repeat(auto-fill, minmax(220px, 1fr))";
  const drawerWidth = sw < 768 ? "100%" : "50%";

  const items = allItems.slice(0, showCount);

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

          {items.length === 0 ? (
            <div className="tech-empty">
              <Empty description="Empty directory" />
            </div>
          ) : (
            <div
              className="tech-grid"
              style={{
                gridTemplateColumns: imgGridCols,
              }}
            >
              {items}
              {showCount < allItems.length && (
                <div ref={sentinelRef} className="tech-sentinel">
                  <Spin />
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </>
  );
}

export default ImageManage;