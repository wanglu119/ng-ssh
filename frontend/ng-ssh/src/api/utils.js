import store from '@/store'
import { renew } from '@/utils/auth'

import axios from 'axios'

// post put delete method
export async function procURL (url, opts) {

    if(url.indexOf('subdomain') < 0 && url.indexOf('?') < 0) {
        url = url+'?subdomain='+store.state.subdomain
    }

    opts = opts || {}
    opts.headers = opts.headers || {}
  
    const { headers,token,data,method, ...rest } = opts
    try {
        const res = await axios(`${url}`, {
            method: method,
            headers: {
              'X-Auth': token?token:store.state.jwt,
              ...headers
            },
            data: data,
            ...rest
          })
        
        if (!res.status === 200) {
            store.commit('setError',{code:res.status, msg:res.statusText})
            return res.data 
        }
        if (res.headers['X-Renew-Token'] === 'true') {
            await renew(store.state.jwt)
        }
        
        store.commit('setError',{code:res.status, msg:res.statusText})
        return res.data
    }catch(e) {
        if(e) {
            if(e.response) {
                store.commit('setError', {code:e.response.status,msg:e.response.statusText})
                throw new Error(e.response.status)
            } else {
                store.commit('setError', {code:999,msg:'unknown error'})
                throw new Error(404)
            }
        } 
    }
}

export function getImageUrl(img) {
    return process.env.VUE_APP_STATIC_URL+img
}