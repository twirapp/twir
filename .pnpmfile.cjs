function readPackage(pkg, context) {
  if (pkg.name === '@twurple/chat' && pkg.version.startsWith('5.1.6')) {
    pkg.dependencies = {
      ...pkg.dependencies,
      ircv3: '^0.29.4'
    }
    context.log('ircv3 locked in dependencies of twurple/chat')
  }


  return pkg
}

module.exports = {
  hooks: {
    readPackage
  }
}