import { shallowMount } from '@vue/test-utils'

import App from '@/App.vue'

describe('App.vue', () => {
  it('renders container', () => {
    const wrapper = shallowMount(App)
    expect(wrapper.find('.c-app--tiles-container')).toBeTruthy()
  })
})
