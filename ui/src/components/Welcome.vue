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
              <svg viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
                <!-- Bind href value to avoid Vite transforms -->
                <use :xlink:href="'./icons.svg#arrow'"/>
              </svg>
              Back
            </button>
            <button v-if="currentStep < 3" @click="currentStep += 1"
                    class="c-monitoror-welcome--button c-monitoror-welcome--button__plain">
              Next step
              <svg viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
                <!-- Bind href value to avoid Vite transforms -->
                <use :xlink:href="'./icons.svg#arrow'"/>
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
            <svg viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
              <!-- Bind href value to avoid Vite transforms -->
              <use :xlink:href="'./icons.svg#cog'"/>
            </svg>
          </a>
        </template>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, ref, onMounted} from 'vue'
import {useStore} from 'vuex'

import Route from '@/enums/route'
import ConfigMetadata from '@/types/configMetadata'

export default defineComponent({
  name: 'MonitororWelcome',
  setup() {
    const store = useStore()
    const currentStep = ref(1)
    const shouldHaveScroll = ref(false)

    const classes = computed((): Record<string, string | boolean> => {
      return {
        'c-monitoror-welcome__no-scroll': !shouldHaveScroll.value,
      }
    })

    const appPort = computed((): string => {
      const protocolPort = window.location.protocol === 'https:' ? '443' : '80'

      return window.location.port || protocolPort
    })

    const configList = computed((): ConfigMetadata[] => {
      return store.state.configList
    })

    const shouldShowWelcomeNextStep = computed((): boolean => {
      const isOnWelcomeNextStepRoute = Route.Welcome === store.getters.currentRoute

      return isOnWelcomeNextStepRoute || store.getters.isNewUser
    })

    onMounted(() => {
      store.dispatch('fetchConfigList')

      setTimeout(() => {
        shouldHaveScroll.value = true
      }, 1650)
    })

    return {
      // data
      currentStep,
      shouldHaveScroll,

      // attributes
      classes,

      // computed
      appPort,
      configList,
      shouldShowWelcomeNextStep,
    }
  }
})
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
