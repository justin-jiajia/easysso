<template>
  <v-app>
    <v-app-bar elevation="5">
      <v-app-bar-title>登录/注册以继续<span v-if="name != ''">前往{{ name }}</span></v-app-bar-title>
    </v-app-bar>
    <v-main>
      <v-card class="login-card">
        <div class="container d-flex justify-center">
          <div class="circle rounded-circle bg-primary text-white d-flex justify-center align-center">
            <v-icon icon="mdi-login-variant" class="login-icon"></v-icon>
          </div>
        </div>
        <v-tabs v-model="tab">
          <v-tab value="login">登录</v-tab>
          <v-tab value="signup">注册</v-tab>
        </v-tabs>
        <v-card-text>
          <v-window v-model="tab">
            <v-window-item value="login">
              <v-form v-model="login_valid" @submit.prevent>
                <v-text-field label="用户名" prepend-icon="mdi-account" v-model="login_username"
                  :rules="usernamerule"></v-text-field>
                <v-text-field label="密码" prepend-icon="mdi-lock"
                  :append-icon="login_passwd_visible ? 'mdi-eye' : 'mdi-eye-off'"
                  :type="login_passwd_visible ? 'text' : 'password'"
                  @click:append="login_passwd_visible = !login_passwd_visible" v-model="login_passwd"
                  :rules="passwdrule"></v-text-field>
                <v-checkbox label="记住我" v-model="login_remember_me"></v-checkbox>
                <v-btn type="submit" @click="login">登录</v-btn></v-form>
              <v-divider></v-divider>
              <v-text-field label="用户名（留空以使用无用户名登录[密钥驻留]）" prepend-icon="mdi-account"
                v-model="login_username"></v-text-field>
              <v-btn color=" primary" @click="usewebauthn">或使用认证器登录</v-btn>
            </v-window-item>
            <v-window-item value="signup">
              <v-form v-model="sign_up_valid" @submit.prevent>
                <v-text-field label="用户名" prepend-icon="mdi-account" v-model="sign_up_username"
                  :rules="usernamerule"></v-text-field>
                <v-text-field label="密码" prepend-icon="mdi-lock"
                  :append-icon="sign_up_passwd_visible ? 'mdi-eye' : 'mdi-eye-off'"
                  :type="sign_up_passwd_visible ? 'text' : 'password'"
                  @click:append="sign_up_passwd_visible = !sign_up_passwd_visible" v-model="sign_up_passwd"
                  :rules="passwsignupdrule"></v-text-field>
                <v-text-field label="确认密码" prepend-icon="mdi-lock"
                  :append-icon="sign_up_passwd_again_visible ? 'mdi-eye' : 'mdi-eye-off'"
                  :type="sign_up_passwd_again_visible ? 'text' : 'password'"
                  @click:append="sign_up_passwd_again_visible = !sign_up_passwd_again_visible"
                  v-model="sign_up_passwd_again" :rules="passwsignupdrule"></v-text-field>
                <v-checkbox label="记住我" v-model="sign_up_remember_me"></v-checkbox>
                <v-btn type="submit" @click="signup">注册</v-btn></v-form>
            </v-window-item>
          </v-window>
        </v-card-text>
      </v-card>
    </v-main>
  </v-app>
