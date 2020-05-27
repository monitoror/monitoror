import axios from 'axios'
import {now} from 'lodash-es'
import {Md5 as md5} from 'ts-md5/dist/md5'
import Vue from 'vue'
import Vuex, {StoreOptions} from 'vuex'

import DISPLAYABLE_SUBTILE_STATUS from '@/constants/displayableSubtileStatus'

import Task from '@/classes/task'
import ConfigErrorId from '@/enums/configErrorId'
import Route from '@/enums/route'
import TaskInterval from '@/enums/taskInterval'
import TaskType from '@/enums/taskType'
import Theme from '@/enums/theme'
import TileStatus from '@/enums/tileStatus'
import getQueryParamValue from '@/helpers/getQueryParamValue'
import getSubTilePreviousOrStatus from '@/helpers/getSubTilePreviousOrStatus'
import mostImportantStatus from '@/helpers/mostImportantStatus'
import Config from '@/interfaces/config'
import ConfigBag from '@/interfaces/configBag'
import ConfigError from '@/interfaces/configError'
import ConfigMetadata from '@/interfaces/configMetadata'
import Info from '@/interfaces/info'
import TaskOptions from '@/interfaces/taskOptions'
import TileConfig from '@/interfaces/tileConfig'
import TileState from '@/interfaces/tileState'

Vue.use(Vuex)

const API_BASE_PATH = '/api/v1'
export const DEFAULT_CONFIG_NAME = 'default'
const INFO_URL = '/info'
const QUERY_PARAM_KEYS = {
  API_BASE_URL: 'apiBaseUrl',
  CONFIG: 'config',
  THEME: 'theme',
}

export interface RootState {
  appVersion: string | undefined,
  configVersion: string | undefined,
  columns: number,
  zoom: number,
  tiles: TileConfig[],
  tilesState: { [key: string]: TileState },
  tasks: Task[],
  errors: ConfigError[],
  online: boolean,
  now: Date,
  lastRefreshDate: Date,
  configList: ConfigMetadata[],
}

