<template>
  <div class="full-width center-content">
    <img src="static/img/logo.png">
    <WelcomeMessage name="World" />
  </div>
</template>

<script>
import WelcomeMessage from 'components/Home/WelcomeMessage'

export default {
  components: {
    WelcomeMessage
  }
}
</script>
