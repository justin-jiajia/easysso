<template>
  <v-alert :title="tiptitle" :type="tiptype" :text="tiptext" class="tip"></v-alert>
</template>
<script setup>
import { notauthedheader, authedheader, getusertoken } from '@/utils/index.js'
import { ref, onMounted } from 'vue'
let tiptitle = ref('请稍等...')
let tiptype = ref('info')
let tiptext = ref('正在加载中')
let query = new URLSearchParams(window.location.search)
let client_id = query.get('client_id')
if (client_id == null) {
  tiptext.value = '空的client_id值'
  tiptype.value = 'error'
  tiptitle.value = '出错了'
} else {
  onMounted(async () => {
    let res, resjson
    try {
      res = await fetch('/api/oath2/getcallback/', {
        method: 'POST',
        headers: notauthedheader(),
        body: JSON.stringify({
          client_id
        })
      })
      resjson = await res.json()
      if (!res.ok) {
        throw new Error(resjson.error)
      }
    } catch (e) {
      tiptext.value = e.toString()
      tiptype.value = 'error'
      tiptitle.value = '出错了'
      throw new Error(e)
    }
    if (getusertoken()) {
      try {
        const callback = resjson.callback
        res = await fetch('/api/user/settings/getcode/', {
          method: 'POST',
          headers: authedheader(),
          body: JSON.stringify({
            client_id,
          })
        })
        resjson = await res.json()
        if (!res.ok) {
          tiptext.value = resjson.error
          tiptype.value = 'error'
          tiptitle.value = '出错了'
        }
        const code = resjson.code
        let nwquery = new URLSearchParams
        nwquery.append('code', code)
        if (query.get('state')) nwquery.append('state', query.get('state'))
        window.location = callback + '?' + nwquery.toString()
      } catch (e) {
        tiptext.value = e.toString()
        tiptype.value = 'error'
        tiptitle.value = '出错了'
        throw new Error(e)
      }
    } else {
      window.location = 'login.html' + window.location.search
    }
  })
}
</script>
<style>
.tip {
  margin-left: 10%;
  margin-right: 10%;
  margin-top: 100px;
}
</style>