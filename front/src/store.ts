import throttle from 'lodash-es/throttle'
import {Md5 as md5} from 'ts-md5/dist/md5'
import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'

import VueInstance from './main'

Vue.use(Vuex)

const INFO_URL = '/info'

export interface InfoInterface {
  version: string,
}

export enum TileCategory {
  Health = 'HEALTH',
  Build = 'BUILD',
  Group = 'GROUP',
}

export enum TileType {
  HttpAny = 'HTTP-ANY',
  HttpRaw = 'HTTP-RAW',
  HttpJson = 'HTTP-JSON',
  HttpYaml = 'HTTP-YAML',
  Ping = 'PING',
  Port = 'PORT',
  Pingdom = 'PINGDOM',
  GitLab = 'GITLAB-BUILD',
  Travis = 'TRAVISCI-BUILD',
  Jenkins = 'JENKINS-BUILD',

  Empty = 'EMPTY',
  Group = 'GROUP',
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

export enum Theme {
  Default = 'DEFAULT',
  Dark = 'DARK',
}

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
  version: string | undefined,
  columns: number,
  tiles: TileConfig[],
  tilesState: { [key: string]: TileState },
}

const store: StoreOptions<RootState> = {
  state: {
    version: undefined,
    columns: 4,
    tiles: [],
    tilesState: {},
  },
  getters: {
    queryParams(): string[] {
      return window.location.search.substr(1).split('&')
    },
    configUrl(state, getters): string {
      let configUrl = ''
      const configQueryParam = getters.queryParams.find((queryParam: string) => {
        return /^config=/.test(queryParam)
      })
      if (configQueryParam) {
        configUrl = configQueryParam.substr(configQueryParam.indexOf('=') + 1)
      }

      return configUrl
    },
    theme(state, getters): Theme {
      let theme = Theme.Default
      const themeQueryParam = getters.queryParams.find((queryParam: string) => {
        return /^theme=/.test(queryParam)
      })
      if (themeQueryParam) {
        const queryTheme = themeQueryParam.substr(themeQueryParam.indexOf('=') + 1)
        if (Object.values(Theme).includes(queryTheme.toUpperCase())) {
          theme = queryTheme.toUpperCase()
        }
      }

      return theme
    },
  },
  mutations: {
    setVersion(state, payload: string): void {
      state.version = payload
    },
    setConfig(state, payload: ConfigInterface): void {
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
  },
  actions: {
    autoUpdate({commit, state, getters}) {
      return VueInstance.$http.get(getters.configUrl.replace(/\/config.*$/, INFO_URL))
        .then(async (data) => {
          const info: InfoInterface = await data.json()

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
      function setTileStateKey(tile: TileConfig) {
        // Create a random identifier
        tile.stateKey = tile.type + '_' + md5.hashStr(JSON.stringify(tile))

        // Set stateKey on group subTiles
        if (tile.tiles) {
          tile.tiles = tile.tiles.map(setTileStateKey)
        }

        return tile
      }

      return VueInstance.$http.get(getters.configUrl)
        .then(async (data) => {
          const config: ConfigInterface = await data.json()

          config.tiles = config.tiles.map(setTileStateKey)

          commit('setConfig', config)
        })
    },
    refreshTiles({commit, state}) {
      function timeout(delay: number = 0) {
        return new Promise(resolve => setTimeout(resolve, delay))
      }
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
      function refreshGroup(groupTile: TileConfig) {
        if (!groupTile.tiles) {
          return
        }

        const groupSubTilesState = groupTile.tiles
          .map((subTile) => subTile.stateKey)
          .map((subTileStateKey) => state.tilesState[subTileStateKey])

        const groupStatus = groupSubTilesState.reduce((worstSubTileStatus, subTileState) => {
          const subTileStatus = subTileState !== undefined ? getPreviousOrStatus(subTileState) : TileStatus.Unknown

          return mostImportantStatus(worstSubTileStatus, subTileStatus)
        }, TileStatus.Unknown)

        const groupSucceededSubTiles = groupSubTilesState.filter((subTileState) => {
          const subTileStatus = subTileState !== undefined ? getPreviousOrStatus(subTileState) : TileStatus.Unknown

          return subTileStatus === TileStatus.Success
        })

        const groupMessage = `${groupSucceededSubTiles.length} / ${groupTile.tiles.length}`

        const groupState = {
          category: TileCategory.Group,
          status: groupStatus,
          message: groupMessage,
        }

        commit('setTileState', {tileStateKey: groupTile.stateKey, tileState: groupState})
      }

      function getPreviousOrStatus(subTileState: TileState): TileStatus {
        let subTileStatus = subTileState.status

        if ([TileStatus.Queued, TileStatus.Running].includes(subTileState.status)) {
          subTileStatus = subTileState.previousStatus as TileStatus
        }

        return subTileStatus
      }

      function mostImportantStatus(status1: TileStatus, status2: TileStatus): TileStatus {
        return ORDERED_TILE_STATUS.indexOf(status1) < ORDERED_TILE_STATUS.indexOf(status2) ? status2 : status1
      }

      // Classic tiles (all except empty and group types)
      state.tiles
        .filter((tile) => !!tile.url)
        .forEach(async (tile) => {
          if (state.tilesState.hasOwnProperty(tile.stateKey)) {
            // Randomize delay for each tile to avoid DoS back-end services
            await timeout(Math.random() * 10000)
          }

          await refreshTile(tile)
        })

      // Group subTiles
      state.tiles.forEach(async (tile) => {
        if (!tile.tiles) {
          return
        }

        if (state.tilesState.hasOwnProperty(tile.stateKey)) {
          // Randomize delay for each group to avoid DoS back-end services
          await timeout(Math.random() * 10000)
        }

        const throttledRefreshGroup = throttle(refreshGroup, 150)
        tile.tiles.map(async (subTile) => {
          await refreshTile(subTile).then(() => {
            throttledRefreshGroup(tile)
          })
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
