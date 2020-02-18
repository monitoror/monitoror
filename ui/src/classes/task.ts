import {now} from 'lodash-es'

export default class Task {
  public readonly id: string
  public readonly interval: number
  private done: boolean = false
  private dead: boolean = false
  private runTime: number
  private readonly task: () => void

  constructor(id: string, task: () => void, interval: number = 0, initialDelay: number = 0) {
    this.id = id
    this.interval = interval
    this.runTime = now() + initialDelay
    this.task = task
  }

  get time() {
    return this.runTime
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

    await this.task()
    this.done = true

    if (this.interval === 0) {
      this.dead = true
    }
  }

  public prepareNextRun() {
    if (this.isDead()) {
      return
    }

    this.done = false
    this.runTime = now() + this.interval
  }
}
