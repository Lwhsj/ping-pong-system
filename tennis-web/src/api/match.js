import request from '@/utils/request'

export function getMatchList(params) {
  return request({
    url: '/matches',
    method: 'get',
    params
  })
}

export function getMatchDetail(id) {
  return request({
    url: `/match/${id}/detail`,
    method: 'get'
  })
}

export function getMatchStats(id) {
  return request({
    url: `/match/${id}/stats`,
    method: 'get'
  })
}

export function getCurrentScore(id) {
  return request({
    url: `/match/${id}/current`,
    method: 'get'
  })
}

export function exportMatch(id) {
  return request({
    url: `/match/${id}/export`,
    method: 'get',
    responseType: 'blob' // Important for downloading files
  })
}
