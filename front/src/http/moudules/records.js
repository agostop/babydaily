import axios from '../axios'

export const addRecord = (type, amount) => {
  return axios({
    url: 'baby/add',
    method: 'post',
    data: {
      type: type,
      amount: amount
    }
  })
}

export const delRecord = (id) => {
  return axios({
    url: 'baby/del',
    method: 'post',
    data: {
      id: id
    }
  })
}

export const getRecord = (userId) => {
  return axios({
    url: 'baby/list',
    method: 'get'
  })
}

export const editRecord = (id, type, amount, dateTime) => {
  return axios({
    url: 'baby/edit',
    method: 'post',
    data: {
      id: id,
      type: type,
      amount: amount,
      time: dateTime
    }
  })
}
