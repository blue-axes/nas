function ApiRequest(url, initRequest) {
  return fetch(url, {
    ...initRequest,
    credentials: "same-origin",
  })
    .then((resp) => {
      if (resp.status === 401 && !url.includes("/api/login")) {
        window.dispatchEvent(new CustomEvent("auth:unauthorized"));
        throw new Error("请先登录");
      }
      return resp.json();
    })
    .then((data) => {
      if (data.Code == "success") {
        return data.Data;
      }
      throw new Error("code:" + data?.Code + " message:" + data?.Message);
    });
}

export default ApiRequest;
