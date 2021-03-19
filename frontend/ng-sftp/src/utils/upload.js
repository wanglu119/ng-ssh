import store from '@/store'
import url from '@/utils/url'

export function checkConflict(files, items) {
  if (typeof items === "undefined" || items === null) {
    items = [];
  }

  const folder_upload = files[0].path !== undefined;

  let conflict = false;
  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    let name = file.name;

    if (folder_upload) {
      const dirs = file.path.split("/");
      if (dirs.length > 1) {
        name = dirs[0];
      }
    }

    const res = items.findIndex(function hasConflict(element) {
      // this is param name
      return element.name === this;
    }, name);

    if (res >= 0) {
      conflict = true;
      break;
    }
  }

  return conflict;
}

export function handleFiles(files, path, overwrite = false) {
  for (let i = 0; i < files.length; i++) {
    const file = files[i]

    const filename = (file.path !== undefined) ? file.path : file.name
    const filenameEncoded = url.encodeRFC5987ValueChars(filename)

    const id = store.state.upload.id

    let itemPath = path +"/"+ filenameEncoded

    if (file.isDir) {
      itemPath = path
      const folders = file.path.split("/")

      for (let i = 0; i < folders.length; i++) {
        const folder = folders[i]
        const folderEncoded = encodeURIComponent(folder)
        itemPath += folderEncoded + "/"
      }
    }

    const item = {
      id,
      path: itemPath,
      file,
      overwrite
    }

    store.dispatch('upload/upload', item);
  }
}