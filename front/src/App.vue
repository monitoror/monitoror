<template>
  <div id="app" :style="cssProperties">
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
    private static readonly refreshTilesDelta: number = 10 // Each 10 seconds

    /*
     * Data
     */

    private refreshTilesCount: number = 0
    private refreshTilesInterval!: number

    /*
     * Computed
     */

    get cssProperties() {
      return {
        '--columns': this.columns,
        '--rows': Math.ceil(this.tiles.length / this.columns),
      }
    }

    get columns(): number {
      return this.$store.state.columns
    }

    get tiles(): TileConfig[] {
      return this.$store.state.tiles
    }

    /*
     * Hooks
     */

    private async mounted() {
      await Vue.nextTick()
      await this.$store.dispatch('loadConfig')
      await this.$store.dispatch('refreshTiles')

      this.refreshTilesInterval = setInterval(() => {
        if (this.refreshTilesCount >= App.refreshTilesDelta) {
          this.refreshTilesCount = 0
          return this.$store.dispatch('refreshTiles')
        }

        this.$store.dispatch('increaseTilesDuration')
        this.refreshTilesCount += 1
      }, 1000)
    }

    private beforeDestroy() {
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
  }

  .c-app--tiles-container {
    display: grid;
    grid-template-columns: repeat(var(--columns), 1fr);
    grid-gap: 6px;
    grid-auto-rows: calc((100vh - 6px) / var(--rows) - 6px);
    margin: 6px;
  }
</style>
