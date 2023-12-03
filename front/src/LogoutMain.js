if (sessionStorage.getItem('token')) {
  sessionStorage.removeItem('exp')
  sessionStorage.removeItem('token')
}
if (localStorage.getItem('token')) {
  localStorage.removeItem('exp')
  localStorage.removeItem('token')
}
window.location = "/ui/login.html"