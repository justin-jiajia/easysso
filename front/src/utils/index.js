function getusertoken() {
if (sessionStorage.getItem('token')) {
  if (sessionStorage.getItem('exp') <= Date.now()/1000) {
    sessionStorage.removeItem('exp')
    sessionStorage.removeItem('token')
  } else {
    return sessionStorage.getItem('token')
  }
} else if (localStorage.getItem('token')) {
  if (parseInt(localStorage.getItem('exp')) <= Date.now()/1000) {
    localStorage.removeItem('exp')
    localStorage.removeItem('token')
  } else {
    return localStorage.getItem('token')
  }
}
  return null
}
function authedheader() {
  return new Headers({ 'Content-Type': 'application/json', 'Authorization': getusertoken() });
}
function notauthedheader() {
  return new Headers({ 'Content-Type': 'application/json'});
}
export {getusertoken,authedheader,notauthedheader}