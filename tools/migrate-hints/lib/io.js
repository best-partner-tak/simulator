const yaml = require('js-yaml')
const {truncateSync, writeFileSync, readFileSync, readdirSync} = require('fs')
const {resolve,join} = require('path')
const {transformV0ToV1} = require('../lib/transforms')

// Loads and parses a `hints.yaml` from the supplied absolute path. 
// Returns an object representing the yaml file
module.exports.loadHintsFile = p => {
  console.log(`Loading ${p}`)
  const doc = yaml.safeLoad(readFileSync(p, 'utf8'));
  return doc
}

// Serializes the supplied `hints` to YAML and overwrites an existing 
// `hints.yaml` file to the supplied path `p` with the YAMl
module.exports.writeHintsFile = (hints, p) => {
  console.log(`Writing transformed file ${p}`)
  const contents = yaml.safeDump(hints)
  truncateSync(p)
  writeFileSync(p, contents)
}

// Given a relative path to a scenario directory, scans for scenarios
// Returns a list of absolute paths to all `hints.yaml` files
module.exports.findScenarioHintsFiles = scenariosDir => {
  const absPath = resolve(scenariosDir)

  return readdirSync(absPath, { withFileTypes: true })
    .filter(dirent => dirent.isDirectory())
    .map(dirent => join(absPath, dirent.name, 'hints.yaml'))
}

