function getusertoken() {
  if (sessionStorage.getItem('token')) {
    if (sessionStorage.getItem('exp') <= Date.now() / 1000) {
      sessionStorage.removeItem('exp')
      sessionStorage.removeItem('token')
    } else {
      return sessionStorage.getItem('token')
    }
  }
  if (localStorage.getItem('token')) {
    if (parseInt(localStorage.getItem('exp')) <= Date.now() / 1000) {
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
  return new Headers({ 'Content-Type': 'application/json' });
}
// copy from https://github.com/koesie10/webauthn/blob/master/webauthn.js start
// Decode a base64 string into a Uint8Array.
function _getBytesFromBase64(value) {
  return Uint8Array.from(atob(value), c => c.charCodeAt(0));
}
// Encode an ArrayBuffer into a base64 string.
function _getBase64FromBytes(buffer) {
  var binary = '';
  var bytes = new Uint8Array(buffer);
  var len = bytes.byteLength;
  for (var i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return window.btoa(binary);
}
//copy from https://github.com/koesie10/webauthn/blob/master/webauthn.js end
function getBytesFromBase64(str) {
  return _getBytesFromBase64(ReplaceURL2Std(str));
}
function getBase64FromBytes(bytes) {
  return ReplaceStd2URL(_getBase64FromBytes(bytes));
}
function ReplaceStd2URL(str) {
  //replace + to - & / to _ & delate the padding
  return str.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}
function ReplaceURL2Std(str) {
  //replace - to + & _ to / & add the padding
  return str.replace(/-/g, '+').replace(/_/g, '/').padEnd(str.length + (str.length % 3), '=');
}
export { getusertoken, authedheader, notauthedheader, getBase64FromBytes, getBytesFromBase64 }  