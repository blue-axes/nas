import { useState, useEffect, useCallback } from "react";
import { Outlet } from "react-router";
import { useNavigate } from "react-router";
import { Button, ConfigProvider, Layout, Menu, theme, Modal, Form, Input, message, Space, Tag } from "antd";
import {
  DeliveredProcedureOutlined,
  PictureOutlined,
  PlaySquareOutlined,
  FileOutlined,
  CloudServerOutlined,
  MenuOutlined,
  CloseOutlined,
  SettingOutlined,
  KeyOutlined,
  UserOutlined,
} from "@ant-design/icons";
import zh_CN from "antd/locale/zh_CN";
import "antd/dist/reset.css";
import { GetCurrentUser, ChangePassword } from "./apis/User.jsx";

const Sider = Layout.Sider;
const Header = Layout.Header;
const Content = Layout.Content;

function getItem(idx, keyPrefix, item) {
  if (!(item instanceof Object)) {
    return [{}, {}];
  }
  let key = keyPrefix + "key" + idx;
  let children = [];
  let keyLink = {};
  if (item?.children) {
    children = item?.children.map((subItem, subIdx) => {
      let [tmpItem, tmpKeyLink] = getItem(subIdx, key, subItem);
      if (tmpItem) {
        keyLink = { ...keyLink, ...tmpKeyLink };
        return tmpItem;
      }
      return null;
    });
  }

  if (children.length == 0) {
    children = null;
  }

  keyLink[key] = item.link;

  return [
    {
      key: key,
      icon: item?.icon,
      label: item.label,
      title: item?.title || item.label,
      children: children,
      link: item.link,
    },
    keyLink,
  ];
}

function buildMenuConfig(user) {
  const fileMgmt = {
    label: "文件管理",
    title: "文件管理",
    icon: <DeliveredProcedureOutlined />,
    children: [
      {
        label: "图片文件",
        title: "图片文件",
        icon: <PictureOutlined />,
        link: "/image",
      },
      {
        label: "视频文件",
        title: "视频文件",
        icon: <PlaySquareOutlined />,
        link: "/video",
      },
      {
        label: "普通文件",
        title: "普通文件",
        icon: <FileOutlined />,
        link: "/file",
      },
    ],
  };

  const config = [fileMgmt];

  if (user && user.IsAdmin) {
    config.push({
      label: "用户管理",
      title: "用户管理",
      icon: <SettingOutlined />,
      link: "/users",
    });
  }

  return config;
}

function initMenu(user) {
  const menuConfig = buildMenuConfig(user);
  let keyLink = {};
  let res = menuConfig.map((item, idx) => {
    let [tmpItem, tmpKeyLink] = getItem(idx, "", item);
    if (tmpItem) {
      keyLink = { ...keyLink, ...tmpKeyLink };
      return tmpItem;
    }
    return null;
  });

  return [res, keyLink];
}

