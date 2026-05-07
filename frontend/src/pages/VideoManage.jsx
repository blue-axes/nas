import {
  Layout,
  Button,
  Upload,
  Drawer,
  message,
  Spin,
  Tooltip,
  Empty,
  Modal,
  Slider,
} from "antd";
import {
  UploadOutlined,
  FolderOutlined,
  HomeOutlined,
  InboxOutlined,
  DeleteOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  ForwardOutlined,
  BackwardOutlined,
  ExpandOutlined,
  CompressOutlined,
} from "@ant-design/icons";
import { useImmerReducer } from "use-immer";
import { useEffect, useState, useRef, useCallback } from "react";
import path from "path-browserify";

import { ImageList, FileList } from "./../reducers/ImageManageReducer";
import { ReadDir, UploadFile, DeleteFile } from "../apis/SimpleUpload";
import PathTravel from "../components/PathTravel";

const Content = Layout.Content;
const pathPrefix = "/simple_upload/object";
const { Dragger } = Upload;
const PAGE_SIZE = 20;

const SPEED_OPTIONS = [0.5, 0.75, 1, 1.25, 1.5, 2, 3];
const SKIP_SECONDS = 10;

const styles = {
  page: {
    height: "100%",
    display: "flex",
    flexDirection: "column",
    background: "#f5f5f5",
  },
  header: {
    background: "#fff",
    padding: "12px 24px",
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    borderBottom: "1px solid #f0f0f0",
    boxShadow: "0 1px 4px rgba(0,0,0,0.04)",
  },
  content: {
    flex: 1,
    overflow: "auto",
    padding: "20px 24px",
  },
  grid: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))",
    gap: "16px",
  },
  card: {
    position: "relative",
    borderRadius: "10px",
    overflow: "hidden",
    background: "#fff",
    boxShadow: "0 1px 3px rgba(0,0,0,0.08)",
    transition: "transform 0.2s ease, box-shadow 0.2s ease",
    cursor: "pointer",
  },
  cardHover: {
    transform: "translateY(-3px)",
    boxShadow: "0 6px 16px rgba(0,0,0,0.12)",
  },
  folderIcon: {
    fontSize: "72px",
    color: "#faad14",
    display: "block",
    lineHeight: "180px",
    textAlign: "center",
  },
  folderName: {
    padding: "10px 12px",
    fontSize: "14px",
    fontWeight: 500,
    textAlign: "center",
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
    borderTop: "1px solid #f5f5f5",
    color: "#333",
  },
  videoThumb: {
    position: "relative",
    aspectRatio: "16 / 9",
    overflow: "hidden",
    background: "#000",
  },
  videoThumbEl: {
    width: "100%",
    height: "100%",
    objectFit: "cover",
  },
  playOverlay: {
    position: "absolute",
    inset: 0,
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    background: "rgba(0,0,0,0.3)",
    transition: "background 0.2s ease",
  },
  playIcon: {
    fontSize: "48px",
    color: "rgba(255,255,255,0.9)",
    filter: "drop-shadow(0 2px 8px rgba(0,0,0,0.3))",
  },
  videoInfo: {
    padding: "10px 12px",
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
  },
  videoName: {
    fontSize: "14px",
    fontWeight: 500,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
    color: "#333",
    flex: 1,
    marginRight: "8px",
  },
  sentinel: {
    width: "100%",
    padding: "24px",
    textAlign: "center",
    gridColumn: "1 / -1",
  },
  emptyWrap: {
    gridColumn: "1 / -1",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    minHeight: "300px",
  },
  playerContainer: {
    width: "100%",
    background: "#000",
    position: "relative",
    borderRadius: "8px",
    overflow: "hidden",
  },
  playerVideo: {
    width: "100%",
    display: "block",
    maxHeight: "70vh",
  },
  controls: {
    position: "absolute",
    bottom: 0,
    left: 0,
    right: 0,
    background: "linear-gradient(transparent, rgba(0,0,0,0.8))",
    padding: "20px 16px 12px",
    display: "flex",
    flexDirection: "column",
    gap: "8px",
    transition: "opacity 0.3s ease",
  },
  controlsRow: {
    display: "flex",
    alignItems: "center",
    gap: "8px",
  },
  controlBtn: {
    color: "#fff",
    border: "none",
    background: "transparent",
    cursor: "pointer",
    fontSize: "16px",
    padding: "4px 8px",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
  timeText: {
    color: "rgba(255,255,255,0.85)",
    fontSize: "13px",
    fontFamily: "monospace",
    minWidth: "100px",
    textAlign: "center",
  },
  speedBtn: {
    color: "#fff",
    border: "1px solid rgba(255,255,255,0.3)",
    borderRadius: "4px",
    background: "transparent",
    cursor: "pointer",
    fontSize: "12px",
    padding: "2px 8px",
    fontWeight: 500,
  },
  speedBtnActive: {
    background: "rgba(255,255,255,0.2)",
    borderColor: "rgba(255,255,255,0.6)",
  },
  progressWrap: {
    flex: 1,
    cursor: "pointer",
    padding: "4px 0",
  },
};

