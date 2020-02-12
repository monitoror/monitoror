import { createLocalVue, shallowMount } from '@vue/test-utils'
import Vuex from 'vuex'

import App from '@/App.vue'
import store from '@/store'

const localVue = createLocalVue()
localVue.use(Vuex)

describe('App.vue', () => {
  it('renders container', () => {
    const wrapper = shallowMount(App, { store, localVue })
    expect(wrapper.find('.c-app--tiles-container')).toBeTruthy()
  })
})
