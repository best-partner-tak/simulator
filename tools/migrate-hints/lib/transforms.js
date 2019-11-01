module.exports.groupHintsByTask = groupHintsByTask = (hints) => {
  return hints.reduce((acc, hint) => {
    // Pull apart the hint text on the first colon
    const [task, hintText] = hint.text.split(': ', 2)

    const sortOrder = Number(task[task.length - 1])

    // Create a new group for the task if it doesn't exist
    if (!acc[task]) acc[task] = { "sort-order": sortOrder, hints: []}

    // Reconstruct a hint without the task prefix and add it to the group
    acc[task].hints.push({ text: hintText })

    return acc
  }, {})
}

// Takes an array of top level yaml properties represented as a JS object and
// Returns a transformed array of top-level properties in the new format
module.exports.transformV0ToV1 = hints => {
  const transformed = {}
  // remove hint count
  delete hints[0]['general_overview']["num-hints"]

  transformed['general_overview'] = hints[0]['general_overview']
  // add version
  transformed.kind = 'cp.simulator/scenario:1.0.0'

  // transform hints
  const hintsList  = Object.entries(hints[1].hints)
    .reduce((acc, [key, val], idx) => { 
      acc.push({ text: val })
      return acc
    }, [])

  transformed.tasks = groupHintsByTask(hintsList)

  return transformed
}
