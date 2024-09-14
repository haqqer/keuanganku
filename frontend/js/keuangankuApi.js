import { API_URL } from "./constants.js";
import { authRefreshToken, storeToken } from "./auth.js";

export const getAll = (month) => {
  return fetch(`${API_URL}/api/tx?month=${month}`, {
    headers: new Headers({
      "Content-Type": "application/json",
      Authorization: `Bearer ${localStorage.getItem("access_token")}`,
    }),
  })
    .then((res) => res.json())
    .then((result) => {
      return result.data;
    })
    .catch((err) => {
      return err;
    });
};

export const postTx = (payload) => {
  return fetch(`${API_URL}/api/tx`, {
    method: "POST",
    body: JSON.stringify(payload),
    headers: new Headers({
      "Content-Type": "application/json",
      Authorization: `Bearer ${localStorage.getItem("access_token")}`,
    }),
  })
    .then((res) => res.json())
    .then((result) => {
      return result.data;
    })
    .catch((err) => {
      return err;
    });
};

export const delTx = (id) => {
  return fetch(`${API_URL}/api/tx/${id}`, {
    method: "DELETE",
    headers: new Headers({
      "Content-Type": "application/json",
      Authorization: `Bearer ${localStorage.getItem("access_token")}`,
    }),
  })
    .then((res) => res.json())
    .then((result) => {
      return result.data;
    })
    .catch((err) => {
      return err;
    });
};

export const getChart = (txType = "out", month) => {
  return fetch(`${API_URL}/api/chart?type=${txType}&month=${month}`, {
    headers: new Headers({
      "Content-Type": "application/json",
      Authorization: `Bearer ${localStorage.getItem("access_token")}`,
    }),
  })
    .then((res) => res.json())
    .then((result) => {
      return result.data;
    })
    .catch((err) => {
      return err;
    });
};
