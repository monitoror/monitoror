import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'

Vue.use(Vuex)

export enum TileType {
  Ping = 'ping',
  Port = 'port',
  Travis = 'travis',

  Empty = 'empty',
  Group = 'group',
}

export interface TileState {
  type: TileType,
  label?: string,
  url?: string,
  subTiles?: TileState[],
}

interface RootState {
  version: string,
  columns: number,
  tiles: TileState[],
}

const store: StoreOptions<RootState> = {
  state: {
    version: 'unknown',
    columns: 4,
    tiles: [],
  },
  mutations: {

  },
  actions: {

  },
}

export default new Vuex.Store(store)
