<template>
  <div id="main">
    <div id="list" ref="list-container">
      <van-list v-model="loading" :finished="finished">
        <van-sticky>
          <van-cell>
            <van-row>
              <van-col span="6">日期</van-col>
              <van-col span="6">时间</van-col>
              <van-col span="6">类型</van-col>
              <van-col span="6">数量</van-col>
            </van-row>
          </van-cell>
        </van-sticky>
        <van-cell v-for="item in recordList" :key="item.time">
          <van-swipe-cell style="line-height: 40px">
            <van-row>
              <van-col span="6">{{ item.dateStr }}</van-col>
              <van-col span="6">{{ item.timeStr }}</van-col>
              <van-col span="6">{{ item.type === 1 ? '奶' : '屎' }}</van-col>
              <van-col span="6">{{ item.amount }}</van-col>
            </van-row>
            <template #right>
              <van-button
                square
                text="删除"
                type="danger"
                class="row-button"
                @click="delBabyRecord(item.id)"
              />
            </template>
            <template #left>
              <van-button
                square
                text="编辑"
                type="info"
                class="row-button"
                @click="showEdit(item)"
              />
            </template>
          </van-swipe-cell>
        </van-cell>
      </van-list>
    </div>
    <div>
      <van-popup v-model="popUpShow" class="popupshow">
        <van-form @submit="editSubmit">
          <van-field v-show="false" name="id" :value="recordId" />
          <van-field name="type" label="类型" :value="recordType">
            <template #input>
              <van-radio-group v-model="recordType" direction="horizontal">
                <van-radio name="1">喂奶</van-radio>
                <van-radio name="2">拉屎</van-radio>
              </van-radio-group>
            </template>
          </van-field>
          <van-field
            v-show="recordType === '1'"
            name="amount"
            label="奶量"
            :value="recordAmount"
          >
            <template #input>
              <van-slider
                v-model="recordAmount"
                @change="sliderOnChange"
                max="300"
                step="10"
              >
                <template #button>
                  <div class="custom-button">{{ recordAmount }} ml</div>
                </template>
              </van-slider>
            </template>
          </van-field>
          <van-field
            readonly
            clickable
            name="time"
            :value="recordTime"
            label="时间选择"
            placeholder="点击选择时间"
            @click="timeShow = true"
          />
          <van-action-sheet v-model="timeShow">
            <van-datetime-picker
              v-model="currentDate"
              type="datetime"
              title="选择完整时间"
              :min-date="minDate"
              :max-date="maxDate"
              @confirm="editDate"
              @cancel="cancelDatePicker"
            />
          </van-action-sheet>
          <div style="margin: 16px;">
            <van-button round block type="info" native-type="submit"
              >提交</van-button
            >
          </div>
        </van-form>
      </van-popup>
    </div>
    <div id="hm">
      <van-cell-group>
        <van-cell center>
          <div id="food-val">
            <van-slider
              active-color="#07c160"
              max="300"
              step="10"
              v-model.number="recordAmount"
              @change="sliderOnChange"
              bar-height="20px"
            >
              <template #button>
                <div class="custom-button">{{ recordAmount }} ml</div>
              </template>
            </van-slider>
          </div>
          <van-button
            :loading="suckleRequesting"
            type="primary"
            size="large"
            @click="addBabyRecord(1)"
            >喂奶了</van-button
          >
        </van-cell>

        <div style="margin-top:20px">
          <van-cell center>
            <van-button
              :loading="shitRequesting"
              type="info"
              size="large"
              @click="addBabyRecord(2)"
              >拉屎了</van-button
            >
          </van-cell>
        </div>
      </van-cell-group>
    </div>
  </div>
</template>

