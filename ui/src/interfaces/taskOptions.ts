import TaskType from '@/enums/taskType'

export default interface TaskOptions {
  id: string,
  type: TaskType,
  executor: () => Promise<void>,
  interval?: number,
  initialDelay?: number,
  retryOnFailInterval?: number,
  onFailedCallback?: (failedAttemptsCount: number) => void,
}
