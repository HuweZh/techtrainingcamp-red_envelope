<template>
  <div id="app">
    请输入url前缀（如：http://127.0.0.1/api）：
    <input type="text" v-model="prefix">
    <br>
    用户id: <input type="number" v-model.number="req.uid">
    <br>
    红包id（仅拆红包输入）: <input type="number" v-model.number="req.envelope_id">
    <br>
    <button @click="snatch">抢红包</button>
    <button @click="open">拆红包</button>
    <button @click="get_wallet_list">获取钱包列表</button>
    <br>
    {{ msg }}
    <table v-if="show_table">
      <tr>
        <th>红包id</th>
        <th>已拆开</th>
        <th>面值</th>
        <th>获取时间</th>
      </tr>
      <tr v-for="(item, index) in envelope_list" :key="index">
        <td>{{ item.envelope_id }}</td>
        <td>{{ item.opened }}</td>
        <td>{{ show_value(item.value) }}</td>
        <td>{{ getLocalTime(item.snatch_time) }}</td>
      </tr>
    </table>
  </div>
</template>

<script>

import axios from "axios";

export default {
  data() {
    return {
      prefix: '',
      resp: '',
      msg: '',
      req: {
        uid: 0,
        envelope_id: 0
      },
      show_table: false,
      envelope_list: []
    }
  },
  methods: {

    show_value(value) {
      if (value === undefined) {
        return ""
      }
      return "" + parseInt(value / 100) + "." + (value % 100)
    },

    getLocalTime(nS) {
      return new Date(parseInt(nS) * 1000).toLocaleString().replace(/:\d{1,2}$/, ' ');
    },

    snatch() {
      this.show_table = false;
      axios.post(this.prefix + '/snatch', this.req)
          .then(resp => {
            if (resp.status !== 200) {
              this.msg = "请求失败，status=" + resp.status
            } else {
              if (resp.data.code !== 0) {
                this.msg = "抢红包失败：" + resp.data.code + ":" + resp.data.msg
              } else {
                this.msg = "抢红包成功，红包id:" + resp.data.data.envelope_id +
                    "\n\n当前第" + resp.data.data.cur_count + "次，最多抢" + resp.data.data.max_count + "次\n"
              }
            }
          });
    },

    open() {
      this.show_table = false;
      axios.post(this.prefix + '/open', this.req)
          .then(resp => {
            if (resp.status !== 200) {
              this.msg = "请求失败，status=" + resp.status
            } else {
              if (resp.data.code !== 0) {
                this.msg = "拆红包失败：" + resp.data.code + ":" + resp.data.msg
              } else {
                this.msg = "拆红包成功，金额：" + this.show_value(resp.data.data.value)
              }
            }
          });
    },

    get_wallet_list() {
      this.show_table = false;
      axios.post(this.prefix + '/get_wallet_list', this.req)
          .then(resp => {
            if (resp.status !== 200) {
              this.msg = "请求失败，status=" + resp.status
            } else {
              if (resp.data.code !== 0) {
                this.msg = "获取钱包失败：" + resp.data.code + ":" + resp.data.msg
              } else {
                this.msg = "钱包总额：" + this.show_value(resp.data.data.amount)
                this.show_table = true
                this.envelope_list = resp.data.data.envelope_list
              }
            }
          });
    }
  }
}
</script>

<style>

</style>