<script>
import {
  Button,
  Slider,
  CellGroup,
  Cell,
  List,
  Row,
  Col,
  SwipeCell,
  Toast,
  Sticky,
  DatetimePicker,
  ActionSheet,
  Form,
  Field,
  Popup,
  Radio,
  RadioGroup
} from 'vant'
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
    [Toast.name]: Toast,
    [Sticky.name]: Sticky,
    [DatetimePicker.name]: DatetimePicker,
    [ActionSheet.name]: ActionSheet,
    [Form.name]: Form,
    [Field.name]: Field,
    [Popup.name]: Popup,
    [Radio.name]: Radio,
    [RadioGroup.name]: RadioGroup
  },
  data() {
    return {
      recordId: 0,
      recordAmount: 50,
      recordTime: null,
      recordType: '1',
      suckleRequesting: false,
      shitRequesting: false,
      showfoodField: 0,
      recordList: [],
      finished: false,
      loading: false,
      timeShow: false,
      minDate: new Date(2022, 0, 1),
      maxDate: new Date(2035, 10, 1),
      currentDate: new Date(),
      popUpShow: false,
      sliderShowOnEdit: true
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
            obj.id = e.id
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

      this.$http.babyRecord
        .addRecord(type, type === 1 ? this.recordAmount : 1)
        .then(res => {
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
    delBabyRecord: function(id) {
      this.$http.babyRecord.delRecord(id).then(res => {
        console.log(res)
        if (res.Response.Data.Status === 'OK') {
          Toast('请求成功')
        } else {
          Toast('请求异常: ' + res.Response.Error)
        }
        this.popUpShow = false
        this.getBabyRecord()
      })
    },
    editSubmit: function(v) {
      var type = Number(v.type)
      this.$http.babyRecord
        .editRecord(
          v.id,
          type,
          type === 2 ? 1 : v.amount,
          Date.parse(v.time) / 1000
        )
        .then(res => {
          console.log(res)
          if (res.Response.Data.Status === 'OK') {
            Toast('请求成功')
          } else {
            Toast('请求异常: ' + res.Response.Error)
          }
          this.popUpShow = false
          this.getBabyRecord()
        })
    },
    sliderOnChange: function(num) {
      this.recordAmount = num
      Toast(num)
    },
    getTimeStr: function(time) {
      var date = new Date(time * 1000)
      return this.formatTime(date)
    },
    getDateStr: function(time) {
      var da = new Date(time * 1000)
      return this.formatDate(da)
    },
    showEdit: function(item) {
      this.recordId = item.id
      this.recordAmount = item.amount
      this.recordTime = item.dateStr + ' ' + item.timeStr
      this.recordType = '' + item.type
      this.currentDate = new Date(this.recordTime)
      this.popUpShow = true
    },
    editDate: function(date) {
      var dateStr = this.formatDate(date)
      var timeStr = this.formatTime(date)
      this.recordTime = dateStr + ' ' + timeStr
      // var ts = Math.floor(date.getTime() / 1000)
      this.timeShow = false
    },
    cancelDatePicker: function() {
      this.timeShow = false
    },
    formatDate: function(da) {
      var year = da.getFullYear()
      var month = da.getMonth() + 1
      var date = da.getDate()
      month = month < 10 ? '0' + month : month
      date = date < 10 ? '0' + date : date
      return year + '-' + month + '-' + date
    },
    formatTime: function(date) {
      var hours = date.getHours()
      var minutes =
        date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()
      var seconds =
        date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds()
      // Will display time in 10:30:23 format
      return hours + ':' + minutes + ':' + seconds
    }
  },
  mounted() {
    this.getBabyRecord()
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
#main {
  padding: 10px;
}
#list {
  max-height: 500px;
  overflow: auto;
}
#hm {
  /* top: 500px; */
  padding-top: 10px;
  /* position: absolute; */
  width: 100%;
}
#food-val {
  padding-top: 20px;
  padding-bottom: 10px;
}
.custom-button {
  /* width: 40px; */
  color: black;
  font-size: 10px;
  line-height: 18px;
  text-align: center;
  background-color: blanchedalmond;
  border-radius: 100px;
}
.van-cell__value--alone {
  color: #323233;
  text-align: center;
}
.row-button {
  height: 100%;
}
.popupshow {
  width: 100%;
  height: 50%;
}
</style>
