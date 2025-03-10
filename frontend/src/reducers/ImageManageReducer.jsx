export function ImageList(list, action) {
  let payload = action.payload;
  let index = -1;
  switch (action.type) {
    case "list":
      list.length = 0;
      for (let item of payload) {
        list.push(item);
      }
      break;
    case "remove":
      index = list.findIndex((item) => item.Name === payload);
      return [...list.slice(0, index), ...list.slice(index + 1, list.length)];
    default:
      return list;
  }
}

export function FileList(list, action) {
  let payload = action.payload;
  let index = -1;
  let item;
  switch (action.type) {
    case "add":
      list.push(payload);
      return;
    case "remove":
      index = list.indexOf(payload);
      return [...list.slice(0, index), ...list.slice(index + 1, list.length)];
    case "clear":
      list.length = 0;
      return;
    case "process":
      index = list.indexOf(payload.file);
      item = list[index];
      item.percent = payload.process;
      item.status = "uploading";
      return [
        ...list.slice(0, index),
        item,
        ...list.slice(index + 1, list.length),
      ];
    case "done":
      index = list.indexOf(payload.file);
      item = list[index];
      item.percent = 100;
      item.status = "done";

      return [
        ...list.slice(0, index),
        item,
        ...list.slice(index + 1, list.length),
      ];
    case "error":
      index = list.indexOf(payload.file);
      item = list[index];
      item.percent = 100;
      item.status = "error";
      return [
        ...list.slice(0, index),
        item,
        ...list.slice(index + 1, list.length),
      ];
  }
}