function App() {
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [siderCollapsed, setSiderCollapsed] = useState(true);
  const [isMobile, setIsMobile] = useState(false);
  const [openKeys, setOpenKeys] = useState(["key0"]);
  const [selectedKeys, setSelectedKeys] = useState(["key0key0"]);
  const [passwordModalOpen, setPasswordModalOpen] = useState(false);
  const [passwordForm] = Form.useForm();

  const [items, keyLink] = initMenu(user);

  const checkMobile = useCallback(() => {
    setIsMobile(window.innerWidth < 768);
  }, []);

  useEffect(() => {
    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, [checkMobile]);

  useEffect(() => {
    GetCurrentUser()
      .then((data) => {
        if (data && data.Username) {
          setUser(data);
        }
      })
      .catch(() => {});
  }, []);

  const handleMenuClick = ({ key }) => {
    const to = keyLink[key];
    if (to) {
      setSelectedKeys([key]);
      navigate(to);
      if (isMobile) {
        setSiderCollapsed(true);
      }
    }
  };

  const handleChangePassword = async () => {
    try {
      const values = await passwordForm.validateFields();
      await ChangePassword(user.Username, values.CurrentPassword, values.NewPassword);
      message.success("Password changed successfully");
      setPasswordModalOpen(false);
      passwordForm.resetFields();
    } catch (e) {
      if (e?.errorFields) return;
      message.error(e.message);
    }
  };

  const themeConfig = {
    algorithm: theme.darkAlgorithm,
    token: {
      colorPrimary: "#00d4ff",
      colorPrimaryBg: "rgba(0, 212, 255, 0.08)",
      colorPrimaryBgHover: "rgba(0, 212, 255, 0.15)",
      colorPrimaryBorder: "rgba(0, 212, 255, 0.3)",
      colorPrimaryHover: "#33ddff",
      colorPrimaryActive: "#00b8e0",
      colorError: "#f43f5e",
      colorErrorBg: "rgba(244, 63, 94, 0.08)",
      colorSuccess: "#22d3ee",
      colorWarning: "#f59e0b",
      colorInfo: "#00d4ff",
      colorTextBase: "#e2e8f0",
      colorBgBase: "#0a0e27",
      colorBorder: "rgba(0, 212, 255, 0.15)",
      colorBorderSecondary: "rgba(0, 212, 255, 0.08)",
      borderRadius: 6,
      wireframe: false,
      fontFamily:
        '"Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
      fontSize: 14,
      controlHeight: 36,
      lineHeight: 1.5715,
    },
    components: {
      Menu: {
        darkItemBg: "transparent",
        darkSubMenuItemBg: "transparent",
        darkItemColor: "#94a3b8",
        darkItemHoverBg: "rgba(0, 212, 255, 0.08)",
        darkItemHoverColor: "#00d4ff",
        darkItemSelectedBg: "rgba(0, 212, 255, 0.12)",
        darkItemSelectedColor: "#00d4ff",
        darkGroupTitleColor: "#64748b",
        itemBorderRadius: 6,
        itemMarginInline: 8,
        itemHeight: 40,
      },
      Layout: {
        siderBg: "#0b1020",
        triggerBg: "#0b1020",
        triggerColor: "#94a3b8",
      },
      Table: {
        headerBg: "rgba(15, 23, 42, 0.9)",
        headerColor: "#94a3b8",
        rowHoverBg: "rgba(0, 212, 255, 0.04)",
        borderColor: "rgba(0, 212, 255, 0.1)",
        colorBgContainer: "transparent",
        colorText: "#e2e8f0",
      },
      Modal: {
        contentBg: "#0f172a",
        headerBg: "#0f172a",
      },
      Drawer: {
        colorBgElevated: "#0f172a",
        colorText: "#e2e8f0",
      },
      Slider: {
        trackBg: "#00d4ff",
        trackHoverBg: "#33ddff",
        railBg: "rgba(255,255,255,0.1)",
        handleColor: "#00d4ff",
        handleActiveColor: "#33ddff",
        dotActiveBorderColor: "#00d4ff",
      },
      Button: {
        defaultBg: "rgba(15, 23, 42, 0.8)",
        defaultBorderColor: "rgba(0, 212, 255, 0.2)",
        defaultColor: "#e2e8f0",
        defaultHoverBg: "rgba(0, 212, 255, 0.08)",
        defaultHoverBorderColor: "rgba(0, 212, 255, 0.4)",
        defaultHoverColor: "#00d4ff",
        defaultActiveBg: "rgba(0, 212, 255, 0.12)",
        dangerColor: "#f43f5e",
        primaryShadow: "0 2px 8px rgba(0, 212, 255, 0.25)",
      },
    },
  };

  const siderStyle = {
    background: "var(--bg-sider)",
    borderRight: "1px solid var(--border-glow)",
  };

  const headerStyle = {
    background: "var(--bg-header)",
    backdropFilter: "blur(20px)",
    WebkitBackdropFilter: "blur(20px)",
    borderBottom: "1px solid var(--border-glow)",
    padding: "0 24px",
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    height: 56,
    flexShrink: 0,
  };

  const mobileHeaderStyle = {
    ...headerStyle,
    padding: "0 12px",
    height: 48,
  };

  const menuBtnStyle = {
    color: "#94a3b8",
    border: "none",
    background: "transparent",
    cursor: "pointer",
    fontSize: "20px",
    padding: "4px",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  };

  const renderSider = () => (
    <>
      <div
        style={{
          padding: "12px 0",
          color: "#64748b",
          fontSize: 11,
          textAlign: "center",
          letterSpacing: 2,
          textTransform: "uppercase",
        }}
      >
        Navigation
      </div>
      <Menu
        mode="inline"
        theme="dark"
        triggerSubMenuAction="click"
        items={items}
        onClick={handleMenuClick}
        openKeys={openKeys}
        selectedKeys={selectedKeys}
        onOpenChange={setOpenKeys}
      />
      <div
        style={{
          position: "absolute",
          bottom: 20,
          left: 0,
          right: 0,
          textAlign: "center",
          color: "#334155",
          fontSize: 11,
        }}
      >
        NAS System
      </div>
    </>
  );

  const renderMobileOverlay = () => (
    <>
      <div
        onClick={() => setSiderCollapsed(true)}
        style={{
          position: "fixed",
          inset: 0,
          background: "rgba(0,0,0,0.5)",
          zIndex: 99,
        }}
      />
      <div
        style={{
          position: "fixed",
          top: 0,
          left: 0,
          bottom: 0,
          width: 200,
          background: "#0b1020",
          borderRight: "1px solid var(--border-glow)",
          zIndex: 100,
          display: "flex",
          flexDirection: "column",
          paddingTop: 12,
        }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            padding: "0 16px 12px",
          }}
        >
          <span
            style={{
              color: "#64748b",
              fontSize: 11,
              letterSpacing: 2,
              textTransform: "uppercase",
            }}
          >
            Navigation
          </span>
          <button
            style={menuBtnStyle}
            onClick={() => setSiderCollapsed(true)}
            aria-label="Close menu"
          >
            <CloseOutlined />
          </button>
        </div>
        <div style={{ flex: 1, overflow: "auto" }}>
          <Menu
            mode="inline"
            theme="dark"
            triggerSubMenuAction="click"
            items={items}
            onClick={handleMenuClick}
            openKeys={openKeys}
            selectedKeys={selectedKeys}
            onOpenChange={setOpenKeys}
          />
          </div>
        <div
          style={{
            padding: "12px 0 20px",
            textAlign: "center",
            color: "#334155",
            fontSize: 11,
          }}
        >
          NAS System
        </div>
      </div>
    </>
  );

  return (
    <>
      <ConfigProvider locale={zh_CN} componentSize="large" theme={themeConfig}>
        <Layout style={{ width: "100vw", height: "100vh" }}>
          <Header style={isMobile ? mobileHeaderStyle : headerStyle}>
            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: isMobile ? 8 : 10,
              }}
            >
              {isMobile && (
                <button
                  style={menuBtnStyle}
                  onClick={() => setSiderCollapsed(false)}
                  aria-label="Open menu"
                >
                  <MenuOutlined />
                </button>
              )}
              <CloudServerOutlined
                style={{
                  fontSize: isMobile ? 18 : 22,
                  color: "var(--accent-cyan)",
                  filter: "drop-shadow(0 0 8px rgba(0,212,255,0.4))",
                }}
              />
              <span className="tech-logo">NAS</span>
            </div>
            {user ? (
              <Space size="small">
                <Tag color="cyan" icon={<UserOutlined />}>
                  {user.Username}
                </Tag>
                <Button
                  icon={<KeyOutlined />}
                  size="small"
                  onClick={() => {
                    passwordForm.resetFields();
                    setPasswordModalOpen(true);
                  }}
                >
                  修改密码
                </Button>
              </Space>
            ) : (
              <span
                style={{
                  color: "#64748b",
                  fontSize: isMobile ? 11 : 13,
                  letterSpacing: 1,
                }}
              >
                v1.0
              </span>
            )}
          </Header>
          <Layout>
            {isMobile ? (
              siderCollapsed ? null : renderMobileOverlay()
            ) : (
              <Sider width={200} style={siderStyle}>
                {renderSider()}
              </Sider>
            )}
            <Content
              style={{
                overflow: "auto",
                background: "var(--bg-primary)",
              }}
            >
              <Outlet />
            </Content>
          </Layout>
        </Layout>
      </ConfigProvider>

      <Modal
        title="修改密码"
        open={passwordModalOpen}
        onOk={handleChangePassword}
        onCancel={() => {
          setPasswordModalOpen(false);
          passwordForm.resetFields();
        }}
        okText="确认"
        cancelText="取消"
      >
        <Form form={passwordForm} layout="vertical">
          <Form.Item
            name="CurrentPassword"
            label="当前密码"
            rules={[{ required: true, message: "请输入当前密码" }]}
          >
            <Input.Password placeholder="当前密码" />
          </Form.Item>
          <Form.Item
            name="NewPassword"
            label="新密码"
            rules={[{ required: true, message: "请输入新密码" }]}
          >
            <Input.Password placeholder="新密码" />
          </Form.Item>
          <Form.Item
            name="ConfirmPassword"
            label="确认新密码"
            dependencies={["NewPassword"]}
            rules={[
              { required: true, message: "请再次输入新密码" },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue("NewPassword") === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error("两次输入的密码不一致"));
                },
              }),
            ]}
          >
            <Input.Password placeholder="确认新密码" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
}

export default App;