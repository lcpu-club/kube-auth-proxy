import { ref } from "vue";

const LOCAL_STORAGE_TOKEN_KEY = "kube-auth-proxy-ui-token";

function loadToken() {
  const token = localStorage.getItem(LOCAL_STORAGE_TOKEN_KEY);
  if (token) {
    setToken(token);
  }
}

export const token = ref("");
export const setToken = (newToken: string) => {
  token.value = newToken;
  localStorage.setItem(LOCAL_STORAGE_TOKEN_KEY, newToken);
};
export const clearToken = () => {
  token.value = "";
  localStorage.removeItem(LOCAL_STORAGE_TOKEN_KEY);
};
export const hasToken = () => {
  return token.value !== "";
};
export const useToken = () => {
  if (!hasToken()) {
    loadToken();
  }

  return token;
};
export const getToken = () => {
  return useToken().value;
};
