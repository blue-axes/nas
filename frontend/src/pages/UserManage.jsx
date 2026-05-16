import { useState, useEffect, useCallback } from "react";
import { Button, Table, Modal, Form, Input, Switch, Popconfirm, message, Space, Tag, Tooltip } from "antd";
import { PlusOutlined, DeleteOutlined, EditOutlined, UserOutlined, CrownOutlined } from "@ant-design/icons";
import { ListUsers, CreateUser, UpdateUser, DeleteUser } from "../apis/User.jsx";
import useScreenWidth from "../hooks/useScreenWidth.js";

export default function UserManage() {
  const sw = useScreenWidth();
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [form] = Form.useForm();

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    try {
      const data = await ListUsers();
      setUsers(data || []);
    } catch (e) {
      message.error("Failed to load users: " + e.message);
    }
    setLoading(false);
  }, []);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const handleCreate = () => {
    setEditingUser(null);
    form.resetFields();
    form.setFieldsValue({ CanRead: true, CanWrite: false });
    setModalOpen(true);
  };

  const handleEdit = (record) => {
    setEditingUser(record);
    form.setFieldsValue({
      Username: record.Username,
      Password: "",
      CanRead: record.CanRead,
      CanWrite: record.CanWrite,
    });
    setModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingUser) {
        await UpdateUser(editingUser.Username, values.CanRead, values.CanWrite);
        message.success("User updated");
      } else {
        await CreateUser(values.Username, values.Password, values.CanRead, values.CanWrite);
        message.success("User created");
      }
      setModalOpen(false);
      fetchUsers();
    } catch (e) {
      if (e?.errorFields) return;
      message.error(e.message);
    }
  };

  const handleDelete = async (record) => {
    try {
      await DeleteUser(record.Username);
      message.success("User deleted");
      fetchUsers();
    } catch (e) {
      message.error(e.message);
    }
  };

  const columns = [
    {
      title: "Username",
      dataIndex: "Username",
      key: "Username",
      render: (text, record) => (
        <Space wrap size={4}>
          {record.IsAdmin ? <CrownOutlined style={{ color: "#f59e0b" }} /> : <UserOutlined />}
          <span>{text}</span>
          {record.IsAdmin && <Tag color="gold">Admin</Tag>}
        </Space>
      ),
    },
    ...(sw >= 480 ? [
      {
        title: "Read",
        dataIndex: "CanRead",
        key: "CanRead",
        render: (v) => (v ? <Tag color="cyan">Yes</Tag> : <Tag>No</Tag>),
      },
      {
        title: "Write",
        dataIndex: "CanWrite",
        key: "CanWrite",
        render: (v) => (v ? <Tag color="cyan">Yes</Tag> : <Tag>No</Tag>),
      },
    ] : []),
    {
      title: "Actions",
      key: "actions",
      render: (_, record) => {
        if (record.IsAdmin) return null;
        return (
        <Space size={4}>
          <Tooltip title="Edit">
            <Button icon={<EditOutlined />} size="small" onClick={() => handleEdit(record)}>
              {sw >= 600 && "Edit"}
            </Button>
          </Tooltip>
          <Popconfirm title="Delete this user?" onConfirm={() => handleDelete(record)} placement="topRight">
            <Tooltip title="Delete">
              <Button icon={<DeleteOutlined />} size="small" danger>
                {sw >= 600 && "Delete"}
              </Button>
            </Tooltip>
          </Popconfirm>
        </Space>
        );
      },
    },
  ];

  return (
    <div style={{ padding: sw < 768 ? 12 : 24, height: "100%", display: "flex", flexDirection: "column" }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16, flexWrap: "wrap", gap: 8 }}>
        <h2 style={{ color: "var(--accent-cyan)", margin: 0, fontSize: sw < 480 ? 16 : 20, fontWeight: 600 }}>User Management</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          {sw >= 480 && "Create User"}
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={users}
        rowKey="Username"
        loading={loading}
        style={{ flex: 1 }}
        pagination={false}
        size={sw < 768 ? "small" : "middle"}
        scroll={{ x: sw < 480 ? 350 : 500 }}
      />
      <Modal
        title={editingUser ? "Edit User" : "Create User"}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        okText={editingUser ? "Save" : "Create"}
        width={sw < 480 ? "95%" : undefined}
        centered={sw < 480}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="Username"
            label="Username"
            rules={[{ required: true, message: "Username is required" }]}
          >
            <Input disabled={!!editingUser} placeholder="username" />
          </Form.Item>
          <Form.Item
            name="Password"
            label={editingUser ? "New Password (leave empty to keep)" : "Password"}
            rules={editingUser ? [] : [{ required: true, message: "Password is required" }]}
          >
            <Input.Password placeholder="password" />
          </Form.Item>
          <Form.Item name="CanRead" label="Read Permission" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="CanWrite" label="Write Permission" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}