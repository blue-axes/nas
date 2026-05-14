import ApiRequest from "./Api.jsx";

export function GetCurrentUser() {
  return ApiRequest("/api/users/me");
}

export function ListUsers() {
  return ApiRequest("/api/users");
}

export function CreateUser(username, password, canRead, canWrite) {
  return ApiRequest("/api/users", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ Username: username, Password: password, CanRead: canRead, CanWrite: canWrite }),
  });
}

export function UpdateUser(username, canRead, canWrite) {
  return ApiRequest(`/api/users/${encodeURIComponent(username)}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ CanRead: canRead, CanWrite: canWrite }),
  });
}

export function DeleteUser(username) {
  return ApiRequest(`/api/users/${encodeURIComponent(username)}`, {
    method: "DELETE",
  });
}

export function ChangePassword(username, currentPassword, newPassword) {
  return ApiRequest(`/api/users/${encodeURIComponent(username)}/password`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ CurrentPassword: currentPassword, NewPassword: newPassword }),
  });
}