function formatTime(seconds) {
  if (!seconds || isNaN(seconds)) return "00:00";
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);
  if (h > 0) {
    return `${h}:${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
  }
  return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
}

function VideoPlayer({ src, open, onClose }) {
  const videoRef = useRef(null);
  const [playing, setPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [speed, setSpeed] = useState(1);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const containerRef = useRef(null);
  const controlsTimerRef = useRef(null);
  const [showControls, setShowControls] = useState(true);

  useEffect(() => {
    if (!open) {
      setPlaying(false);
      setCurrentTime(0);
      setDuration(0);
      setSpeed(1);
    }
  }, [open]);

  const startControlsTimer = useCallback(() => {
    setShowControls(true);
    if (controlsTimerRef.current) clearTimeout(controlsTimerRef.current);
    controlsTimerRef.current = setTimeout(() => {
      if (playing) setShowControls(false);
    }, 3000);
  }, [playing]);

  useEffect(() => {
    if (playing) {
      startControlsTimer();
    } else {
      setShowControls(true);
      if (controlsTimerRef.current) clearTimeout(controlsTimerRef.current);
    }
    return () => {
      if (controlsTimerRef.current) clearTimeout(controlsTimerRef.current);
    };
  }, [playing, startControlsTimer]);

  const togglePlay = () => {
    const video = videoRef.current;
    if (!video) return;
    if (video.paused) {
      video.play();
      setPlaying(true);
    } else {
      video.pause();
      setPlaying(false);
    }
  };

  const skip = (seconds) => {
    const video = videoRef.current;
    if (!video) return;
    video.currentTime = Math.max(0, Math.min(video.duration, video.currentTime + seconds));
  };

  const changeSpeed = (s) => {
    const video = videoRef.current;
    if (!video) return;
    video.playbackRate = s;
    setSpeed(s);
  };

  const onTimeUpdate = () => {
    const video = videoRef.current;
    if (!video) return;
    setCurrentTime(video.currentTime);
  };

  const onLoadedMetadata = () => {
    const video = videoRef.current;
    if (!video) return;
    setDuration(video.duration);
  };

  const onProgressChange = (value) => {
    const video = videoRef.current;
    if (!video) return;
    video.currentTime = value;
    setCurrentTime(value);
  };

  const toggleFullscreen = () => {
    const container = containerRef.current;
    if (!container) return;
    if (!document.fullscreenElement) {
      container.requestFullscreen();
      setIsFullscreen(true);
    } else {
      document.exitFullscreen();
      setIsFullscreen(false);
    }
  };

  useEffect(() => {
    const handler = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };
    document.addEventListener("fullscreenchange", handler);
    return () => document.removeEventListener("fullscreenchange", handler);
  }, []);

  const handleKeyDown = (e) => {
    if (e.key === " " || e.key === "k") {
      e.preventDefault();
      togglePlay();
    } else if (e.key === "ArrowRight") {
      e.preventDefault();
      skip(SKIP_SECONDS);
    } else if (e.key === "ArrowLeft") {
      e.preventDefault();
      skip(-SKIP_SECONDS);
    } else if (e.key === "f") {
      e.preventDefault();
      toggleFullscreen();
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      const idx = SPEED_OPTIONS.indexOf(speed);
      if (idx < SPEED_OPTIONS.length - 1) changeSpeed(SPEED_OPTIONS[idx + 1]);
    } else if (e.key === "ArrowDown") {
      e.preventDefault();
      const idx = SPEED_OPTIONS.indexOf(speed);
      if (idx > 0) changeSpeed(SPEED_OPTIONS[idx - 1]);
    }
  };

  return (
    <Modal
      open={open}
      onCancel={onClose}
      footer={null}
      width="80%"
      styles={{ body: { padding: 0, background: "#000" } }}
      destroyOnClose
    >
      <div
        ref={containerRef}
        style={styles.playerContainer}
        onKeyDown={handleKeyDown}
        tabIndex={0}
        onMouseMove={startControlsTimer}
        onClick={togglePlay}
      >
        <video
          ref={videoRef}
          src={src}
          style={styles.playerVideo}
          onTimeUpdate={onTimeUpdate}
          onLoadedMetadata={onLoadedMetadata}
          onEnded={() => setPlaying(false)}
          onClick={(e) => e.stopPropagation()}
        />
        <div
          style={{
            ...styles.controls,
            opacity: showControls ? 1 : 0,
            pointerEvents: showControls ? "auto" : "none",
          }}
          onClick={(e) => e.stopPropagation()}
        >
          <div style={styles.progressWrap}>
            <Slider
              min={0}
              max={duration || 0}
              value={currentTime}
              onChange={onProgressChange}
              tooltip={{ formatter: formatTime }}
              styles={{
                track: { background: "#1677ff" },
                rail: { background: "rgba(255,255,255,0.2)" },
                handle: { borderColor: "#1677ff" },
              }}
            />
          </div>
          <div style={styles.controlsRow}>
            <button style={styles.controlBtn} onClick={togglePlay} title="Play/Pause (Space)">
              {playing
                ? <PauseCircleOutlined style={{ fontSize: "24px" }} />
                : <PlayCircleOutlined style={{ fontSize: "24px" }} />}
            </button>
            <button style={styles.controlBtn} onClick={() => skip(-SKIP_SECONDS)} title={`Rewind ${SKIP_SECONDS}s (Left Arrow)`}>
              <BackwardOutlined />
            </button>
            <button style={styles.controlBtn} onClick={() => skip(SKIP_SECONDS)} title={`Forward ${SKIP_SECONDS}s (Right Arrow)`}>
              <ForwardOutlined />
            </button>
            <span style={styles.timeText}>
              {formatTime(currentTime)} / {formatTime(duration)}
            </span>
            <div style={{ flex: 1 }} />
            {SPEED_OPTIONS.map((s) => (
              <button
                key={s}
                style={{
                  ...styles.speedBtn,
                  ...(speed === s ? styles.speedBtnActive : {}),
                }}
                onClick={() => changeSpeed(s)}
                title={`Speed ${s}x (Up/Down Arrow)`}
              >
                {s}x
              </button>
            ))}
            <div style={{ flex: 1 }} />
            <button style={styles.controlBtn} onClick={toggleFullscreen} title="Fullscreen (F)">
              {isFullscreen ? <CompressOutlined /> : <ExpandOutlined />}
            </button>
          </div>
        </div>
      </div>
    </Modal>
  );
}

function VideoCard({ item, currentDir, onDelete, onPlay }) {
  const [hover, setHover] = useState(false);
  const src = path.join(pathPrefix, currentDir, item.Name);

  return (
    <div
      style={{
        ...styles.card,
        ...(hover ? styles.cardHover : {}),
      }}
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
    >
      <div style={styles.videoThumb} onClick={() => onPlay(item.Name)}>
        <video
          src={src}
          style={styles.videoThumbEl}
          preload="metadata"
          muted
        />
        <div
          style={{
            ...styles.playOverlay,
            background: hover ? "rgba(0,0,0,0.5)" : "rgba(0,0,0,0.3)",
          }}
        >
          <PlayCircleOutlined style={styles.playIcon} />
        </div>
      </div>
      <div style={styles.videoInfo}>
        <span style={styles.videoName} title={item.Name}>
          {item.Name}
        </span>
        <Tooltip title="Delete">
          <Button
            type="text"
            size="small"
            danger
            icon={<DeleteOutlined />}
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
  const [hover, setHover] = useState(false);

  return (
    <div
      style={{
        ...styles.card,
        ...(hover ? styles.cardHover : {}),
      }}
      onClick={() => onEnter(item.Name)}
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
    >
      <FolderOutlined style={styles.folderIcon} />
      <div style={styles.folderName} title={item.Name}>
        {item.Name}
      </div>
    </div>
  );
}

function VideoManage() {
  const [messageApi, contextHolder] = message.useMessage();

  const [currentDir, setCurrentDir] = useState("/video/");
  const [showUploadDrawer, setShowUploadDrawer] = useState(false);
  const [pathItems, setPathItems] = useState([
    {
      title: <HomeOutlined />,
      path: "/",
    },
    {
      title: "video",
      path: "/video/",
    },
  ]);
  const [showCount, setShowCount] = useState(PAGE_SIZE);

  const [imageList, dispatch] = useImmerReducer(ImageList, []);
  const [fileList, dispatchFileList] = useImmerReducer(FileList, []);
  const sentinelRef = useRef(null);

  const [playerOpen, setPlayerOpen] = useState(false);
  const [playerSrc, setPlayerSrc] = useState("");

  const loadMore = useCallback(() => {
    setShowCount((prev) => prev + PAGE_SIZE);
  }, []);

  useEffect(() => {
    const el = sentinelRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          loadMore();
        }
      },
      { rootMargin: "300px" },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, [loadMore, imageList.length]);

  useEffect(() => {
    ReadDir(currentDir).then((data) => {
      setShowCount(PAGE_SIZE);
      dispatch({
        type: "list",
        payload: data.List ? data.List : [],
      });
    });
  }, []);

  const changeDir = (dirName) => {
    const nextDir = path.normalize(path.join(currentDir, dirName));
    setCurrentDir(nextDir);
    setPathItems([
      ...pathItems,
      {
        title: dirName,
        path: nextDir,
      },
    ]);
    setShowCount(PAGE_SIZE);
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
    let absPath = "";
    const newPathItems = [];

    for (const item of pathItems) {
      absPath = path.join(absPath, item.path);
      newPathItems.push(item);
      if (absPath === absolutePath) {
        break;
      }
    }
    setPathItems(newPathItems);
    setShowCount(PAGE_SIZE);

    ReadDir(absolutePath).then((data) => {
      dispatch({
        type: "list",
        payload: data?.List,
      });
    });
  };

  const deleteFile = (filepath) => {
    DeleteFile(filepath)
      .then(() => {
        dispatch({
          type: "remove",
          payload: path.basename(filepath),
        });
      })
      .catch((err) => {
        messageApi.open({
          type: "error",
          content: "Delete failed: " + err.message,
          duration: 5,
        });
      });
  };

  const openPlayer = (fileName) => {
    setPlayerSrc(path.join(pathPrefix, currentDir, fileName));
    setPlayerOpen(true);
  };

  const uploadFile = () => {
    if (fileList.length == 0) {
      messageApi.open({
        type: "error",
        content: "Please select files",
      });
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
                if (percent > 100) {
                  percent = 100;
                }
                dispatchFileList({
                  type: "process",
                  payload: {
                    file,
                    process: percent,
                  },
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
      dispatchFileList({
        type: "process",
        payload: {
          file,
          process: 0,
        },
      });
      UploadFile(path.join(currentDir, file.name), formData)
        .then(() => {
          dispatchFileList({
            type: "done",
            payload: {
              file,
            },
          });
        })
        .catch((err) => {
          console.log(err);
          dispatchFileList({
            type: "error",
            payload: {
              file,
            },
          });
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
      <VideoCard
        key={item.Name}
        item={item}
        currentDir={currentDir}
        onDelete={deleteFile}
        onPlay={openPlayer}
      />
    );
  });

  const items = allItems.slice(0, showCount);

  return (
    <>
      {contextHolder}
      <div style={styles.page}>
        <div style={styles.header}>
          <PathTravel items={pathItems} onClick={changeDirAbsolute} />
          <Button
            type="primary"
            icon={<UploadOutlined />}
            onClick={() => setShowUploadDrawer(true)}
          >
            Upload
          </Button>
        </div>

        <div style={styles.content}>
          <Drawer
            open={showUploadDrawer}
            width="50%"
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
                onRemove={(file) => {
                  dispatchFileList({ type: "remove", payload: file });
                }}
                fileList={fileList}
                multiple={true}
                listType="picture"
                accept="video/*"
              >
                <p className="ant-upload-drag-icon">
                  <InboxOutlined />
                </p>
                <p style={{ color: "#999" }}>
                  Drag video files here or click to select
                </p>
              </Dragger>
            </div>
          </Drawer>

          {items.length === 0 ? (
            <div style={styles.emptyWrap}>
              <Empty description="Empty directory" />
            </div>
          ) : (
            <div style={styles.grid}>
              {items}
              {showCount < allItems.length && (
                <div ref={sentinelRef} style={styles.sentinel}>
                  <Spin />
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      <VideoPlayer
        src={playerSrc}
        open={playerOpen}
        onClose={() => setPlayerOpen(false)}
      />
    </>
  );
}

export default VideoManage;
