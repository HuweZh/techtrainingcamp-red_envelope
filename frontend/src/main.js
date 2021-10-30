import Vue from 'vue'
import App from './App.vue'

require('./mock.js')

Vue.config.productionTip = false

new Vue({
  render: h => h(App),
}).$mount('#app')
