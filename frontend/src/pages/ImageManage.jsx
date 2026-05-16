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
  TagsOutlined,
  FolderAddOutlined,
  LeftOutlined,
  RightOutlined,
} from "@ant-design/icons";
import { useImmerReducer } from "use-immer";
import { useEffect, useState, useRef, useCallback } from "react";
import path from "path-browserify";

import { ImageList, FileList } from "./../reducers/ImageManageReducer";
import { ReadDir, UploadFile, DeleteFile, SearchFiles, Mkdir } from "../apis/SimpleUpload";
import PathTravel from "../components/PathTravel";
import SearchBar from "../components/SearchBar";
import TagEditor from "../components/TagEditor";
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
    background: "linear-gradient(transparent 55%, rgba(0,0,0,0.7))",
    opacity: 0,
    transition: "opacity 0.25s ease",
    display: "flex",
    alignItems: "flex-end",
    justifyContent: "space-between",
    padding: "10px",
    pointerEvents: "none",
  },
  imageOverlayVisible: { opacity: 1 },
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
  actionBtn: { pointerEvents: "auto" },
  toolbarBtn: {
    color: "#fff",
    background: "rgba(255,255,255,0.1)",
    border: "1px solid rgba(255,255,255,0.25)",
    borderRadius: 6,
    cursor: "pointer",
    fontSize: "13px",
    padding: "4px 12px",
    display: "flex",
    alignItems: "center",
    gap: 4,
  },
};

