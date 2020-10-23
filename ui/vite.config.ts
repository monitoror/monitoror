import {UserConfig} from 'vite'
// @ts-ignore
import tsResolver from 'vite-tsconfig-paths'

const config: UserConfig = {
  optimizeDeps: {
    include: [
      'ts-md5/dist/md5',
    ],
    exclude: [
      'ts-md5',
    ],
  },
  port: 8000,
  proxy: {
    '/api': 'http://localhost:8080',
  },
  resolvers: [
    tsResolver,
  ],
}

export default config
