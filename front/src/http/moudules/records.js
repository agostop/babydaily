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

export const delRecord = (ts) => {
  return axios({
    url: 'baby/del',
    method: 'post',
    data: {
      ts: ts
    }
  })
}

export const getRecord = (userId) => {
  return axios({
    url: 'baby/list',
    method: 'get'
  })
}