</template>
<script setup>
import swal from 'sweetalert';
import { getusertoken, getBase64FromBytes, getBytesFromBase64 } from '@/utils/index.js'
import { ref } from 'vue'
//public
let tab = ref(null)
let nooath2 = ref(false)
//public end
// get client id from url query
let query = new URLSearchParams(window.location.search)
let clientid = query.get('client_id')
let name = ref('')
if (!clientid) {
  nooath2.value = true
} else {
  fetch('/api/oath2/getcallback/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      client_id: clientid
    })
  })
    .then((res) => {
      res.json().then((resjson) => {
        if (!res.ok) {
          throw new Error(resjson.error)
        } else {
          name.value = resjson.name
        }
      })
    })
    .catch((err) => {
      console.log(err)
      swal('获取应用详情错误', err.message, 'error')
    })
}
//test if login
if (getusertoken()) {
  if (nooath2.value) {
    window.location.href = 'settings.html'
  } else {
    window.location.href = 'authorize.html' + window.location.search
  }
}
//login
let login_passwd = ref('')
let login_username = ref('')
let login_remember_me = ref(false)
let login_passwd_visible = ref(false)
let login_valid = ref(false)
let usernamerule = [
  (v) => {
    if (v) return true
    return '用户名不能为空'
  }
]
let passwdrule = [
  (v) => {
    if (v) return true
    return '密码不能为空'
  }
]
function login() {
  if (!login_valid.value) return
  fetch('/api/user/sign_in/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      username: login_username.value,
      password: login_passwd.value
    })
  })
    .then((res) => {
      res.json().then((resjson) => {
        if (!res.ok) {
          swal('登录失败', resjson.error, 'error')
          throw new Error(resjson.error)
        } else {
          console.log(resjson)
          let wheretostorage
          if (login_remember_me.value) {
            wheretostorage = localStorage
          } else {
            wheretostorage = sessionStorage
          }
          wheretostorage.setItem('token', resjson.token)
          wheretostorage.setItem('exp', resjson.expire)
          if (nooath2.value) {
            window.location.href = 'settings.html'
          } else {
            window.location.href = 'authorize.html' + window.location.search
          }
        }

      })
    })
    .catch((err) => {
      console.log(err)
      swal('登录失败', err.message, 'error')
    })
}
function usewebauthn() {
  fetch('/api/user/startwlogin/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      username: login_username.value
    })
  })
    .then(res => {
      if (!res.ok) {
        res.json().then(resjson => {
          swal('登录失败', resjson.error, 'error')
        })
        throw new Error("!!!");
      } else {
        return res
      }
    })
    .then(res => res.json())
    .then(res => res.options)
    .then(resjson => {
      console.log(resjson)
      resjson.publicKey.challenge = getBytesFromBase64(resjson.publicKey.challenge);
      if (resjson.publicKey.allowCredentials) {
        for (var i = 0; i < resjson.publicKey.allowCredentials.length; i++) {
          resjson.publicKey.allowCredentials[i].id = getBytesFromBase64(resjson.publicKey.allowCredentials[i].id);
        }
      }
      return resjson
    })
    .then(c => navigator.credentials.get(c))
    .then(res => { console.log(res); return res })
    .then(res => fetch('/api/user/finishwlogin/', {
      method: 'POST',
      body: JSON.stringify({
        rawId: getBase64FromBytes(res.rawId),
        type: res.type,
        id: res.id,
        response: {
          authenticatorData: getBase64FromBytes(res.response.authenticatorData),
          signature: getBase64FromBytes(res.response.signature),
          userHandle: getBase64FromBytes(res.response.userHandle),
          clientDataJSON: getBase64FromBytes(res.response.clientDataJSON),
        }
      })
    }))
    .then(res => {
      if (!res.ok) {
        res.json().then(resjson => {
          swal('登录失败', resjson.error, 'error')
        })
        throw new Error("!!!");
      }
      return res
    })
    .then(res => res.json())
    .then(resjson => {
      console.log(resjson)
      let wheretostorage
      if (login_remember_me.value) {
        wheretostorage = localStorage
      } else {
        wheretostorage = sessionStorage
      }
      wheretostorage.setItem('token', resjson.token)
      wheretostorage.setItem('exp', resjson.expire)
      if (nooath2.value) {
        window.location.href = 'settings.html'
      } else {
        window.location.href = 'authorize.html' + window.location.search
      }
    })
    .catch(nerr => {
      if (nerr.message == "!!!") return
      swal('登录失败', nerr.message, 'error')
      throw nerr
    })
}
//login end
//signup
let sign_up_passwd = ref('')
let sign_up_passwd_again = ref('')
let sign_up_passwd_visible = ref(false)
let sign_up_username = ref('')
let sign_up_remember_me = ref(false)
let sign_up_passwd_again_visible = ref(false)
let sign_up_valid = ref(false)
let passwsignupdrule = [
  (v) => {
    if (v && v == sign_up_passwd.value) return true
    if (!v) return '请再次输入密码'
    return '两次输入密码不一致'
  }
]
function signup() {
  if (!sign_up_valid.value) return
  fetch('/api/user/sign_up/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      username: sign_up_username.value,
      password: sign_up_passwd.value
    })
  })
    .then((res) => {
      res.json().then((resjson) => {
        if (!res.ok) {
          // console.log(resjson.error);
          tipshow.value = true
          tiptitle.value = '注册失败'
          tiptype.value = 'error'
          tiptext.value = resjson.error
          reset_tip()
        } else {
          console.log(resjson)
          let wheretostorage
          if (sign_up_remember_me.value) {
            wheretostorage = localStorage
          } else {
            wheretostorage = sessionStorage
          }
          wheretostorage.setItem('token', resjson.token)
          wheretostorage.setItem('exp', resjson.expire)
          if (nooath2.value) {
            window.location.href = 'settings.html'
          } else {
            window.location.href = 'authorize.html' + window.location.search
          }
        }
      })
    })
    .catch((err) => {
      console.log(err)
      tipshow.value = true
      tiptitle.value = '注册失败'
      tiptype.value = 'error'
      tiptext.value = '网络错误'
    })
}
</script>
<style>
.login-card {
  margin-top: 100px;
  margin-right: 10vw;
  margin-left: 10vw;
}

.container {
  margin-top: 20px;
  width: 100%;
}

.circle {
  width: 75px;
  height: 75px;
}

.login-icon {
  font-size: 45px !important;
}

.tip {
  position: fixed !important;
  right: 10vw !important;
  left: 10vw !important;
  z-index: 999;
}

@keyframes fade {
  to {
    opacity: 0%;
  }
}

.fade-enter-active {
  animation-direction: reverse;
  animation-name: fade;
  animation-duration: 1s;
}

.fade-leave-active {
  animation-name: fade;
  animation-duration: 3s;
}
</style>
