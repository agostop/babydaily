<template>
  <div>
    <div id="list">
      <van-list v-model="loading" :finished="finished">
        <van-cell>
          <van-row>
            <van-col span="6">日期</van-col>
            <van-col span="6">时间</van-col>
            <van-col span="6">类型</van-col>
            <van-col span="6">数量</van-col>
          </van-row>
        </van-cell>
        <van-cell v-for="item in recordList" :key="item.time">
          <van-swipe-cell style="line-height: 40px">
            <van-row>
              <van-col span="6">{{item.dateStr}}</van-col>
              <van-col span="6">{{item.timeStr}}</van-col>
              <van-col span="6">{{item.type === 1 ? '奶' : '屎'}}</van-col>
              <van-col span="6">{{item.amount}}</van-col>
            </van-row>
            <template #right>
              <van-button square text="删除" type="danger" class="delete-button" @click="delBabyRecord(item.time)"/>
            </template>
          </van-swipe-cell>
        </van-cell>
      </van-list>
    </div>
    <div id="hm">
      <van-cell-group>
        <div style="margin-top:20px">
          <van-cell center>
            <van-button :loading="shitRequesting" type="info" size="large" @click="addBabyRecord(2)">拉屎了</van-button>
          </van-cell>
        </div>

        <van-cell center>
          <van-button :loading="suckleRequesting" type="primary" size="large" @click="addBabyRecord(1)">喂奶了</van-button>
          <div id="food-val">
            <van-slider active-color="#07c160" max=300 step=10 v-model.number="sliderValue" @change="sliderOnChange">
              <template #button>
                <div class="custom-button">{{ sliderValue }} ml</div>
              </template>
            </van-slider>
          </div>
        </van-cell>

      </van-cell-group>

    </div>
  </div>
</template>

<script>
import { Button, Slider, CellGroup, Cell, List, Row, Col, SwipeCell, Toast } from 'vant'
export default {
  components: {
    [Button.name]: Button,
    [Slider.name]: Slider,
    [CellGroup.name]: CellGroup,
    [Cell.name]: Cell,
    [List.name]: List,
    [Row.name]: Row,
    [Col.name]: Col,
    [SwipeCell.name]: SwipeCell,
    [Toast.name]: Toast
  },
  data() {
    return {
      sliderValue: 50,
      suckleRequesting: false,
      shitRequesting: false,
      showfoodField: 0,
      recordList: [],
      finished: false,
      loading: false
    }
  },
  methods: {
    getBabyRecord: function() {
      this.recordList = []
      this.loading = true
      this.$http.babyRecord.getRecord().then(res => {
        if (res.Response.GroupList != null) {
          res.Response.GroupList.forEach(e => {
            var obj = {}
            obj.timeStr = this.getTimeStr(e.time)
            obj.dateStr = this.getDateStr(e.time)
            obj.type = e.type
            obj.amount = e.amount
            obj.time = e.time
            this.recordList.push(obj)
          })
        }
      })
      this.finished = true
      this.loading = false
    },
    addBabyRecord: function(type) {
      if (type === 1) {
        this.suckleRequesting = true
      }
      if (type === 2) {
        this.shitRequesting = true
      }

      this.$http.babyRecord.addRecord(type, type === 1 ? this.sliderValue : 1).then(res => {
        console.log(res)
        if (res.Response.Data.Status === 'OK') {
          Toast('请求成功')
        } else {
          Toast('请求异常: ' + res.Response.Error)
        }
        this.getBabyRecord()
      })

      this.shitRequesting = false
      this.suckleRequesting = false
    },
    delBabyRecord: function(ts) {
      this.$http.babyRecord.delRecord(ts).then(res => {
        console.log(res)
        if (res.Response.Data.Status === 'OK') {
          Toast('请求成功')
        } else {
          Toast('请求异常: ' + res.Response.Error)
        }
        this.getBabyRecord()
      })
    },
    sliderOnChange: function(num) {
      this.sliderValue = num
    },
    suckleOnChange: function() {
      // this.showfoodField = this.showfoodField ^ 1
    },
    getTimeStr: function(time) {
      var date = new Date(time * 1000)
      var hours = date.getHours()
      var minutes = date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()
      var seconds = date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds()
      // Will display time in 10:30:23 format
      return hours + ':' + minutes + ':' + seconds
    },
    getDateStr: function(time) {
      var da = new Date(time * 1000)
      var year = da.getFullYear()
      var month = da.getMonth() + 1
      var date = da.getDate()
      month = month < 10 ? '0' + month : month
      date = date < 10 ? '0' + date : date
      return year + '.' + month + '.' + date
    }

  },
  mounted() {
    this.getBabyRecord()
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
#list {
  margin: 0 auto;
}
#hm {
  top: 500px;
  padding-top: 10px;
  position: absolute;
  width: 100%;
}
#food-val {
  padding-top: 20px;
  padding-bottom: 10px;
}
.custom-button {
  width: 40px;
  color: black;
  font-size: 10px;
  line-height: 18px;
  text-align: center;
  background-color:blanchedalmond;
  border-radius: 100px;
  }
.van-cell__value--alone{
  color: #323233;
  text-align: center;
}
</style>
