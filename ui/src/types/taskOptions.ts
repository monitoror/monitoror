import TaskType from '@/enums/taskType'

type TaskOptions = {
  id: string,
  type: TaskType,
  executor: () => Promise<void>,
  interval?: number,
  initialDelay?: number,
  retryOnFailInterval?: number,
  onFailedCallback?: (failedAttemptsCount: number) => void,
}

export default TaskOptions
