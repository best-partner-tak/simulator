const { loadYamlFile, getProgress, saveProgress } = require('./io')
const { createLogger } = require('./logger')
const { TASKS_FILE_PATH, PROGRESS_FILE_PATH } = require('./constants')

const logger = createLogger({})

// All the optional parameters in these functions are to enable injection
// for testing

function showHint (task, index, taskspath = TASKS_FILE_PATH, log = logger) {
  const { tasks } = loadYamlFile(taskspath)

  if (!tasks[task]) {
    return log.warn('Cannot find task')
  }

  if (index >= tasks[task].hints.length) {
    return log.warn('There are no more hints for this task')
  }

  const hint = tasks[task].hints[index]

  log.info(hint.text)
}

function showHints (task, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  const progress = getProgress(progresspath)
  const { tasks } = loadYamlFile(taskspath)

  if (!tasks[task]) {
    return logger.warn('Cannot find task')
  }

  if (progress[task] === undefined) {
    return logger.info('You have not seen any hints for this task')
  }

  const lastSeenHintIndex = progress[task]
  const hintcount = tasks[task].hints.length

  for (let i = 0; i <= lastSeenHintIndex && i < hintcount; i++) {
    logger.info(tasks[task].hints[i].text)
  }

  if (lastSeenHintIndex >= hintcount - 1) {
    return logger.info('You have seen all the hints for this task')
  }
}

function nextHint (task, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  if (!task) {
    return logger.error('No task provided to nextHint')
  }

  const progress = getProgress(progresspath)
  let hintIndex

  if (progress[task] === undefined) {
    hintIndex = progress[task] = 0
  } else {
    const { tasks } = loadYamlFile(taskspath)

    hintIndex = progress[task] = ++progress[task]
    if (hintIndex >= tasks[task].hints.length) {
      return logger.info('You have seen all the hints for this task')
    }
  }
  saveProgress(progress, progresspath)

  showHint(task, hintIndex, taskspath, log)
}

module.exports = {
  showHint,
  showHints,
  nextHint
}
