function ApiRequest(url, initRequest) {
  return fetch(url, initRequest)
    .then((resp) => resp.json())
    .then((data) => {
      if (data.Code == "success") {
        return data.Data;
      }
      throw new Error("code:" + data?.Code + " message:" + data?.Message);
    });
}

export default ApiRequest;
