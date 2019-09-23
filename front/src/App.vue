<template>
  <div id="app" :class="classes" :style="cssProperties" @mousemove="resetShowCursorTimeout">
    <div class="c-app--tiles-container">
      <monitoror-tile v-for="tileConfig in tiles" :key="tileConfig.stateKey" :config="tileConfig"></monitoror-tile>
    </div>
  </div>
</template>

<script lang="ts">
  import {Component, Vue} from 'vue-property-decorator'

  import {TileConfig} from '@/store'
  import MonitororTile from './components/Tile.vue'

  @Component({
    components: {
      MonitororTile,
    },
  })
  export default class App extends Vue {
    private static readonly SHOW_CURSOR_DELAY: number = 10 // 10 seconds
    private static readonly REFRESH_TILES_DELTA: number = 10 // Each 10 seconds

    /*
     * Data
     */

    private showCursor: boolean = true
    private showCursorTimeout!: number
    private autoUpdateInterval!: number
    private loadConfigurationInterval!: number
    private refreshTilesCount: number = 0
    private refreshTilesInterval!: number

    /*
     * Computed
     */

    get classes() {
      return {
        'c-app__show-cursor': this.showCursor,
      }
    }

    get cssProperties() {
      const tilesCount = this.tiles.reduce((accumulator, tile) => {
        return accumulator + (tile.rowSpan || 1) * (tile.columnSpan || 1)
      }, 0)

      return {
        '--columns': this.columns,
        '--rows': Math.ceil(tilesCount / this.columns),
      }
    }

    get columns(): number {
      return this.$store.state.columns
    }

    get tiles(): TileConfig[] {
      return this.$store.state.tiles
    }

    /*
     * Methods
     */

    private resetShowCursorTimeout() {
      clearTimeout(this.showCursorTimeout)
      this.showCursor = true
      this.showCursorTimeout = setTimeout(() => {
        this.showCursor = false
      }, App.SHOW_CURSOR_DELAY * 1000)
    }

    /*
     * Hooks
     */

    private async mounted() {
      await Vue.nextTick()
      await this.$store.dispatch('autoUpdate')
      await this.$store.dispatch('loadConfiguration')
      await this.$store.dispatch('refreshTiles')

      this.autoUpdateInterval = setInterval(async () => {
        await this.$store.dispatch('autoUpdate')
      }, 60000)

      this.loadConfigurationInterval = setInterval(async () => {
        await this.$store.dispatch('loadConfiguration')
        await this.$store.dispatch('refreshTiles')
      }, 10000)

      this.refreshTilesInterval = setInterval(() => {
        if (this.refreshTilesCount >= App.REFRESH_TILES_DELTA) {
          this.refreshTilesCount = 0
          return this.$store.dispatch('refreshTiles')
        }

        this.$store.dispatch('increaseTilesDuration')
        this.refreshTilesCount += 1
      }, 1000)
    }

    private beforeDestroy() {
      clearInterval(this.autoUpdateInterval)
      clearInterval(this.loadConfigurationInterval)
      clearInterval(this.refreshTilesInterval)
    }
  }
</script>

<style lang="scss">
  #app {
    height: 100%;
    width: 100%;

    --columns: 1;
    --rows: 1;

    &:not(.c-app__show-cursor) {
      cursor: none;
    }
  }

  .c-app--tiles-container {
    display: grid;
    grid-template-columns: repeat(var(--columns), 1fr);
    grid-gap: 6px;
    grid-auto-rows: calc((100vh - 6px) / var(--rows) - 6px);
    margin: 6px;
  }
</style>
