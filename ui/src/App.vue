<template>
  <div id="app" :class="classes" :style="cssProperties" @mousemove="resetShowCursorTimeout">
    <div class="c-app--tiles-container">
      <monitoror-tile v-for="tileConfig in tiles" :key="tileConfig.stateKey" :config="tileConfig"></monitoror-tile>
    </div>

    <transition name="loading-fade">
      <div class="c-app--loading" v-if="showLoading" :class="appLoadingClasses">
        <div class="c-app--loading-container">
          <div class="c-app--loading-progress">
            <div class="c-app--loading-progress-bar" :style="loadingProgressBarStyle"></div>
          </div>
          <div class="c-app--logo">
            <svg viewBox="0 0 1220 300" xmlns="http://www.w3.org/2000/svg">
              <defs>
                <linearGradient id="monitoror-logo-line-gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                  <stop offset="0%" stop-color="var(--color-logo-line-start, #3d5a61)"/>
                  <stop offset="100%" stop-color="var(--color-logo-background, #87d7af)"/>
                </linearGradient>
                <mask id="monitoror-logo-line-mask">
                  <path fill="#fff"
                        d="M1190.18 188.973l2.522 12h-49.334v-12h46.812zm-79.771 21.881l-14.452 31.169-1.094 1.609c-1.578 1.196-3.466 2.099-5.342 1.785-1.912-.321-3.346-1.756-4.467-3.429l-18.692-41.015h-25.556c-1.233-1.172-2.334-2.533-3.326-4.066a36.452 36.452 0 01-3.844-7.934h36.585l1.692.243c.518.237 1.075.402 1.554.711 1.3.836 1.405 1.297 2.214 2.558l14.879 32.649 15.155-32.685.926-1.425c.43-.369.811-.803 1.289-1.108 1.294-.827 1.754-.732 3.229-.943h14.173v12h-10.342l-3.044 6.565-1.537.264v3.052zm-95.192-21.881c1.205 4.333 2.915 8.338 5.166 12H974.46v-12h40.757zm-58.793 12h-2.119l-1.717-.251c-1.378-.641-1.859-.682-2.871-1.883-.372-.442-.609-.983-.914-1.474l-17.081-39.288c-1.493-9.812-5.139-18.201-10.943-25.17l-10.023-23.054-37.946 87.507-.914 1.476-1.095.972c-.827-.9-1.573-1.87-2.229-2.901-1.885-2.91-3.368-6.153-4.422-9.739l41.095-94.767 1.077-1.663c.52-.407.98-.905 1.561-1.221a5.994 5.994 0 013.854-.647c1.293.216 2.386 1.044 3.436 1.864l1.079 1.662 40.172 92.401v16.176zm-110.115-12c1.204 4.333 2.915 8.338 5.165 12h-35.087c-.782-1.839-1.17-4.232-1.17-7.179v-4.821h31.092zm-49.128 0v4.821c0 2.578.161 4.969.502 7.179h-33.616v-12h33.114zm-91.838 0h40.688v12h-40.688v-12zm-95.306 29.505l-7.715-17.505h-8.4c2.223-3.666 3.951-7.666 5.145-12h7.167l1.71.248c.522.242 1.085.411 1.568.726 1.309.854 1.409 1.323 2.213 2.606l6.128 13.906-7.816 1.343v10.676zm39.614-17.505l-18.68 40.3-1.101 1.617c-1.048.789-2.14 1.589-3.421 1.788-1.97.306-4.034-.446-5.342-1.871-.441-.481-.714-1.092-1.07-1.638l-9.622-21.831h13.114l2.136 4.847 2.247-4.847 12.464-26.889.926-1.425c.429-.369.81-.803 1.288-1.108l.406-.224.046-.027c.967-.558 1.479-.506 2.777-.692h41.488v12h-37.656zm-69.038-12a36.343 36.343 0 01-3.87 7.934c-.999 1.533-2.101 2.894-3.33 4.066h-32.574c-1.234-1.172-2.334-2.533-3.327-4.066a36.449 36.449 0 01-3.843-7.934h46.944zm-90.021 0h24.658c1.204 4.333 2.915 8.338 5.165 12h-33.848l-1.747-.26c-1.397-.663-1.89-.708-2.901-1.947-.373-.456-.603-1.011-.904-1.517l-20.764-50.654 6.434-15.941 23.907 58.319zm-78.981-27.916l28.017-69.269 1.062-1.718c.524-.423.984-.94 1.572-1.269 2.405-1.345 5.633-.748 7.41 1.248.447.503.713 1.142 1.07 1.713l7.48 18.248-6.334 16.185-6.661-16.248-27.244 67.358-6.372-16.248zm-92.776 27.916h81.485l3.158-7.81 6.468 16.012-.019.048-.9 1.527c-1.145 1.049-1.373 1.491-2.908 1.961-.547.167-1.129.174-1.754.262h-88.051l2.521-12z"/>
                </mask>
              </defs>

              <path style="fill: var(--color-logo-background, #87d7af)"
                    d="M159.577 259.753l.782 3.75h-40.987l.855-3.75h39.35zm38.349.016c18.098.071 36.198.328 54.292-.018 2.465-.078 4.64-2.235 4.657-4.786v-49.61h3.75s.108 33.151-.001 49.726c-.082 4.393-3.83 8.333-8.421 8.421l-55.128.001.851-3.734zM253.594 41.25c3.651.002 6.654 3.005 6.656 6.656v146.438l-4.669.014-22-105.817-16.121 5.191-31.918 139.758-30.408-145.588-15.447 4.974-35.942 158.479-13.064-62.747H69.768l15.458 74.52H45.029c-3.651-.002-6.655-3.006-6.657-6.657V47.906c.002-3.651 3.006-6.654 6.657-6.656h208.565z"
                    fill-rule="nonzero"/>
              <path
                d="M166.345 292.218l-31.356-150.414-32.908 144.473-21.115 6.798-18.27-87.754H24.751l4.661-20.463h50.019l13.064 62.747 35.942-158.479 15.447-4.974 30.408 145.588L206.21 89.982l16.121-5.191 20.752 99.976h28.498l4.279 20.554-49.472.097-13.448-64.594-33.56 147.197-13.035 4.197zM474.467 112.587l-.458-.091-43.121 106.842h-11.993l-43.122-106.751-.458.092 1.832 53.924v38.635l14.923 2.564v11.536h-47.882v-11.536l14.923-2.564v-105.01l-14.923-2.564V86.037h37.994l42.48 108.307h.55l42.389-108.307h38.086v11.627l-14.924 2.564v105.01l14.924 2.564v11.536h-47.883v-11.536l14.924-2.564v-38.635l1.739-54.016zm38.368 56.305c0-14.648 3.967-26.718 11.902-36.209 7.934-9.491 18.707-14.237 32.318-14.237 13.672 0 24.49 4.731 32.455 14.191 7.966 9.461 11.948 21.546 11.948 36.255v2.014c0 14.771-3.967 26.856-11.902 36.255-7.934 9.4-18.707 14.099-32.318 14.099-13.733 0-24.567-4.715-32.501-14.145-7.935-9.43-11.902-21.499-11.902-36.209v-2.014zm18.036 2.014c0 10.498 2.212 19.165 6.637 26.001 4.426 6.836 11.002 10.254 19.73 10.254 8.545 0 15.045-3.418 19.501-10.254 4.455-6.836 6.683-15.503 6.683-26.001v-2.014c0-10.376-2.228-19.012-6.683-25.909-4.456-6.897-11.017-10.346-19.684-10.346-8.667 0-15.198 3.449-19.592 10.346-4.395 6.897-6.592 15.533-6.592 25.909v2.014zm79.162 36.896l14.923-2.564v-70.77l-14.923-2.563v-11.627h31.036l1.282 14.74c3.296-5.249 7.431-9.324 12.405-12.223 4.974-2.899 10.635-4.349 16.983-4.349 10.681 0 18.951 3.129 24.811 9.385 5.859 6.256 8.789 15.915 8.789 28.976v48.431l14.923 2.564v11.536H672.38v-11.536l14.923-2.564v-48.065c0-8.728-1.724-14.923-5.173-18.585-3.448-3.662-8.712-5.493-15.793-5.493-5.188 0-9.78 1.251-13.778 3.754-3.998 2.502-7.187 5.92-9.567 10.254v58.135l14.923 2.564v11.536h-47.882v-11.536zm121.071 0l14.923-2.564v-70.77l-14.923-2.563v-11.627h32.959v84.96l14.923 2.564v11.536h-47.882v-11.536zm32.959-112.885h-18.036V76.515h18.036v18.402zm51.15 1.465v23.896h18.768v13.366h-18.768v60.15c0 4.639.961 7.904 2.884 9.797 1.923 1.892 4.471 2.838 7.645 2.838 1.587 0 3.372-.138 5.356-.412 1.983-.275 3.646-.565 4.989-.87l2.472 12.36c-1.709 1.098-4.211 1.998-7.507 2.7a47.357 47.357 0 01-9.888 1.053c-7.324 0-13.153-2.212-17.486-6.637-4.334-4.425-6.501-11.368-6.501-20.829v-60.15h-15.655v-13.366h15.655V96.382h18.036zm28.681 72.51c0-14.648 3.967-26.718 11.902-36.209 7.935-9.491 18.707-14.237 32.318-14.237 13.672 0 24.49 4.731 32.456 14.191 7.965 9.461 11.947 21.546 11.947 36.255v2.014c0 14.771-3.967 26.856-11.902 36.255-7.934 9.4-18.707 14.099-32.318 14.099-13.733 0-24.566-4.715-32.501-14.145-7.935-9.43-11.902-21.499-11.902-36.209v-2.014zm18.036 2.014c0 10.498 2.213 19.165 6.638 26.001 4.425 6.836 11.001 10.254 19.729 10.254 8.545 0 15.045-3.418 19.501-10.254 4.455-6.836 6.683-15.503 6.683-26.001v-2.014c0-10.376-2.228-19.012-6.683-25.909-4.456-6.897-11.017-10.346-19.684-10.346-8.667 0-15.198 3.449-19.592 10.346-4.395 6.897-6.592 15.533-6.592 25.909v2.014zm79.567-39.001v-11.627h31.036l1.74 14.373c2.807-5.066 6.271-9.033 10.391-11.902 4.12-2.868 8.835-4.303 14.145-4.303 1.403 0 2.853.107 4.348.321 1.496.214 2.64.442 3.434.687l-2.381 16.754-10.254-.55c-4.76 0-8.758 1.114-11.993 3.342-3.235 2.228-5.737 5.356-7.507 9.384v56.854l14.923 2.564v11.536h-47.882v-11.536l14.923-2.564v-70.77l-14.923-2.563zm71.305 36.987c0-14.648 3.968-26.718 11.902-36.209 7.935-9.491 18.707-14.237 32.318-14.237 13.672 0 24.491 4.731 32.456 14.191 7.965 9.461 11.947 21.546 11.947 36.255v2.014c0 14.771-3.967 26.856-11.902 36.255-7.934 9.4-18.707 14.099-32.318 14.099-13.733 0-24.566-4.715-32.501-14.145-7.934-9.43-11.902-21.499-11.902-36.209v-2.014zm18.036 2.014c0 10.498 2.213 19.165 6.638 26.001 4.425 6.836 11.001 10.254 19.729 10.254 8.545 0 15.046-3.418 19.501-10.254 4.456-6.836 6.683-15.503 6.683-26.001v-2.014c0-10.376-2.227-19.012-6.683-25.909-4.455-6.897-11.017-10.346-19.684-10.346-8.667 0-15.197 3.449-19.592 10.346s-6.592 15.533-6.592 25.909v2.014zm79.567-39.001v-11.627h31.036l1.74 14.373c2.807-5.066 6.271-9.033 10.391-11.902 4.12-2.868 8.835-4.303 14.145-4.303 1.404 0 2.853.107 4.349.321 1.495.214 2.639.442 3.433.687l-2.381 16.754-10.253-.55c-4.761 0-8.759 1.114-11.994 3.342-3.235 2.228-5.737 5.356-7.507 9.384v56.854l14.923 2.564v11.536h-47.882v-11.536l14.923-2.564v-70.77l-14.923-2.563z"
                fill="#ffffff" fill-rule="nonzero"/>
              <path
                d="M316.036 194.973h88.325l40.825-100.935 41.377 100.935H606.23l19.293 43.776 20.292-43.776h221.487l43.443-100.185 43.556 100.185h115.916l20.292 44.526 20.646-44.526h81.809"
                fill="none" stroke="url(#monitoror-logo-line-gradient)" mask="url(#monitoror-logo-line-mask)"
                stroke-width="12"
                class="c-app--logo-line"/>
            </svg>
          </div>
          <monitoror-errors class="c-app--loading-errors"></monitoror-errors>
        </div>
      </div>
    </transition>
  </div>
</template>

<script lang="ts">
  import {Component, Vue} from 'vue-property-decorator'

  import CONFIG_VERIFY_ERRORS from '@/constants/configVerifyErrors'

  import MonitororErrors from '@/components/Errors.vue'
  import MonitororTile from '@/components/Tile.vue'
  import ConfigError from '@/interfaces/configError'
  import TileConfig from '@/interfaces/tileConfig'

  @Component({
    components: {
      MonitororErrors,
      MonitororTile,
    },
  })
  export default class App extends Vue {
    private static readonly SHOW_CURSOR_DELAY: number = 10 // 10 seconds

    /*
     * Data
     */

    private showCursor: boolean = true
    private showCursorTimeout!: number
    private taskRunnerInterval!: number

    /*
     * Computed
     */

    get classes() {
      return {
        'c-app__show-cursor': this.hasConfigVerifyErrors || this.showCursor,
        'c-app__no-scroll': !this.hasConfigVerifyErrors,
        'c-app__config-verify-errors': this.hasConfigVerifyErrors,
      }
    }

    get cssProperties() {
      const tilesCount = this.tiles.reduce((accumulator, tile) => {
        return accumulator + (tile.rowSpan || 1) * (tile.columnSpan || 1)
      }, 0)

      return {
        '--columns': this.columns,
        '--rows': Math.ceil(tilesCount / this.columns),
        '--zoom': this.zoom,
      }
    }

    get appLoadingClasses() {
      return {
        'c-app--loading__error': this.hasErrors,
        'c-app--loading__warning': !this.isOnline,
        'c-app--loading__config-verify-errors': this.hasConfigVerifyErrors,
      }
    }

    get columns(): number {
      return this.$store.state.columns
    }

    get zoom(): number {
      return this.$store.state.zoom
    }

    get tiles(): TileConfig[] {
      return this.$store.state.tiles
    }

    get errors(): ConfigError[] {
      return this.$store.state.errors
    }

    get hasErrors(): boolean {
      return this.errors.length > 0
    }

    get hasConfigVerifyErrors(): boolean {
      return this.errors.filter((error) => CONFIG_VERIFY_ERRORS.includes(error.id)).length > 0
    }

    get loadingProgress(): number {
      return this.$store.getters.loadingProgress
    }

    get loadingProgressBarStyle() {
      return {
        transform: `translateX(-${100 - this.loadingProgress * 100}%)`,
      }
    }

    get showLoading(): boolean {
      return this.loadingProgress < 1 || !this.isOnline || this.hasErrors
    }

    get isOnline(): boolean {
      return this.$store.state.online
    }

    /*
     * Methods
     */

    public resetShowCursorTimeout() {
      clearTimeout(this.showCursorTimeout)
      this.showCursor = true
      this.showCursorTimeout = setTimeout(() => {
        this.showCursor = false
      }, App.SHOW_CURSOR_DELAY * 1000)
    }

    public dispatchUpdateNetworkState() {
      return this.$store.dispatch('updateNetworkState')
    }

    /*
     * Hooks
     */

    private async mounted() {
      await Vue.nextTick()

      window.addEventListener('online', this.dispatchUpdateNetworkState)
      window.addEventListener('offline', this.dispatchUpdateNetworkState)

      this.taskRunnerInterval = setInterval(() => {
        this.$store.dispatch('runTasks')
      }, 50)

      await this.$store.dispatch('init')
    }

    private beforeDestroy() {
      window.removeEventListener('online', this.dispatchUpdateNetworkState)
      window.removeEventListener('offline', this.dispatchUpdateNetworkState)

      clearTimeout(this.showCursorTimeout)
      clearInterval(this.taskRunnerInterval)

      this.$store.dispatch('killAllTasks', [])
    }
  }
</script>

<style lang="scss">
  #app {
    height: 100%;
    width: 100%;

    --columns: 1;
    --rows: 1;
    --zoom: 1;

    &:not(.c-app__show-cursor) {
      cursor: none;
      user-select: none;
    }

    @media screen and (max-width: 750px) {
      --columns: 2 !important;
      --rows: 0 !important;
      --zoom: 0.65 !important;
    }

    @media screen and (max-width: 500px) {
      --columns: 1 !important;
    }

    &.c-app__no-scroll {
      @media screen and (min-width: 751px) {
        overflow: hidden;
      }
    }
  }

  .c-app--tiles-container {
    display: grid;
    grid-template-columns: repeat(var(--columns), 1fr);
    grid-gap: 6px;
    grid-auto-rows: calc((100vh - 6px) / var(--rows) - 6px);
    margin: 6px;
  }

  .c-app--loading {
    --color-logo-background: #87d7af;
    --color-logo-m: #ffffff;

    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    color: var(--color-spindle);
    background: var(--color-background);
    text-align: center;
    z-index: 50;
    will-change: opacity;
    transition: opacity 500ms;
  }

  .c-app--loading__warning {
    --color-logo-background: var(--color-warning);
    --color-logo-line-start: #3f4b51;
  }

  .c-app--loading__error {
    --color-logo-background: var(--color-failed);
    --color-logo-line-start: #3f4252;
  }

  .c-app--loading__config-verify-errors {
    position: absolute;
    text-align: left;
    bottom: initial;
    min-height: 100vh;
    background: var(--color-docs-background);
  }

  .c-app--loading__config-verify-errors .c-app--loading-container {
    position: relative;
    width: 1000px;
    max-width: 100%;
    padding: 100px 0;
    margin: 0 auto;

    hr {
      margin: 150px auto;
    }
  }

  .c-app--logo {
    position: absolute;
    top: 50vh;
    left: 50%;
    width: 80%;
    transform: translate(-50%, -50%);
    will-change: transform, opacity, width;
    animation: fadeIn 1s, logoSlideIn 1s;
    transition: transform 750ms, opacity 750ms, width 300ms, top 300ms;

    path {
      will-change: fill;
      transition: fill 300ms;
    }

    stop {
      will-change: stop-color;
      transition: stop-color 300ms;
    }
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @keyframes logoSlideIn {
    from {
      transform: translate(-50%, -40%);
    }
    to {
      transform: translate(-50%, -50%);
    }
  }

  .c-app--logo-line {
    stroke-dasharray: 0, 1265, 0;
    animation: appLogoLine 7s both 300ms infinite;
  }

  @keyframes appLogoLine {
    0% {
      stroke-dashoffset: 2535;
    }
    20%,
    80% {
      stroke-dashoffset: 1265;
    }
    100% {
      stroke-dashoffset: 0;
    }
  }

  .c-app--loading__config-verify-errors .c-app--logo {
    top: 120px;
    width: 700px;
    max-width: calc(100% - 100px);
  }

  .c-app--loading-errors {
    position: absolute;
    top: 70%;
    left: 50%;
    padding: 0 50px;
    width: 100%;
    transform: translateX(-50%);
    font-weight: 400;
    font-size: 24px;
    animation: fadeIn 300ms ease 500ms;
    animation-fill-mode: both;
  }

  .c-app--loading__config-verify-errors .c-app--loading-errors {
    position: relative;
    top: initial;
    left: initial;
    margin-top: 130px;
    transform: initial;
  }

  .loading-fade-enter,
  .loading-fade-leave-to {
    opacity: 0;
  }

  .loading-fade-leave-to .c-app--logo {
    transform: translate(-50%, -60%);
    opacity: 0;
  }

  .loading-fade-leave-to .c-app--loading-progress-bar {
    transform: translateX(0);
  }

  .c-app--loading__config-verify-errors .c-app--loading-progress {
    display: none;
  }

  .c-app--loading-progress {
    position: absolute;
    top: 0;
    right: 0;
    left: 0;
    display: block;
    height: 5px;
    margin: 0 auto;
    background: rgba(0, 0, 0, 0.2);
  }

  .c-app--loading-progress-bar {
    display: block;
    width: 100%;
    height: 100%;
    background: var(--color-logo-background);
    transition: transform 150ms ease-in-out;
  }
</style>
