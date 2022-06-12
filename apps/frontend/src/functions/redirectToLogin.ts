export function redirectToLogin() {
  window.location.replace(`/api/auth?state=${window.btoa(window.location.origin)}`);
}