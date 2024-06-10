<template>
    <v-app>
        <v-app-bar elevation="5">
            <v-app-bar-title>个人设置</v-app-bar-title>
            <template v-slot:append>
                <v-btn icon="mdi-logout" @click="logout()"></v-btn>
            </template>
        </v-app-bar>
        <v-main>
            <v-container>
                <v-row justify="center">
                    <v-card>
                        <v-card-item>
                            <v-card-title>认证器设置</v-card-title>
                        </v-card-item>
                        <v-card-text>
                            <v-container>
                                <v-row>
                                    <v-col>
                                        <v-card>
                                            <v-card-item>
                                                <v-card-title>添加认证器</v-card-title>
                                            </v-card-item>
                                            <v-card-text>
                                                配置后你将能够用支持的设备通过指纹、FaceID等方式进行快速登录
                                            </v-card-text>
                                            <v-card-actions>
                                                <v-btn prepend-icon="mdi-add" @click="add_authenticator">添加认证器</v-btn>
                                            </v-card-actions>
                                        </v-card>
                                    </v-col>
                                    <v-col>
                                        <v-card>
                                            <v-card-item>
                                                <v-card-title>管理认证器<v-btn icon="mdi-refresh" @click="refresh_aus()"
                                                        :loading="refreshloading"></v-btn></v-card-title>
                                            </v-card-item>
                                            <v-card-actions>
                                                <v-table>
                                                    <thead>
                                                        <tr>
                                                            <th>名称</th>
                                                            <th>创建日期</th>
                                                            <th>上次使用</th>
                                                            <th>无用户名登录</th>
                                                            <th>操作</th>
                                                        </tr>
                                                    </thead>
                                                    <tbody>
                                                        <tr v-for="item in aus">
                                                            <th>{{ item.name }}</th>
                                                            <th>{{ getdatestring(item.created) }}</th>
                                                            <th>{{ getdatestring(item.last_used) }}</th>
                                                            <th><v-icon
                                                                    :icon="item.username_less ? 'mdi-check' : 'mdi-close'"></v-icon>
                                                            </th>
                                                            <th><v-btn icon="mdi-pencil"
                                                                    @click="editau(item.id)"></v-btn>
                                                                <v-btn icon="mdi-delete"
                                                                    @click="deleteau(item.id)"></v-btn>
                                                            </th>
                                                        </tr>
                                                    </tbody>
                                                </v-table>
                                            </v-card-actions>
                                        </v-card>
                                    </v-col>
                                </v-row>
                            </v-container>
                        </v-card-text>
                    </v-card>
                </v-row>
                <v-row justify="center">
                    <v-col>
                        <v-card>
                            <v-card-item>
                                <v-card-title>安全设置</v-card-title>
                            </v-card-item>
                            <v-card-actions>
                                <v-btn prepend-icon="mdi-form-textbox-password" @click="changepasswd">修改密码</v-btn>
                                <v-btn prepend-icon="mdi-delete" @click="removeaccount">注销账户</v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-col>
                    <!-- <v-spacer></v-spacer> -->
                    <v-col>
                        <v-card>
                            <v-card-item>
                                <v-card-title>行为日志<v-btn icon="mdi-refresh" @click="refresh_log()"
                                        :loading="logloading"></v-btn></v-card-title>
                            </v-card-item>
                            <v-card-actions>
                                <v-table>
                                    <thead>
                                        <tr>
                                            <th>操作名称</th>
                                            <th>日期</th>
                                            <th>设备</th>
                                            <th>IP</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr v-for="item in log">
                                            <th>{{ item.action }}</th>
                                            <th>{{ getdatestring(item.actiontime) }}</th>
                                            <th>{{ item.ua }}</th>
                                            <th>{{ item.ip }}</th>
                                        </tr>
                                    </tbody>
                                </v-table>
                            </v-card-actions>
                        </v-card>
                    </v-col>
                </v-row>
            </v-container>
        </v-main>
    </v-app>
