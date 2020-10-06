import store from '@/store'
import { Base64 } from 'js-base64'

import axios from 'axios'

export function parseToken (token) {

  const parts = token.split('.')

  if (parts.length !== 3) {
    throw new Error('token malformed')
  }

  const data = JSON.parse(Base64.decode(parts[1]))

  if (Math.round(new Date().getTime() / 1000) > data.exp) {
    throw new Error('token expired')
  }

  if(data.user) {
    store.commit('setUser', data.user)
  }
  localStorage.setItem('jwt', token)
  store.commit('setJwt', token)
}

export async function validateLogin () {
  try {
    if (localStorage.getItem('jwt')) {
      await renew(localStorage.getItem('jwt'))
    }
  } catch (_) {
    console.warn('Invalid JWT token in storage') // eslint-disable-line
  }
}

export async function login (serverUrl, subdomain) {
  
  const jwt = localStorage.getItem('jwt')
  if(!serverUrl && jwt) {
    parseToken(jwt)
  }
  let res = null
  if(serverUrl) {
    let url = '/ngssh_api/auth/renew?serverUrl='+serverUrl
    if(subdomain) {
      url += '&subdomain='+subdomain
    }
    res = await axios.get(url, {
      headers: {
        'X-Auth': jwt,
      }
    })
  } else {
    res = await axios.get('/ngssh_api/auth/renew', {
      headers: {
        'X-Auth': jwt,
      }
    })
  }
  
  const body = res.data

  if (res.status === 200) {
    parseToken(body)
  } else {
    throw new Error(body)
  }
}

export async function renew () {

  const res = await axios.get('/ngssh_api/auth/renew', {
    headers: {
      'X-Auth': store.state.jwt,
    }
  })

  const body = await res.data

  if (res.status === 200) {
    parseToken(body)
  } else {
    throw new Error(body)
  }
}

export function logout () {
  store.commit('setJwt', '')
  localStorage.setItem('jwt', null)
  window.location = getLoginUrl()
}

export function getLoginUrl() {

  let toLoginPath = process.env.VUE_APP_LOGIN_URL+"?mode=ng-sftp"
  if(store.state.user && store.state.user.subdomain) {
    toLoginPath = toLoginPath+'&subdomain='+store.state.user.subdomain
  }

  return toLoginPath
}

export function getQueryString (name) {
  const reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
  const r = window.location.search.substr(1).match(reg);
  if (r != null) return unescape(r[2]); return null;
}

export async function healthcheck(serverUrl, auth) {
  const tmp = axios.create()
  const healthcheck = serverUrl+'/tunnel_api/healthcheck?auth='+auth
  await tmp.get(healthcheck, {
    timeout: 1000*3,
  })
}