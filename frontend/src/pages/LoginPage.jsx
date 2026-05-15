import { useState } from "react";
import { Form, Input, Button, message, Typography } from "antd";
import { UserOutlined, LockOutlined, CloudServerOutlined } from "@ant-design/icons";
import { Login } from "../apis/User.jsx";

export default function LoginPage({ onLoginSuccess }) {
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();

  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      const data = await Login(values.Username, values.Password);
      message.success("登录成功");
      onLoginSuccess(data);
    } catch (e) {
      message.error(e.message);
    }
    setLoading(false);
  };

  return (
    <div
      style={{
        width: "100vw",
        height: "100vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        background: "var(--bg-primary)",
      }}
    >
      <div
        style={{
          width: 360,
          padding: "40px 32px",
          background: "var(--bg-glass)",
          borderRadius: 12,
          border: "1px solid var(--border-glow)",
          boxShadow: "0 0 40px rgba(0, 212, 255, 0.08)",
        }}
      >
        <div style={{ textAlign: "center", marginBottom: 32 }}>
          <CloudServerOutlined
            style={{
              fontSize: 48,
              color: "var(--accent-cyan)",
              filter: "drop-shadow(0 0 12px rgba(0,212,255,0.4))",
            }}
          />
          <Typography.Title level={3} style={{ color: "var(--accent-cyan)", marginBottom: 4 }}>
            NAS
          </Typography.Title>
          <Typography.Text type="secondary">请登录以继续</Typography.Text>
        </div>
        <Form form={form} onFinish={handleSubmit} size="large">
          <Form.Item name="Username" rules={[{ required: true, message: "请输入用户名" }]}>
            <Input prefix={<UserOutlined />} placeholder="用户名" />
          </Form.Item>
          <Form.Item name="Password" rules={[{ required: true, message: "请输入密码" }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              登 录
            </Button>
          </Form.Item>
        </Form>
      </div>
    </div>
  );
}