import axios from 'axios'
import {Md5 as md5} from 'ts-md5/dist/md5'
import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'

import DISPLAYABLE_SUBTILE_STATUS from '@/constants/displayableSubtileStatus'

import Task from '@/classes/task'
import TaskInterval from '@/enums/taskInterval'
import TaskType from '@/enums/taskType'
import Theme from '@/enums/theme'
import TileStatus from '@/enums/tileStatus'
import getQueryParamValue from '@/helpers/getQueryParamValue'
import getSubTilePreviousOrStatus from '@/helpers/getSubTilePreviousOrStatus'
import mostImportantStatus from '@/helpers/mostImportantStatus'
import Config from '@/interfaces/config'
import Info from '@/interfaces/info'
import TileConfig from '@/interfaces/tileConfig'
import TileState from '@/interfaces/tileState'
import {now} from 'lodash-es'

Vue.use(Vuex)

const API_BASE_PATH = '/api/v1'
const INFO_URL = '/info'

export interface RootState {
  version: string | undefined,
  columns: number,
  tiles: TileConfig[],
  tilesState: { [key: string]: TileState },
  tasks: Task[],
  online: boolean,
}

const store: StoreOptions<RootState> = {
  state: {
    version: undefined,
    columns: 4,
    tiles: [],
    tilesState: {},
    tasks: [],
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
    tileStateKeys(state): string[] {
      const tileStateKeys: string[] = []
      state.tiles.forEach((tile: TileConfig) => {
        tileStateKeys.push(tile.stateKey)

        // Add group subTiles stateKeys
        if (tile.tiles) {
          tile.tiles.forEach((subTile) => {
            tileStateKeys.push(subTile.stateKey)
          })
        }
      })

      return tileStateKeys
    },
    taskIds(state): string[] {
      return state.tasks.map((task: Task) => task.id)
    },
    loadedTilesCount(state): number {
      const loadedTilesCount = Object.keys(state.tilesState).length

      return loadedTilesCount
    },
    loadableTilesCount(state): number {
      const loadableTilesStateKeys: string[] = []

      function addLoadableTileStateKey(stateKey: string) {
        if (loadableTilesStateKeys.includes(stateKey)) {
          return
        }

        loadableTilesStateKeys.push(stateKey)
      }

      state.tiles.forEach((tile) => {
        if (tile.url) {
          addLoadableTileStateKey(tile.stateKey)
        }

        if (tile.tiles) {
          addLoadableTileStateKey(tile.stateKey)
          tile.tiles.forEach((subTile) => {
            addLoadableTileStateKey(subTile.stateKey)
          })
        }
      })

      const loadableTilesCount = loadableTilesStateKeys.length

      return loadableTilesCount
    },
    loadingProgress(state, getters): number {
      const loadingProgress = getters.loadedTilesCount / getters.loadableTilesCount

      if (!loadingProgress) {
        return 0
      }

      return loadingProgress
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
    setTasks(state, payload: Task[]): void {
      state.tasks = payload
    },
    addTask(state, payload: Task): void {
      state.tasks.push(payload)
    },
  },
  actions: {
    async autoUpdate({commit, state, getters}) {
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
    async loadConfiguration({commit, state, getters, dispatch}) {
      function hydrateTile(tile: TileConfig, groupTile?: TileConfig) {
        // Create a identifier based on tile configuration
        tile.stateKey = tile.type + '_' + md5.hashStr(JSON.stringify(tile))

        if (tile.url) {
          // Prefix URL with api base URL
          tile.url = getters.apiBaseUrl + tile.url

          // Create a task for this tile
          const initialDelay = Math.random() * (tile.initialMaxDelay || 0)
          const refreshTileTask = new Task(
            tile.stateKey,
            TaskType.RefreshTile,
            async () => {
              await dispatch('refreshTile', tile)

              if (groupTile !== undefined) {
                dispatch('refreshGroup', groupTile)
              }
            },
            10 * TaskInterval.Second,
            initialDelay,
          )
          dispatch('addTask', refreshTileTask)
        }

        // Set stateKey on group subTiles
        if (tile.tiles) {
          tile.tiles = tile.tiles.map((subTile) => hydrateTile(subTile, tile))
        }

        return tile
      }

      return axios.get(getters.proxyfiedConfigUrl)
        .then((response) => {
          const config: Config = response.data

          // Kill old refreshTile tasks
          state.tasks
            .filter((task) => task.type === TaskType.RefreshTile && !getters.tileStateKeys.includes(task.id))
            .map((task) => task.kill())

          config.tiles = config.tiles.map((tile) => hydrateTile(tile))

          commit('setConfig', config)
        })
    },
    async refreshTile({commit}, tile: TileConfig) {
      if (!tile.url) {
        return Promise.resolve()
      }

      return axios.get(tile.url)
        .then(async (response) => {
          const tileState = response.data

          commit('setTileState', {tileStateKey: tile.stateKey, tileState})
        })
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
        if (tileState.build !== undefined && tileState.build.duration) {
          tileState.build.duration += 1
          commit('setTileState', {tileStateKey, tileState})
        }
      })
    },
    updateNetworkState({commit}) {
      commit('setOnline', navigator.onLine)
    },
    addTask({getters, commit}, task: Task) {
      // Avoid adding multiple task with the same ID
      if (getters.taskIds.includes(task.id)) {
        return
      }

      commit('addTask', task)
    },
    runTasks({state, dispatch}) {
      const nowTime = now()
      const shouldRunTask = (task: Task) => !task.isDone() && task.time <= nowTime

      const taskToRun = state.tasks.filter(shouldRunTask)
      Promise.all(taskToRun.map((task: Task) => {
        return task.run()
      })).then(() => {
        dispatch('updateTasks')
      })
    },
    updateTasks({commit, state}) {
      const {taskList, hasChanged} = state.tasks.reduce(
        (previousValue: { taskList: Task[], hasChanged: boolean }, task: Task) => {
          // Remove dead tasks from task list
          if (task.isDead()) {
            previousValue.hasChanged = true
            return previousValue
          }

          // Update outdated recurring tasks
          if (task.isDone()) {
            task.prepareNextRun()
          }

          previousValue.taskList.push(task)

          return previousValue
        },
        {taskList: [], hasChanged: false},
      )

      if (hasChanged) {
        commit('setTasks', taskList)
      }
    },
    killAllTasks({commit, state}) {
      state.tasks.map((task: Task) => task.kill())
      commit('setTasks', [])
    },
  },
}

export default new Vuex.Store(store)
