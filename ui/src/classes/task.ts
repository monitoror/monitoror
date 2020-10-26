import {noop, now} from 'lodash-es'

import TaskType from '@/enums/taskType'
import TaskOptions from '@/types/taskOptions'

export default class Task {
  public readonly id: string
  public readonly type: TaskType
  public readonly interval: number
  public readonly retryOnFailInterval: number
  private running: boolean = false
  private done: boolean = false
  private dead: boolean = false
  private internalTime: number
  private failedAttemptsCount: number = 0
  private readonly executor: () => Promise<void>
  private readonly onFailedCallback: (failedAttemptsCount: number) => void

  constructor(
    {
      id,
      type,
      executor,
      interval = 0,
      initialDelay = 0,
      retryOnFailInterval = 0,
      onFailedCallback = noop,
    }: TaskOptions,
  ) {
    this.id = id
    this.type = type
    this.interval = interval
    this.retryOnFailInterval = retryOnFailInterval
    this.internalTime = now() + initialDelay
    this.executor = executor
    this.onFailedCallback = onFailedCallback
  }

  get time() {
    return this.internalTime
  }

  public isRunning() {
    return this.running
  }

  public isDone() {
    return this.done
  }

  public isDead() {
    return this.dead
  }

  public kill() {
    this.dead = true
  }

  public async run() {
    if (this.isDead()) {
      return
    }

    this.running = true
    try {
      await this.executor()
    } catch (e) {
      this.failedAttemptsCount += 1
      this.onFailedCallback(this.failedAttemptsCount)
    }
    this.done = true
    this.running = false

    if (this.interval === 0) {
      this.dead = true
    }
  }

  public prepareNextRun() {
    if (this.isDead() || this.isRunning()) {
      return
    }

    let interval = this.interval
    if (this.retryOnFailInterval > 0) {
      interval = this.retryOnFailInterval
    }

    this.done = false
    this.internalTime = now() + interval
  }
}
