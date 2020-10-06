
import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);


export const ttyTabStore = new Vuex.Store({
    state: {
        ttys: [],
        currTtyId: null,
        geometry: null,
    },
    mutations: {
        setCurrTtyId(state, id) {
            state.currTtyId = id 
        },
        addTty(state, tty) {
            state.ttys.push(tty)
        },
        removeTty(state, id) {
            for(const i in state.ttys) {
                if(state.ttys[i].id === id) {
                    state.ttys.splice(i,1)
                    break
                }
            }
        },
        setGeometry(state, geometry) {
            state.geometry = geometry
        }
    },
    actions: {
        setCurrTtyId({commit, state}, id) {
            commit('setCurrTtyId', id)
        },
        addTty({commit, state}, tty) {
            commit('addTty', tty)
        },
        removeTty({commit, state}, id) {
            commit('removeTty',id)
        },
        setGeometry({commit, state},geometry) {
            // console.log('tty_store setGeometry: ', geometry)
            commit('setGeometry',geometry)
        }
    },
    getters: {
        getCurrTtyId(state) {
            return state.currTtyId
        },
        getGeometry(state) {
            return state.geometry
        }
    }

})