</template>
<script setup>
import swal from 'sweetalert';
import { authedheader, getusertoken, getBase64FromBytes, getBytesFromBase64 } from '@/utils/index.js'
import { ref } from 'vue'
let aus = ref([])
let log = ref([])
let refreshloading = ref(false)
let logloading = ref(false)
if (getusertoken() == null) {
    window.location = 'login.html'
}
function add_authenticator() {
    swal({
        title: "是否绑定为密钥驻留模式？",
        text: "你可以无用户名登录",
        buttons: true,
    }).then((v) =>
        fetch('/api/user/settings/webauthn/startregistration' + (v ? '?discovered=1' : ''), {
            method: 'GET',
            headers: authedheader(),
        }))
        .then(res => {
            if (!res.ok) {
                res.json().then(resjson => {
                    swal('绑定失败', resjson.error, 'error')
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
            resjson.publicKey.user.id = getBytesFromBase64(resjson.publicKey.user.id);
            if (resjson.publicKey.excludeCredentials) {
                for (var i = 0; i < resjson.publicKey.excludeCredentials.length; i++) {
                    resjson.publicKey.excludeCredentials[i].id = getBytesFromBase64(resjson.publicKey.excludeCredentials[i].id);
                }
            }
            return resjson
        })
        .then(c => navigator.credentials.create(c))
        .then(res => { console.log(res); return res })
        .then(res => fetch('/api/user/settings/webauthn/finishregistration/', {
            method: 'POST',
            headers: authedheader(),
            body: JSON.stringify({
                rawId: getBase64FromBytes(res.rawId),
                type: res.type,
                id: res.id,
                response: {
                    attestationObject: getBase64FromBytes(res.response.attestationObject),
                    clientDataJSON: getBase64FromBytes(res.response.clientDataJSON),
                }
            })
        }))
        .then(res => {
            if (!res.ok) {
                res.json().then(resjson => {
                    swal('绑定失败', resjson.error, 'error')
                })
                throw new Error("!!!");
            } else {
                swal("成功", "认证器设置成功，要不重命名为更好认的名字？", "success")
                refresh_aus()
            }
        })
        .catch(nerr => {
            if (nerr.message == "!!!") return
            swal('注册失败', nerr.message, 'error')
            throw nerr
        })
}
function refresh_aus() {
    refreshloading.value = true
    fetch('/api/user/settings/webauthn/list/', {
        method: 'GET',
        headers: authedheader(),
    })
        .then(res => {
            if (!res.ok) {
                res.json().then(resjson => {
                    swal('失败', resjson.error, 'error')
                })
                throw new Error("!!!");
            } else {
                return res
            }
        })
        .then(res => res.json())
        .then(res => aus.value = res)
        .catch(nerr => {
            if (nerr.message == "!!!") return
            swal('查询失败', nerr.message, 'error')
            throw nerr
        }).finally(() => refreshloading.value = false)
}
function editau(id) {
    swal({
        title: '重命名',
        content: {
            element: "input",
            attributes: {
                placeholder: "在此输入新名字...",
            },
        },
    }).then((res) =>
        fetch('/api/user/settings/webauthn/edit/', {
            method: 'POST',
            headers: authedheader(),
            body: JSON.stringify({
                id,
                'new_name': res,
            })
        })).then(res => {
            if (!res.ok) {
                swal('失败', err, 'error')
            } else {
                swal("成功", "重命名成功", "success")
                refresh_aus()
            }
        })
}
function deleteau(id) {
    swal({
        title: '真的要删除吗？',
        buttons: ["别", "删除吧"],
        dangerMode: true,
    }).then((v) => {
        if (!v) { throw Error("user cancelled") }
    })
        .then(() => fetch('/api/user/settings/webauthn/delete/', {
            method: 'POST',
            headers: authedheader(),
            body: JSON.stringify({
                id,
            })
        })).then(res => {
            if (!res.ok) {
                swal('失败', err, 'error')
            } else {
                swal("成功", "删除成功", "success")
                refresh_aus()
            }
        })
}
function getdatestring(val) {
    return new Date(val * 1000).toLocaleString()
}
function changepasswd() {
    var newpasswd;
    swal({
        title: '修改密码',
        content: {
            element: "input",
            attributes: {
                type: "password",
                placeholder: "请输入新密码...",
            }
        },
        buttons: true
    }).then((res) => {
        if (!res) {
            throw Error("user cancelled")
        }
        newpasswd = res
    })
        .then(() => swal({
            title: '修改密码',
            content: {
                element: "input",
                attributes: {
                    placeholder: "请再次输入新密码...",
                    type: "password",
                }
            }
        })).then((res) => {
            console.log(newpasswd)
            console.log(res)
            if (newpasswd != res) {
                throw Error("两次密码不一致")
            }
        }).then(() => fetch('/api/user/settings/change_passwd/', {
            method: 'POST',
            headers: authedheader(),
            body: JSON.stringify({
                "new_passwd": newpasswd
            })
        })).then(res => {
            if (!res.ok) {
                throw Error(res)
            }
            else {
                swal("成功", "修改密码成功", "success")
            }
        }).catch((err) => {
            if (err.message == "user cancelled") return
            swal('修改密码失败', err.message, 'error')
        })
}
function removeaccount() {
    swal({
        title: '注销账户',
        text: '为确认，请输入"unregister"(不带标点符号):',
        content: {
            element: "input",
            attributes: {
                placeholder: "请输入...",
            }
        },
        buttons: true
    }).then((res) => {
        if ('unregister' != res) {
            throw Error("user cancelled")
        }
    }).then(() => fetch('/api/user/settings/remove/', {
        method: 'POST',
        headers: authedheader(),
        body: JSON.stringify({})
    })).then(res => {
        console.log(res)
        if (!res.ok) {
            res.json().then(resjson => {
                swal('注销失败', resjson.error, 'error')
            })
            throw Error("user cancelled")
        }
        else {
            swal("成功", "注销成功", "success").then(() => window.location = 'logout.html')
        }
    }).catch((err) => {
        if (err.message == "user cancelled") return
        swal('注销失败', err.message, 'error')
    })
}
function logout() {
    window.location = '/ui/logout.html';
}
function refresh_log() {
    logloading.value = true
    fetch('/api/user/settings/webauthn/log/', {
        method: 'GET',
        headers: authedheader(),
    })
        .then(res => {
            if (!res.ok) {
                res.json().then(resjson => {
                    swal('失败', resjson.error, 'error')
                })
                throw new Error("!!!");
            } else {
                return res
            }
        })
        .then(res => res.json())
        .then(res => log.value = res)
        .catch(nerr => {
            if (nerr.message == "!!!") return
            swal('查询失败', nerr.message, 'error')
            throw nerr
        }).finally(() => logloading.value = false)
}
refresh_aus()
refresh_log()
</script>