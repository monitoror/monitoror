import axios from 'axios'
import throttle from 'lodash-es/throttle'
import {Md5 as md5} from 'ts-md5/dist/md5'
import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'

import DISPLAYABLE_SUBTILE_STATUS from '@/constants/displayableSubtileStatus'
import Theme from '@/enums/theme'
import TileStatus from '@/enums/tileStatus'
import getQueryParamValue from '@/helpers/getQueryParamValue'
import getSubTilePreviousOrStatus from '@/helpers/getSubTilePreviousOrStatus'
import mostImportantStatus from '@/helpers/mostImportantStatus'
import timeout from '@/helpers/timeout'
import Config from '@/interfaces/config'
import Info from '@/interfaces/info'
import TileConfig from '@/interfaces/tileConfig'
import TileState from '@/interfaces/tileState'

Vue.use(Vuex)

const API_BASE_PATH = '/api/v1'
const INFO_URL = '/info'

interface RootState {
  version: string | undefined,
  columns: number,
  tiles: TileConfig[],
  tilesState: { [key: string]: TileState },
  online: boolean,
}

const store: StoreOptions<RootState> = {
  state: {
    version: undefined,
    columns: 4,
    tiles: [],
    tilesState: {},
    online: true,
  },
  getters: {
    apiBaseUrl(): string {
      const defaultApiBaseUrl = window.location.origin
      let apiBaseUrl = getQueryParamValue('apiBaseUrl', defaultApiBaseUrl)

      apiBaseUrl = apiBaseUrl.replace(/\/+$/, '')

      return apiBaseUrl
    },
    configPath(): string {
      const configPath = getQueryParamValue('configPath')

      return configPath
    },
    configUrl(): string {
      const configUrl = getQueryParamValue('configUrl')

      return configUrl
    },
    proxyfiedConfigUrl(state, getters): string {
      const configProxyUrl = `${getters.apiBaseUrl}${API_BASE_PATH}/config`

      if (getters.configUrl !== '') {
        return `${configProxyUrl}?url=${getters.configUrl}`
      }

      if (getters.configPath !== '') {
        return `${configProxyUrl}?path=${getters.configPath}`
      }

      return ''
    },
    theme(): Theme {
      let theme = Theme.Default
      const queryTheme = getQueryParamValue('theme', theme)

      if (Object.values(Theme).includes(queryTheme.toUpperCase() as Theme)) {
        theme = queryTheme.toUpperCase() as Theme
      }

      return theme
    },
  },
  mutations: {
    setVersion(state, payload: string): void {
      state.version = payload
    },
    setConfig(state, payload: Config): void {
      state.columns = payload.columns
      state.tiles = payload.tiles
    },
    setTileState(state, payload: { tileStateKey: string, tileState: TileState }): void {
      if (!state.tilesState.hasOwnProperty(payload.tileStateKey)) {
        Vue.set(state.tilesState, payload.tileStateKey, payload.tileState)
      } else {
        state.tilesState[payload.tileStateKey] = payload.tileState
      }
    },
    setOnline(state, payload: boolean): void {
      state.online = payload
    },
  },
  actions: {
    autoUpdate({commit, state, getters}) {
      const infoUrl = getters.apiBaseUrl + API_BASE_PATH + INFO_URL

      return axios.get(infoUrl)
        .then((response) => {
          const info: Info = response.data

          if (state.version === undefined) {
            commit('setVersion', info.version)
            return
          }

          if (info.version !== state.version) {
            window.location.reload()
          }
        })
    },
    loadConfiguration({commit, getters}) {
      function hydrateTile(tile: TileConfig) {
        // Create a random identifier
        tile.stateKey = tile.type + '_' + md5.hashStr(JSON.stringify(tile))

        // Prefix URL with api base URL
        if (tile.url) {
          tile.url = getters.apiBaseUrl + tile.url
        }

        // Set stateKey on group subTiles
        if (tile.tiles) {
          tile.tiles = tile.tiles.map(hydrateTile)
        }

        return tile
      }

      return axios.get(getters.proxyfiedConfigUrl)
        .then((response) => {
          const config: Config = response.data

          config.tiles = config.tiles.map(hydrateTile)

          commit('setConfig', config)
        })
    },
    refreshTiles({state, dispatch}) {
      // Classic tiles (all except empty and group types)
      state.tiles
        .filter((tile) => !!tile.url)
        .forEach(async (tile) => {
          // Randomize delay for each tile to avoid DoS back-end services
          await timeout(Math.random() * 10000)

          await dispatch('refreshTile', tile)
        })

      // Group subTiles
      state.tiles.forEach(async (groupTile) => {
        if (!groupTile.tiles) {
          return
        }

        // Randomize delay for each group to avoid DoS back-end services
        await timeout(Math.random() * 10000)

        const throttledDispatch = throttle(dispatch, 150)
        groupTile.tiles.map(async (subTile) => {
          await dispatch('refreshTile', subTile).then(() => {
            throttledDispatch('refreshGroup', groupTile)
          })
        })
      })
    },
    refreshTile({commit}, tile: TileConfig): Promise<void> {
      if (!tile.url) {
        return Promise.resolve()
      }

      return axios.get(tile.url)
        .then(async (response) => {
          const tileState = response.data

          commit('setTileState', {tileStateKey: tile.stateKey, tileState})
        }) as Promise<void>
    },
    refreshGroup({state, commit}, groupTile: TileConfig) {
      if (!groupTile.tiles) {
        return
      }

      const groupSubTilesState = groupTile.tiles
        .map((subTile) => subTile.stateKey)
        .map((subTileStateKey) => state.tilesState[subTileStateKey])

      const groupStatus = groupSubTilesState.reduce((worstSubTileStatus, subTileState) => {
        const subTileStatus = getSubTilePreviousOrStatus(subTileState)

        return mostImportantStatus(worstSubTileStatus, subTileStatus)
      }, TileStatus.Unknown)

      const groupNonDisplayedSubTiles = groupSubTilesState.filter((subTileState) => {
        const subTileStatus = getSubTilePreviousOrStatus(subTileState)

        return !DISPLAYABLE_SUBTILE_STATUS.includes(subTileStatus)
      })

      const groupMessage = `${groupNonDisplayedSubTiles.length} / ${groupTile.tiles.length}`

      const groupState = {
        status: groupStatus,
        message: groupMessage,
      }

      commit('setTileState', {tileStateKey: groupTile.stateKey, tileState: groupState})
    },
    increaseTilesDuration({commit, state}) {
      Object.keys(state.tilesState).forEach((tileStateKey) => {
        const tileState = state.tilesState[tileStateKey]
        if (tileState.duration) {
          tileState.duration += 1
          commit('setTileState', {tileStateKey, tileState})
        }
      })
    },
    updateNetworkState({commit}) {
      commit('setOnline', navigator.onLine)
    },
  },
}

export default new Vuex.Store(store)
