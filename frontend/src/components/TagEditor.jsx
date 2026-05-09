import { useState } from "react";
import { Button, Input, Tag, message } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { UpdateFileTags } from "../apis/SimpleUpload";

export default function TagEditor({ filepath, currentTags, onUpdated }) {
  const [input, setInput] = useState("");
  const [saving, setSaving] = useState(false);

  const saveTags = (newTags) => {
    setSaving(true);
    UpdateFileTags(filepath, newTags)
      .then(() => onUpdated(newTags))
      .catch((err) => message.error("Failed to update tags: " + err.message))
      .finally(() => setSaving(false));
  };

  const addTag = () => {
    const t = input.trim();
    if (!t || (currentTags || []).includes(t)) return;
    saveTags([...(currentTags || []), t]);
    setInput("");
  };

  return (
    <div style={{ minWidth: 180 }}>
      <div style={{ display: "flex", gap: 4, marginBottom: 8 }}>
        <Input
          size="small"
          placeholder="Add tag"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onPressEnter={addTag}
          style={{ flex: 1 }}
        />
        <Button size="small" icon={<PlusOutlined />} onClick={addTag} loading={saving} />
      </div>
      <div style={{ display: "flex", flexWrap: "wrap", gap: 4 }}>
        {(currentTags || []).length === 0 && (
          <span style={{ color: "#64748b", fontSize: 12 }}>No tags</span>
        )}
        {(currentTags || []).map((t) => (
          <Tag
            key={t}
            closable
            onClose={() => saveTags((currentTags || []).filter((x) => x !== t))}
            color="cyan"
          >
            {t}
          </Tag>
        ))}
      </div>
    </div>
  );
}