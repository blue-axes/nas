import { useState, useCallback } from "react";
import { Input, Tag } from "antd";
import { SearchOutlined, CloseOutlined } from "@ant-design/icons";

export default function SearchBar({ onSearch }) {
  const [keyword, setKeyword] = useState("");
  const [tag, setTag] = useState("");
  const [tagInput, setTagInput] = useState("");
  const [tags, setTags] = useState([]);

  const doSearch = useCallback(
    (kw, tg) => {
      onSearch(kw, tg);
    },
    [onSearch],
  );

  const handleSearch = (value) => {
    setKeyword(value);
    doSearch(value, tag);
  };

  const addTag = (t) => {
    if (!t || tags.includes(t)) return;
    const newTags = [...tags, t];
    setTags(newTags);
    const tagStr = newTags.join(",");
    setTag(tagStr);
    setTagInput("");
    doSearch(keyword, tagStr);
  };

  const removeTag = (t) => {
    const newTags = tags.filter((x) => x !== t);
    setTags(newTags);
    const tagStr = newTags.join(",");
    setTag(tagStr);
    doSearch(keyword, tagStr);
  };

  const clearAll = () => {
    setKeyword("");
    setTag("");
    setTags([]);
    setTagInput("");
    doSearch("", "");
  };

  const hasFilter = keyword || tags.length > 0;

  const searchStyle = {
    display: "flex",
    alignItems: "center",
    gap: 8,
    flex: 1,
    maxWidth: 420,
  };

  const inputStyle = {
    background: "rgba(15, 23, 42, 0.8)",
    borderColor: "rgba(0, 212, 255, 0.2)",
    color: "#e2e8f0",
    borderRadius: 20,
    fontSize: 13,
    height: 34,
  };

  const tagInputStyle = {
    ...inputStyle,
    width: 80,
  };

  return (
    <div style={searchStyle}>
      <Input
        placeholder="Search files..."
        prefix={<SearchOutlined style={{ color: "#64748b" }} />}
        value={keyword}
        onChange={(e) => handleSearch(e.target.value)}
        onPressEnter={() => doSearch(keyword, tag)}
        allowClear={{ clearIcon: <CloseOutlined style={{ color: "#64748b", fontSize: 12 }} /> }}
        style={inputStyle}
        size="small"
      />
      <Input
        placeholder="+ tag"
        value={tagInput}
        onChange={(e) => setTagInput(e.target.value)}
        onPressEnter={() => addTag(tagInput.trim())}
        style={tagInputStyle}
        size="small"
      />
      {tags.map((t) => (
        <Tag
          key={t}
          closable
          onClose={() => removeTag(t)}
          color="cyan"
          style={{
            margin: 0,
            fontSize: 11,
            borderRadius: 10,
            background: "rgba(0, 212, 255, 0.12)",
            border: "1px solid rgba(0, 212, 255, 0.25)",
            color: "#00d4ff",
          }}
        >
          {t}
        </Tag>
      ))}
      {hasFilter && (
        <span
          onClick={clearAll}
          style={{
            cursor: "pointer",
            color: "#f43f5e",
            fontSize: 12,
            whiteSpace: "nowrap",
          }}
        >
          Clear
        </span>
      )}
    </div>
  );
}