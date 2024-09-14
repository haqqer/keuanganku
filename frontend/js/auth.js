import { API_URL } from "./constants.js";

const checkToken = async () => {
  const accessToken = localStorage.getItem("access_token");
  const refreshToken = localStorage.getItem("refresh_token");
  if (accessToken === null && refreshToken === null) {
    location.pathname = "/login";
  }
  const res = await authMe();
  if (res.status != 200) {
    const resRefreshToken = await authRefreshToken();
    const newToken = await resRefreshToken.json();
    storeToken(newToken.data.token);
  }
};

export const authRefreshToken = () => {
  return fetch(`${API_URL}/refresh-token`, {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
      Authorization: `Bearer ${localStorage.getItem("access_token")}`,
    }),
    body: JSON.stringify({
      refresh_token: localStorage.getItem("refresh_token"),
    }),
  });
};

export const authMe = () => {
  return fetch(`${API_URL}/auth/me`, {
    headers: new Headers({
      "Content-Type": "application/json",
      Authorization: `Bearer ${localStorage.getItem("access_token")}`,
    }),
  });
};

export const logout = () => {
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  location.pathname = "/login";
};

export const storeToken = (token) => {
  localStorage.setItem("access_token", token.access_token);
  localStorage.setItem("refresh_token", token.refresh_token);
};
document.addEventListener("DOMContentLoaded", checkToken);
