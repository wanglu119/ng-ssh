import axios from "axios";

import store from "@/store";

import { procURL } from './utils'

export async function fetch (dirPath) {
  const sftpConfig = store.state.currTabConf;

  const opts = { 
    method: "POST",
    data: {
      name: sftpConfig.sshConfigName,
      dir_path: dirPath
    }
  }

  const res = await procURL("/ngssh_api/sftp/listFiles", opts);

  return res
}

export async function move (items, overwrite = false, rename = false) {
  const sftpConfig = store.state.currTabConf;

  const opts = { 
    method: "PUT",
    data: {
      name: sftpConfig.sshConfigName
    }
  }
  for (const item of items) {
    opts.data = {
      name: sftpConfig.sshConfigName,
      old_full_path: item.from,
      new_full_path: item.to,
    }
    await procURL(`/ngssh_api/sftp/rename?override=${overwrite}&rename=${rename}`, opts);
  }
}

export async function copy (items, overwrite = false, rename = false) {
  const sftpConfig = store.state.currTabConf;

  const opts = { 
    method: "POST",
    data: {
      name: sftpConfig.sshConfigName
    }
  }
  for (const item of items) {
    opts.data = {
      name: sftpConfig.sshConfigName,
      old_full_path: item.from,
      new_full_path: item.to,
    }
    await procURL(`/ngssh_api/sftp/copy?override=${overwrite}&rename=${rename}`, opts);
  }
}

export function download(format, isDir, ...files) {
  let url = `${axios.defaults.baseURL}/ngssh_api/sftp/raw`;

  const sftpConfig = store.state.currTabConf;
  if (files.length === 1) {
    url += files[0] + `?is_dir=${isDir}&name=${sftpConfig.sshConfigName}&`;
  } else {
    let arg = "";

    for (const file of files) {
      arg += file + ",";
    }

    arg = arg.substring(0, arg.length - 1);
    arg = encodeURIComponent(arg);
    url += `/?files=${arg}&is_dir=${isDir}&name=${sftpConfig.sshConfigName}&`;
  }

  if (format !== null) {
    url += `algo=${format}&`;
  }

  url += `subdomain=` + store.state.subdomain;
  url += `&auth=${store.state.jwt}`;

  console.log(url);
  window.open(url);
}

export async function post (url, content = '', overwrite = false, onupload) {
  const sftpConfig = store.state.currTabConf;
  return new Promise((resolve, reject) => {
    const request = new XMLHttpRequest()
    request.open('POST', `${axios.defaults.baseURL}/ngssh_api/sftp/resources${url}?name=${sftpConfig.sshConfigName}&override=${overwrite}&subdomain=`+store.state.subdomain, true)
    request.setRequestHeader('X-Auth', store.state.jwt)

    if (typeof onupload === 'function') {
      request.upload.onprogress = onupload
    }

    request.onload = () => {
      if (request.status === 200) {
        resolve(request.responseText)
      } else if (request.status === 409) {
        reject(request.status)
      } else {
        reject(request.responseText)
      }
    }

    request.onerror = (error) => {
      reject(error)
    }

    request.send(content)
  })
}

async function resourceAction (url, method, content) {

  const opts = { method }

  if (content) {
    // opts.body = content
    opts.data = content
  }
  const res = await procURL(`/ngssh_api/sftp/resources${url}`, opts)

  if (res.status !== 200) {
    throw new Error(res.responseText)
  } else {
    return res
  }
}

export async function removeFile (url) {
  const sftpConfig = store.state.currTabConf;
  const opts = { 
    method: "DELETE",
    data: {
      name: sftpConfig.sshConfigName,
      old_full_path: url
    }
  }

  const res = await procURL(`/ngssh_api/sftp/removeFile`, opts)

  return res
}

export async function removeDir (url) {
  const sftpConfig = store.state.currTabConf;
  const opts = { 
    method: "DELETE",
    data: {
      name: sftpConfig.sshConfigName,
      old_full_path: url
    }
  }

  const res = await procURL(`/ngssh_api/sftp/removeDir`, opts)

  return res
}