// state
const state = {
  show: null,
  showMessage: null,
  showConfirm: null,
  multiple: false,
  selected: [],
  req: {},
  currPath: null,
};

// mutations
const mutations = {
  setReq: (state: any, {tabId, data}: any) => {
    state.req[tabId] = data;
  },
  removeReq: (state: any, tabId: any) => {
    delete state.req[tabId]
  },
  closeHovers: (state: any) => {
    state.show = null;
    state.showMessage = null;
    state.multiple = false;
  },
  showHover: (state: any, value: any) => {
    if (typeof value !== "object") {
      state.show = value;
      return;
    }
    
    state.show = value.prompt
    state.showMessage = value.message
    state.showConfirm = value.confirm
  },
  setMultiple: (state: any, value: any) => {
    state.multiple = value;
  },
  resetSelected(state: any) {
    state.selected = [];
  },
  removeSelected(state: any, index: any) {
    const i = state.selected.indexOf(index);
    state.selected.splice(i, 1);
  },
  addSelected(state: any, index: any) {
    state.selected.push(index);
  },
  setCurrPath(state: any, currPath: any) {
    state.currPath = currPath
  },
};

// getter
const getters = {
  selectedCount: (state: any) => state.selected.length
};

export default {
  namespaced: true,
  state: state,
  getters: getters,
  mutations: mutations
};
