import { authedheader } from '@/utils/index.js'
if (sessionStorage.getItem('token')) {
  sessionStorage.removeItem('exp')
  sessionStorage.removeItem('token')
}
if (localStorage.getItem('token')) {
  localStorage.removeItem('exp')
  localStorage.removeItem('token')
}
fetch('/api/user/settings/remove_token', {
  headers: authedheader(),
  method: 'POST'
}).then((res) => {
  if (!res.ok) {
    res.json().then((resjson) => {
      swal('登出失败', resjson.error, 'error')
    })
  }
}).then(res.json()).then((resjson) => {
  if (resjson.status != "ok") {
    swal('登出失败', resjson.error, 'error')
  } else {
    swal('Bye')
  }
})
window.location = "/ui/login.html"