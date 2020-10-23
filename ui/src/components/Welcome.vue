<template>
  <div class="c-monitoror-welcome" :class="classes">
    <template v-if="shouldShowWelcomeNextStep">
      <h1 class="c-monitoror-welcome--title">Welcome</h1>
      <h2 class="c-monitoror-welcome--subtitle">You're almost there!</h2>

      <div class="c-monitoror-welcome--content">
        <div class="c-monitoror-welcome--next-step-content">
          <div class="c-monitoror-welcome--next-step-sup-title">
            Step {{currentStep}} of 3
          </div>

          <!-- STEP 1 -->
          <template v-if="currentStep === 1">
            <h3 class="c-monitoror-welcome--next-step-title">
              Create your UI configuration
            </h3>
            <p>
              As a starting point, save this in a <code>config.json</code> file next to <code>monitoror</code> binary:
            </p>
            <pre><code>{
  "version": "2.0",
  "columns": 2,
  "tiles": [
    { "type": "PORT", "label": "Welcome config example", "params": { "hostname": "127.0.0.1", "port": {{appPort}} } },
    { "type": "HTTP-RAW", "label": "Monitoror stars", "params": { "url": "https://github.com/monitoror/monitoror", "regex": "(\\d+) users starred" } }
  ]
}</code></pre>
            <p class="note">
              <strong>Hey, what's that "UI configuration"?</strong> <br>
              It's a JSON file where you will define what you want to display.
            </p>
          </template>

          <!-- STEP 2 -->
          <template v-if="currentStep === 2">
            <h3 class="c-monitoror-welcome--next-step-title">
              Setup Core configuration
            </h3>
            <p>
              You need to tell to the Core that you have a new UI configuration.
              To do so, put the following in a file named <code>.env</code> next to <code>monitoror</code> binary file:
            </p>
            <pre><code>MO_CONFIG="./config.json"</code></pre>
            <p class="note">
              <strong>Another configuration?</strong> <br>
              Yes! This one is mostly here to register UI configurations and to setup <em>monitorables</em>,
              which are things that you can use in your UI configurations.
            </p>
          </template>

          <!-- STEP 3 -->
          <template v-if="currentStep === 3">
            <h3 class="c-monitoror-welcome--next-step-title">
              Restart the Core
            </h3>
            <p>
              As you have updated the <code>.env</code> file, you need to stop and start again the
              <code>monitoror</code> command, and you should see that in start logs:
            </p>
            <pre><code>AVAILABLE NAMED CONFIGURATIONS

  default -> ./config.json</code></pre>
            <p class="note">
              <strong>Wait, when did I need to restart?</strong> <br>
              You need to restart the Core (the <code>monitoror</code> binary) each time you update the
              <code>.env</code> file.
            </p>
          </template>

          <!-- Step navigation -->
          <div class="c-monitoror-welcome--nav">
            <button v-if="currentStep > 1" @click="currentStep -= 1"
                    class="c-monitoror-welcome--button c-monitoror-welcome--button__back">
              <svg viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                <path
                  d="M31.2 14L18.7 1.6c-1.1-1.1-3-1.1-4.1 0-1.1 1.1-1.1 2.9 0 4.1L25 16 14.6 26.3c-1.1 1.1-1.1 2.9 0 4.1 1.1 1.1 3 1.1 4.1 0L31.2 18c.6-.6.9-1.3.8-2 0-.7-.3-1.5-.8-2z"
                  fill="currentColor"></path>
                <path
                  d="M2.3 13h22.9c1.3 0 2.3 1 2.3 2.3v1.5c0 1.3-1 2.3-2.3 2.3H2.3c-1.3 0-2.3-1-2.3-2.3v-1.5C0 14.1 1 13 2.3 13z"
                  fill="currentColor"></path>
              </svg>
              Back
            </button>
            <button v-if="currentStep < 3" @click="currentStep += 1"
                    class="c-monitoror-welcome--button c-monitoror-welcome--button__plain">
              Next step
              <svg viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                <path
                  d="M31.2 14L18.7 1.6c-1.1-1.1-3-1.1-4.1 0-1.1 1.1-1.1 2.9 0 4.1L25 16 14.6 26.3c-1.1 1.1-1.1 2.9 0 4.1 1.1 1.1 3 1.1 4.1 0L31.2 18c.6-.6.9-1.3.8-2 0-.7-.3-1.5-.8-2z"
                  fill="currentColor"></path>
                <path
                  d="M2.3 13h22.9c1.3 0 2.3 1 2.3 2.3v1.5c0 1.3-1 2.3-2.3 2.3H2.3c-1.3 0-2.3-1-2.3-2.3v-1.5C0 14.1 1 13 2.3 13z"
                  fill="currentColor"></path>
              </svg>
            </button>
            <button disabled class="c-monitoror-welcome--button" v-if="currentStep === 3">
              You're done! &#x1F389; <!-- Party Popper Emoji 🎉 -->
            </button>
          </div>
        </div>
      </div>
    </template>
    <template v-else>
      <h1 class="c-monitoror-welcome--title">Welcome <span class="hide-on-mobile">back</span></h1>
      <h2 class="c-monitoror-welcome--subtitle">Choose a configuration</h2>

      <div class="c-monitoror-welcome--content c-monitoror-welcome--content__choose-configuration">
        <template v-for="configMetadata in configList">
          <a :href="configMetadata.uiUrl" class="c-monitoror-welcome--button c-monitoror-welcome--button__card">
            <code>{{configMetadata.name}}</code>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentcolor">
              <path fill-rule="evenodd"
                    d="M11.31 2.525a9.648 9.648 0 011.38 0c.055.004.135.05.162.16l.351 1.45c.153.628.626 1.08 1.173 1.278.205.074.405.157.6.249a1.832 1.832 0 001.733-.074l1.275-.776c.097-.06.186-.036.228 0 .348.302.674.628.976.976.036.042.06.13 0 .228l-.776 1.274a1.832 1.832 0 00-.074 1.734c.092.195.175.395.248.6.198.547.652 1.02 1.278 1.172l1.45.353c.111.026.157.106.161.161a9.653 9.653 0 010 1.38c-.004.055-.05.135-.16.162l-1.45.351a1.833 1.833 0 00-1.278 1.173 6.926 6.926 0 01-.25.6 1.832 1.832 0 00.075 1.733l.776 1.275c.06.097.036.186 0 .228a9.555 9.555 0 01-.976.976c-.042.036-.13.06-.228 0l-1.275-.776a1.832 1.832 0 00-1.733-.074 6.926 6.926 0 01-.6.248 1.833 1.833 0 00-1.172 1.278l-.353 1.45c-.026.111-.106.157-.161.161a9.653 9.653 0 01-1.38 0c-.055-.004-.135-.05-.162-.16l-.351-1.45a1.833 1.833 0 00-1.173-1.278 6.928 6.928 0 01-.6-.25 1.832 1.832 0 00-1.734.075l-1.274.776c-.097.06-.186.036-.228 0a9.56 9.56 0 01-.976-.976c-.036-.042-.06-.13 0-.228l.776-1.275a1.832 1.832 0 00.074-1.733 6.948 6.948 0 01-.249-.6 1.833 1.833 0 00-1.277-1.172l-1.45-.353c-.111-.026-.157-.106-.161-.161a9.648 9.648 0 010-1.38c.004-.055.05-.135.16-.162l1.45-.351a1.833 1.833 0 001.278-1.173 6.95 6.95 0 01.249-.6 1.832 1.832 0 00-.074-1.734l-.776-1.274c-.06-.097-.036-.186 0-.228.302-.348.628-.674.976-.976.042-.036.13-.06.228 0l1.274.776a1.832 1.832 0 001.734.074 6.95 6.95 0 01.6-.249 1.833 1.833 0 001.172-1.277l.353-1.45c.026-.111.106-.157.161-.161zM12 1c-.268 0-.534.01-.797.028-.763.055-1.345.617-1.512 1.304l-.352 1.45c-.02.078-.09.172-.225.22a8.45 8.45 0 00-.728.303c-.13.06-.246.044-.315.002l-1.274-.776c-.604-.368-1.412-.354-1.99.147-.403.348-.78.726-1.129 1.128-.5.579-.515 1.387-.147 1.99l.776 1.275c.042.069.059.185-.002.315a8.45 8.45 0 00-.302.728c-.05.135-.143.206-.221.225l-1.45.352c-.687.167-1.249.749-1.304 1.512a11.149 11.149 0 000 1.594c.055.763.617 1.345 1.304 1.512l1.45.352c.078.02.172.09.22.225.09.248.191.491.303.729.06.129.044.245.002.314l-.776 1.274c-.368.604-.354 1.412.147 1.99.348.403.726.78 1.128 1.129.579.5 1.387.515 1.99.147l1.275-.776c.069-.042.185-.059.315.002.237.112.48.213.728.302.135.05.206.143.225.221l.352 1.45c.167.687.749 1.249 1.512 1.303a11.125 11.125 0 001.594 0c.763-.054 1.345-.616 1.512-1.303l.352-1.45c.02-.078.09-.172.225-.22.248-.09.491-.191.729-.303.129-.06.245-.044.314-.002l1.274.776c.604.368 1.412.354 1.99-.147.403-.348.78-.726 1.129-1.128.5-.579.515-1.387.147-1.99l-.776-1.275c-.042-.069-.059-.185.002-.315.112-.237.213-.48.302-.728.05-.135.143-.206.221-.225l1.45-.352c.687-.167 1.249-.749 1.303-1.512a11.125 11.125 0 000-1.594c-.054-.763-.616-1.345-1.303-1.512l-1.45-.352c-.078-.02-.172-.09-.22-.225a8.469 8.469 0 00-.303-.728c-.06-.13-.044-.246-.002-.315l.776-1.274c.368-.604.354-1.412-.147-1.99-.348-.403-.726-.78-1.128-1.129-.579-.5-1.387-.515-1.99-.147l-1.275.776c-.069.042-.185.059-.315-.002a8.465 8.465 0 00-.728-.302c-.135-.05-.206-.143-.225-.221l-.352-1.45c-.167-.687-.749-1.249-1.512-1.304A11.149 11.149 0 0012 1zm2.5 11a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0zm1.5 0a4 4 0 11-8 0 4 4 0 018 0z"/>
            </svg>
          </a>
        </template>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
  import {Options, Vue} from 'vue-class-component'

  import Route from '@/enums/route'
  import ConfigMetadata from '@/interfaces/configMetadata'

  @Options({})
  export default class MonitororWelcome extends Vue {
    /*
     * Data
     */
    private currentStep: number = 1
    private shouldHaveScroll: boolean = false

    /*
     * Computed
     */

    get classes() {
      return {
        'c-monitoror-welcome__no-scroll': !this.shouldHaveScroll,
      }
    }

    get shouldShowWelcomeNextStep(): boolean {
      const isOnWelcomeNextStepRoute = Route.Welcome === this.$store.getters.currentRoute

      return isOnWelcomeNextStepRoute || this.$store.getters.isNewUser
    }

    get configList(): ConfigMetadata[] {
      return this.$store.state.configList
    }

    get appPort(): string {
      const protocolPort = window.location.protocol === 'https:' ? '443' : '80'

      return window.location.port || protocolPort
    }

    /*
     * Methods
     */

    mounted() {
      this.$store.dispatch('fetchConfigList')

      setTimeout(() => {
        this.shouldHaveScroll = true
      }, 1650)
    }
  }
