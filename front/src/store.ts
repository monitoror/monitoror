import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'

import VueInstance from './main'

Vue.use(Vuex)

export enum TileCategory {
  Health = 'HEALTH',
  Build = 'BUILD',
  Group = 'GROUP',
}

export enum TileType {
  Ping = 'ping',
  Port = 'port',
  Travis = 'travis',

  Empty = 'empty',
  Group = 'group',
}

export enum TileStatus {
  Success = 'SUCCESS',
  Failed = 'FAILURE',
  Warning = 'WARNING',
  Running = 'RUNNING',
  Queued = 'QUEUED',
  Canceled = 'CANCELED',
  Unknown = 'UNKNOWN',
}

const ORDERED_TILE_STATUS = [
  TileStatus.Unknown,
  TileStatus.Success,
  TileStatus.Canceled,
  TileStatus.Warning,
  TileStatus.Failed,
  TileStatus.Queued,
  TileStatus.Running,
]

interface ConfigInterface {
  columns: number,
  tiles: TileConfig[],
}

export interface TileConfig {
  type: TileType,
  label?: string,
  columnSpan?: number,
  rowSpan?: number,
  url?: string,
  stateKey: string,
  tiles?: TileConfig[],
}

export interface TileAuthor {
  name: string,
  avatarUrl: string,
}

export interface TileState {
  category: TileCategory,
  label?: string,
  status: TileStatus,
  previousStatus?: TileStatus,
  message?: string,
  author?: TileAuthor,
  duration?: number,
  estimatedDuration?: number,
  startedAt?: number,
  finishedAt?: number,
}

interface RootState {
  version: string,
  columns: number,
  tiles: TileConfig[],
  tilesState: { [key: string]: TileState },
}

const store: StoreOptions<RootState> = {
  state: {
    version: 'unknown',
    columns: 4,
    tiles: [],
    tilesState: {},
  },
  getters: {
    configUrl(): string {
      let configUrl = ''
      const configQueryParam = window.location.search.substr(1).split('&').find((queryParam: string) => {
        return /^config=/.test(queryParam)
      })
      if (configQueryParam) {
        configUrl = configQueryParam.substr(configQueryParam.indexOf('=') + 1)
      }

      return configUrl
    },
  },
  mutations: {
    setConfig(state, payload: ConfigInterface): void {
      state.columns = payload.columns
      state.tiles = payload.tiles
      state.tilesState = {}
    },
    setTileState(state, payload: { tileStateKey: string, tileState: TileState }): void {
      if (!state.tilesState.hasOwnProperty(payload.tileStateKey)) {
        Vue.set(state.tilesState, payload.tileStateKey, payload.tileState)
      } else {
        state.tilesState[payload.tileStateKey] = payload.tileState
      }
    },
  },
  actions: {
    loadConfig({commit, getters}) {
      return VueInstance.$http.get(getters.configUrl)
        .then(async (data) => {
          const config: ConfigInterface = await data.json()

          config.tiles = config.tiles.map((tile) => {
            // Create a random identifier
            tile.stateKey = tile.type + '_' + Math.random().toString(36).substr(2, 9)

            return tile
          })

          commit('setConfig', config)
        })
    },
    refreshTiles({commit, state}) {
      function refreshTile(tile: TileConfig): Promise<void> {
        if (!tile.url) {
          return Promise.resolve()
        }

        return VueInstance.$http.get(tile.url)
          .then(async (data) => {
            const tileState = await data.json()
            commit('setTileState', {tileStateKey: tile.stateKey, tileState})
          }) as Promise<void>
      }

      function mostImportantStatus(status1: TileStatus, status2: TileStatus) {
        return ORDERED_TILE_STATUS.indexOf(status1) < ORDERED_TILE_STATUS.indexOf(status2) ? status2 : status1
      }

      // Classic tiles (all except empty and group types)
      state.tiles
        .filter((tile) => !!tile.url)
        .forEach(refreshTile)

      // Group subTiles
      state.tiles.forEach((tile) => {
        if (!tile.tiles) {
          return
        }

        Promise.all(tile.tiles.map(refreshTile)).then(() => {
          if (!tile.tiles) {
            return
          }

          const groupSubTilesState = tile.tiles
            .map((subTile) => subTile.stateKey)
            .map((subTileStateKey) => state.tilesState[subTileStateKey])

          const groupStatus = groupSubTilesState.reduce((worstSubTileStatus, subTileState) => {
            let subTileStatus = subTileState.status

            if ([TileStatus.Queued, TileStatus.Running].includes(subTileState.status)) {
              subTileStatus = subTileState.previousStatus as TileStatus
            }

            return mostImportantStatus(worstSubTileStatus, subTileStatus)
          }, TileStatus.Unknown)

          const groupSucceededSubTiles = groupSubTilesState.filter((subTileState) => {
            let subTileStatus = subTileState.status

            if ([TileStatus.Queued, TileStatus.Running].includes(subTileState.status)) {
              subTileStatus = subTileState.previousStatus as TileStatus
            }

            return subTileStatus === TileStatus.Success
          })

          const groupMessage = `${groupSucceededSubTiles.length} / ${tile.tiles.length}`

          const groupState = {
            category: TileCategory.Group,
            status: groupStatus,
            message: groupMessage,
          }

          commit('setTileState', {tileStateKey: tile.stateKey, tileState: groupState})
        })
      })
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
  },
}

export default new Vuex.Store(store)