const store: StoreOptions<RootState> = {
  state: {
    appVersion: undefined,
    configVersion: undefined,
    columns: 4,
    zoom: 1,
    tiles: [],
    tilesState: {},
    tasks: [],
    errors: [],
    online: true,
    now: new Date(),
    lastRefreshDate: new Date(),
    configList: [],
  },
  getters: {
    apiBaseUrl(): string {
      const defaultApiBaseUrl = window.location.origin + window.location.pathname
      let apiBaseUrl = getQueryParamValue(QUERY_PARAM_KEYS.API_BASE_URL, defaultApiBaseUrl) as string

      apiBaseUrl = apiBaseUrl.replace(/\/+$/, '')

      return apiBaseUrl
    },
    configParam(): string {
      return getQueryParamValue(QUERY_PARAM_KEYS.CONFIG, DEFAULT_CONFIG_NAME) as string
    },
    configProxyUrl(state, getters): string {
      return `${getters.apiBaseUrl}${API_BASE_PATH}/configs`
    },
    proxyfiedConfigUrl(state, getters): string | undefined {
      const urlEncodedConfigParam = encodeURIComponent(getters.configParam)

      return `${getters.configProxyUrl}/${urlEncodedConfigParam}`
    },
    currentRoute(): Route | undefined {
      const queryHash = window.location.hash.replace(/^#/, '')
      let route

      if (Object.values(Route).includes(queryHash.toLowerCase() as Route)) {
        route = queryHash.toLowerCase() as Route
      }

      return route
    },
    theme(): Theme {
      let theme = Theme.Default
      const queryTheme = getQueryParamValue(QUERY_PARAM_KEYS.THEME, theme) as string

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
    hasUnknownDefaultConfigError(state, getters): boolean {
      if (state.errors.length === 0) {
        return false
      }

      const error = state.errors[0]

      const isDefaultConfig = getters.configParam === DEFAULT_CONFIG_NAME
      const isUnknownNamedConfigError = error.id === ConfigErrorId.UnknownNamedConfig

      return isDefaultConfig && isUnknownNamedConfigError
    },
    isNewUser(state, getters): boolean {
      if (!getters.hasUnknownDefaultConfigError) {
        return false
      }

      const error = state.errors[0]
      const hasNamedConfig = error.data.expected !== undefined

      return !hasNamedConfig
    },
    shouldShowWelcomePage(state, getters): boolean {
      const welcomePageRoutes = [
        Route.Welcome,
        Route.ChooseConfiguration,
      ]
      const isOnWelcomePageRoute = welcomePageRoutes.includes(getters.currentRoute)

      return isOnWelcomePageRoute || getters.isNewUser || getters.hasUnknownDefaultConfigError
    },
    shouldInit(state, getters): boolean {
      const isOnChooseConfigurationPageRoute = Route.ChooseConfiguration === getters.currentRoute

      return !isOnChooseConfigurationPageRoute
    },
  },
  mutations: {
    setAppVersion(state, payload: string): void {
      state.appVersion = payload
    },
    setConfig(state, payload: Config): void {
      state.configVersion = payload.version
      state.columns = payload.columns
      if (payload.zoom !== undefined) {
        state.zoom = payload.zoom
      }
      state.tiles = payload.tiles
    },
    setConfigList(state, payload: ConfigMetadata[]): void {
      state.configList = payload
    },
    setErrors(state, payload: ConfigError[]): void {
      state.errors = payload
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
    setLastRefreshDate(state, payload: Date): void {
      state.lastRefreshDate = payload
    },
    setNow(state, payload: Date): void {
      state.now = payload
    },
  },
  actions: {
    async autoUpdate({commit, state, getters}) {
      const infoUrl = getters.apiBaseUrl + API_BASE_PATH + INFO_URL

      return axios.get(infoUrl)
        .then((response) => {
          const info: Info = response.data

          if (state.appVersion === undefined) {
            commit('setAppVersion', info.version)
            return
          }

          if (info.version !== state.appVersion) {
            window.location.reload()
          }
        })
    },
    async fetchConfigList({commit, getters}) {
      return axios.get(getters.configProxyUrl)
        .then((response) => {
          const configList: ConfigMetadata[] = response.data.map((configMetadata: ConfigMetadata) => {
            let uiUrl = `?${QUERY_PARAM_KEYS.CONFIG}=${configMetadata.name}`

            if (getters.theme !== Theme.Default) {
              uiUrl += `&${QUERY_PARAM_KEYS.THEME}=${getters.theme.toLowerCase()}`
            }

            configMetadata.uiUrl = uiUrl

            return configMetadata
          }).sort((configMetadataA: ConfigMetadata, configMetadataB: ConfigMetadata) => {
            return configMetadataA.name.localeCompare(configMetadataB.name)
          })

          commit('setConfigList', configList)
        })
    },
    async fetchConfiguration({commit, state, getters, dispatch}) {
      const hydrateTile = (tile: TileConfig, groupTile?: TileConfig) => {
        // Create a identifier based on tile configuration
        tile.stateKey = tile.type + '_' + md5.hashStr(JSON.stringify(tile))

        if (tile.url) {
          // Prefix URL with api base URL
          tile.url = getters.apiBaseUrl + tile.url

          // Create a task for this tile
          dispatch('createRefreshTileTask', {tile, groupTile})
        }

        // Set stateKey on group subTiles
        if (tile.tiles) {
          tile.tiles = tile.tiles.map((subTile) => hydrateTile(subTile, tile))
        }

        return tile
      }

      commit('setLastRefreshDate', new Date())

      return axios.get(getters.proxyfiedConfigUrl)
        .then((response) => {
          const configBag: ConfigBag = response.data

          // Kill old refreshTile tasks
          state.tasks
            .filter((task) => task.type === TaskType.RefreshTile && !getters.tileStateKeys.includes(task.id))
            .map((task) => task.kill())

          if (configBag.errors !== undefined) {
            commit('setErrors', configBag.errors)
          } else {
            commit('setErrors', [])

            if (configBag.config !== undefined) {
              configBag.config.tiles = configBag.config.tiles.map((tile) => hydrateTile(tile))
              commit('setConfig', configBag.config)
            }
          }
        })
    },
    createRefreshTileTask({dispatch}, {tile, groupTile}: { tile: TileConfig, groupTile?: TileConfig }) {
      dispatch('addTask', {
        id: tile.stateKey,
        type: TaskType.RefreshTile,
        executor: async () => {
          await dispatch('refreshTile', tile)

          if (groupTile !== undefined) {
            await dispatch('refreshGroup', groupTile)
          }
        },
        interval: 10 * TaskInterval.Second,
        initialDelay: Math.random() * (tile.initialMaxDelay || 0),
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
    updateNetworkState({commit}) {
      commit('setOnline', navigator.onLine)
    },
    addTask({getters, commit}, taskOptions: TaskOptions) {
      // Avoid adding multiple task with the same ID
      if (getters.taskIds.includes(taskOptions.id)) {
        return
      }

      commit('addTask', new Task(taskOptions))
    },
    runTasks({state, dispatch}) {
      const nowTime = now()
      const shouldRunTask = (task: Task) => !task.isDone() && !task.isRunning() && task.time <= nowTime

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
    init({state, commit, dispatch}) {
      // Run auto-update each minute
      dispatch('addTask', {
        id: 'autoUpdate',
        type: TaskType.Root,
        executor: async () => {
          await dispatch('autoUpdate')
        },
        interval: 1 * TaskInterval.Minute,
      })

      // Fetch configuration each minute
      dispatch('addTask', {
        id: 'fetchConfiguration',
        type: TaskType.Root,
        executor: async () => {
          await dispatch('fetchConfiguration')
        },
        interval: 1 * TaskInterval.Minute,
        retryOnFailInterval: 5 * TaskInterval.Second,
        onFailedCallback: () => {
          // When offline, we know why we cannot get a response from the Core
          if (!state.online) {
            return
          }

          const configCannotBeFetch: ConfigError = {
            id: ConfigErrorId.CannotBeFetched,
            message: 'Configuration cannot be fetch from Monitoror Core',
            data: {},
          }
          commit('setErrors', [configCannotBeFetch])
        },
      })

      // Update "now" each second
      dispatch('addTask', {
        id: 'updateNow',
        type: TaskType.Root,
        executor: async () => {
          commit('setNow', new Date())
        },
        interval: 1 * TaskInterval.Second,
      })
    },
  },
}

export default new Vuex.Store(store)