</script>

<style lang="scss">
  .c-monitoror-welcome {
    position: relative;
    padding-top: 175px;
    font-weight: 400;
    font-size: 24px;
    text-align: center;
    animation: welcomeSlideIn 525ms ease 1.15s both;

    &::before {
      content: "";
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 330px;
      background: var(--color-dark-background);
      transform: skewY(-7deg) scaleY(2);
      transform-origin: bottom center;
      z-index: -1;
    }

    &__no-scroll .c-monitoror-welcome--next-step-content {
      display: none;
    }

    ::-webkit-scrollbar {
      width: 10px;
      height: 10px;
    }

    ::-webkit-scrollbar-track {
      background: var(--color-code-background);
    }

    ::-webkit-scrollbar-thumb {
      background: var(--color-cello);
      border: 2px solid var(--color-code-background);
      border-top: 0;
      border-bottom: 0;
    }

    ::-webkit-scrollbar-thumb:hover {
      background: var(--color-action-required);
    }
  }

  .c-monitoror-welcome--title {
    display: inline-block;
    margin: 0 auto;
    font-weight: 700;
    font-size: 64px;
    color: transparent;
    background: linear-gradient(to right, var(--color-succeeded), var(--color-action-required));
    -webkit-background-clip: text;
    background-size: 200% 100%;
    animation: welcomeTitleGradient 10s ease 1.7s infinite;
  }

  .c-monitoror-welcome--subtitle {
    font-weight: 300;
    font-size: 28px;
    margin: 0;
  }

  .c-monitoror-welcome--content {
    width: 750px;
    max-width: calc(100% - 50px);
    margin: 0 auto;
    padding: 90px 0 0;
    font-size: 20px;
    line-height: 1.3;

    &__choose-configuration {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(230px, 1fr));
      grid-gap: 20px;
      padding-top: 130px;
      padding-bottom: 30px;
    }
  }

  .c-monitoror-welcome--next-step-sup-title {
    text-transform: uppercase;
    font-size: 18px;
    line-height: 1.2;
    letter-spacing: 1.1px;
    color: var(--color-succeeded);
    text-indent: 1px;
    text-align: center;
    font-family: "JetBrains Mono", monospace;
  }

  .c-monitoror-welcome--next-step-title {
    font-size: 36px;
    line-height: 1.3;
    font-weight: 300;
    margin: 0 0 20px;
    color: #ffffff;
    text-align: center;
  }

  .c-monitoror-welcome--next-step-content {
    position: relative;
    text-align: left;
    animation: stepSlideIn 700ms both;

    p {
      margin: 0 0 20px;

      &.note {
        border: 2px solid var(--color-cello);
        border-radius: 4px;
        padding: 20px 25px;
        margin-bottom: 20px;

        strong {
          color: var(--color-action-required);
        }
      }

      code {
        padding: 2px 5px;
        background: var(--color-background);
        border-radius: 4px;
      }
    }

    pre {
      display: block;
      margin: 20px auto;
      background: var(--color-background);
      border: 0;
      border-radius: 4px;
      text-align: left;

      code {
        display: block;
        padding: 15px 20px;
        color: var(--color-spring-wood);
        font-size: 16px;
        font-family: 'JetBrains Mono', monospace;
        overflow: auto;
      }
    }

    .c-monitoror-welcome--nav {
      display: flex;
      justify-content: space-between;
    }

    .c-monitoror-welcome--button {
      display: inline-block;
      cursor: pointer;
      padding: 15px 25px;
      margin-left: auto;
      margin-bottom: 20px;

      &__back {
        margin-left: initial;

        svg {
          transform: rotate(180deg);
          margin-left: 0;
          margin-right: 10px;
        }
      }
    }
  }

  .c-monitoror-welcome--button {
    display: inline-block;
    text-align: center;
    font-size: 16px;
    font-weight: 600;
    line-height: 1.2;
    letter-spacing: 0.3px;
    text-decoration: none;
    padding: 20px 30px;
    border-radius: 4px;
    border: 1px solid hsla(211, 52%, 88%, 0.15);
    color: var(--color-spindle);
    background: rgba(207, 223, 240, 0.03);
    transition: color 200ms, background 200ms, border-color 200ms;

    &::after {
      content: initial;
    }

    &:not(:disabled):hover {
      color: var(--color-background);
      background: var(--color-spindle);
      border-color: var(--color-spindle);
    }

    &:disabled {
      cursor: default;
    }

    &:focus {
      outline: none;
    }

    svg {
      display: inline-block;
      width: 13px;
      margin-left: 10px;
      transform: translateY(1px);
    }

    &__plain {
      color: var(--color-succeeded-dark);
      background: var(--color-succeeded);
      box-shadow: rgba(136, 216, 176, 0.35) 0 5px 50px, var(--color-succeeded-dark) 0 3px 5px;
    }

    &__card {
      width: 100%;
      height: 150px;
      text-align: left;
      font-size: 24px;
      font-weight: normal;
      padding: 22px 25px;
      overflow: hidden;

      svg {
        position: absolute;
        right: -20px;
        bottom: -20px;
        width: 100px;
        opacity: 0.1;
      }
    }
  }

  @media (max-width: 500px) {
    .hide-on-mobile {
      display: none !important;
    }
  }

  @keyframes welcomeTitleGradient {
    0%,
    100% {
      background-position: 0 50%;
    }
    50% {
      background-position: 100% 50%;
    }
  }

  @keyframes welcomeSlideIn {
    from {
      opacity: 0;
      transform: translateY(10%);
    }
  }

  @keyframes stepSlideIn {
    from {
      opacity: 0;
      transform: translateY(2%);
    }
  }
</style>
