export const AUTH_TOKEN = 'kubricks-token';

export function wipeToken() {
  localStorage.removeItem(AUTH_TOKEN);
}

export function fetchToken() {
  return localStorage.getItem(AUTH_TOKEN);
}

export function setToken(token: string) {
  localStorage.setItem(AUTH_TOKEN, token);
}
