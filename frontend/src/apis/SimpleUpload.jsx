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