function ImageCard({ item, currentDir, onDelete, onTagsUpdated, sw }) {
  const [showLayer, setShowLayer] = useState(false);
  const [tagOpen, setTagOpen] = useState(false);
  const src = path.join(pathPrefix, currentDir, item.Name);
  const filepath = path.join(currentDir, item.Name);

  return (
    <div
      className="tech-card"
      onMouseEnter={() => setShowLayer(true)}
      onMouseLeave={() => { if (!tagOpen) setShowLayer(false); }}
      onTouchStart={() => setShowLayer(true)}
    >
      <div style={styles.imageWrap}>
        <Image
          src={src}
          alt={item.Name}
          width="100%"
          height="100%"
          style={{ objectFit: "cover" }}
        />
      </div>
      <div style={{ ...styles.imageOverlay, ...(showLayer || tagOpen ? styles.imageOverlayVisible : {}) }}>
        <div style={{ flex: 1, minWidth: 0 }}>
          <span style={styles.imageName} title={item.Name}>
            {item.Name}
          </span>
          {(item.Tags || []).length > 0 && (
            <div style={{ display: "flex", flexWrap: "wrap", gap: 2, marginTop: 4 }}>
              {(item.Tags || []).slice(0, 3).map((t) => (
                <Tag key={t} color="cyan" style={{ margin: 0, fontSize: 10, lineHeight: "16px", padding: "0 4px", borderRadius: 4 }}>
                  {t}
                </Tag>
              ))}
            </div>
          )}
        </div>
        <div style={{ display: "flex", gap: 2 }}>
          <Popover
            open={tagOpen}
            onOpenChange={setTagOpen}
            trigger="click"
            placement={sw < 480 ? "top" : "left"}
            content={
              <TagEditor
                filepath={filepath}
                currentTags={item.Tags}
                onUpdated={(newTags) => { onTagsUpdated(item.Name, newTags); }}
              />
            }
          >
            <Tooltip title="Edit Tags">
              <Button style={styles.actionBtn} type="text" size="small" icon={<TagsOutlined style={{ color: "#00d4ff" }} />} />
            </Tooltip>
          </Popover>
          <Tooltip title="Delete">
            <Button
              style={styles.actionBtn}
              type="text"
              size="small"
              danger
              icon={<DeleteOutlined style={{ color: "#f43f5e" }} />}
              onClick={(e) => { e.stopPropagation(); onDelete(filepath); }}
            />
          </Tooltip>
        </div>
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
    { title: <HomeOutlined />, path: "/img/" },
    { title: "img", path: "/img/" },
  ]);
  const [showCount, setShowCount] = useState(PAGE_SIZE);

  const [imageList, dispatch] = useImmerReducer(ImageList, []);
  const [fileList, dispatchFileList] = useImmerReducer(FileList, []);
  const sentinelRef = useRef(null);

  const [searching, setSearching] = useState(false);
  const [showNewFolder, setShowNewFolder] = useState(false);
  const [newFolderName, setNewFolderName] = useState("");
  const [creatingFolder, setCreatingFolder] = useState(false);
  const [previewIndex, setPreviewIndex] = useState(0);
  const [previewVisible, setPreviewVisible] = useState(false);
  const [tagModalVisible, setTagModalVisible] = useState(false);

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
        messageApi.open({ type: "error", content: "Delete failed: " + err.message, duration: 5 });
      });
  };

  const imageFiles = useCallback(() => {
    return imageList.filter((item) => item.FileType !== "dir");
  }, [imageList]);

  const currentPreviewFile = () => {
    const files = imageFiles();
    if (previewIndex >= 0 && previewIndex < files.length) {
      return files[previewIndex];
    }
    return null;
  };

  const handlePreviewDelete = () => {
    const file = currentPreviewFile();
    if (!file) return;
    const filepath = path.join(currentDir, file.Name);
    setPreviewVisible(false);
    setTimeout(() => deleteFile(filepath), 300);
  };

  const handlePreviewTag = () => {
    setPreviewVisible(false);
    setTimeout(() => setTagModalVisible(true), 300);
  };

  const handleTagsUpdated = (itemName, newTags) => {
    dispatch({ type: "updateTags", payload: { name: itemName, tags: newTags } });
    setTagModalVisible(false);
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

  const handleSearch = useCallback((keyword, tag) => {
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
  }, [currentDir, dispatch]);

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
      .catch((err) => messageApi.open({ type: "error", content: "Create folder failed: " + err.message }))
      .finally(() => setCreatingFolder(false));
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
        onTagsUpdated={handleTagsUpdated}
        sw={sw}
      />
    );
  });

  const imgGridCols =
    sw < 360 ? "repeat(1, 1fr)" : sw < 480 ? "repeat(2, 1fr)" : sw < 768 ? "repeat(auto-fill, minmax(150px, 1fr))" : "repeat(auto-fill, minmax(220px, 1fr))";
  const drawerWidth = sw < 768 ? "100%" : "50%";

  const items = allItems.slice(0, showCount);

  const toolbarRender = () => {
    const total = imageFiles().length;
    const cur = previewIndex + 1;
    return (
      <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
        <span style={{ color: "rgba(255,255,255,0.65)", fontSize: 13 }}>
          {cur} / {total}
        </span>
        <div style={{ display: "flex", gap: 6 }}>
          <button
            style={styles.toolbarBtn}
            onClick={() => setPreviewIndex((p) => Math.max(0, p - 1))}
            title="Previous (Left Arrow)"
          >
            <LeftOutlined style={{ fontSize: 12 }} /> Prev
          </button>
          <button
            style={styles.toolbarBtn}
            onClick={() => setPreviewIndex((p) => Math.min(total - 1, p + 1))}
            title="Next (Right Arrow)"
          >
            Next <RightOutlined style={{ fontSize: 12 }} />
          </button>
        </div>
        <div style={{ width: 1, height: 20, background: "rgba(255,255,255,0.2)" }} />
        <button style={{ ...styles.toolbarBtn, borderColor: "rgba(0, 212, 255, 0.4)" }} onClick={handlePreviewTag}>
          <TagsOutlined style={{ color: "#00d4ff" }} /> Tag
        </button>
        <button style={{ ...styles.toolbarBtn, borderColor: "rgba(244, 63, 94, 0.4)", color: "#f43f5e" }} onClick={handlePreviewDelete}>
          <DeleteOutlined /> Delete
        </button>
      </div>
    );
  };

  const previewFile = currentPreviewFile();
  const previewFilepath = previewFile
    ? path.join(currentDir, previewFile.Name)
    : "";

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
            onClose={() => { setShowUploadDrawer(false); dispatchFileList({ type: "clear" }); }}
            extra={<Button type="primary" onClick={uploadFile}>Start Upload</Button>}
          >
            <div style={{ height: "15%" }}>
              <Dragger
                beforeUpload={(file) => { dispatchFileList({ type: "add", payload: file }); return false; }}
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
            <Input
              placeholder="Folder name"
              value={newFolderName}
              onChange={(e) => setNewFolderName(e.target.value)}
              onPressEnter={handleCreateFolder}
            />
          </Modal>

          <Modal
            open={tagModalVisible}
            onCancel={() => setTagModalVisible(false)}
            footer={null}
            title={previewFile?.Name || "Edit Tags"}
          >
            {previewFile && (
              <TagEditor
                filepath={previewFilepath}
                currentTags={previewFile.Tags || []}
                onUpdated={(newTags) => handleTagsUpdated(previewFile.Name, newTags)}
              />
            )}
          </Modal>

          {items.length === 0 ? (
            <div className="tech-empty">
              <Empty description={searching ? "No results found" : "Empty directory"} />
            </div>
          ) : (
            <Image.PreviewGroup
              preview={{
                visible: previewVisible,
                current: previewIndex,
                onChange: (current) => setPreviewIndex(current),
                onVisibleChange: (visible) => {
                  setPreviewVisible(visible);
                  if (!visible) setPreviewIndex(0);
                },
                toolbarRender,
              }}
            >
              <div className="tech-grid" style={{ gridTemplateColumns: imgGridCols }}>
                {items}
                {!searching && showCount < allItems.length && (
                  <div ref={sentinelRef} className="tech-sentinel"><Spin /></div>
                )}
              </div>
            </Image.PreviewGroup>
          )}
        </div>
      </div>
    </>
  );
}

export default ImageManage;