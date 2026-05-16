import { message } from "antd";
import path from "path-browserify";
import ApiRequest from "./Api";

export function ReadDir(dir) {
  return ApiRequest(path.join("/simple_upload/objects", dir), {
    method: "GET",
  }).catch((e) => {
    message.error(e);
  });
}

export function UploadFile(uploadFilepath, formdata) {
  return ApiRequest(path.join("/simple_upload/object", uploadFilepath), {
    method: "POST",
    body: formdata,
  });
}

export function DeleteFile(filepath) {
  return ApiRequest(path.join("/simple_upload/object", filepath), {
    method: "DELETE",
  });
}

export function SearchFiles(keyword, tag) {
  const params = new URLSearchParams();
  if (keyword) params.set("Keyword", keyword);
  if (tag) params.set("Tag", tag);
  return ApiRequest("/simple_upload/search?" + params.toString(), {
    method: "GET",
  }).catch((e) => {
    message.error(e);
  });
}

export function UpdateFileTags(filepath, tags) {
  return ApiRequest(path.join("/simple_upload/object", filepath), {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ Tags: tags }),
  });
}

export function Mkdir(dirpath) {
  return ApiRequest(path.join("/simple_upload/mkdir", dirpath), {
    method: "POST",
  });
}

export function ScanFiles() {
  return ApiRequest("/api/scanfs", {
    method: "POST",
  });
}