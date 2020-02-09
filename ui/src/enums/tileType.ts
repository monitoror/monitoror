export enum TileType {
  HttpAny = 'HTTP-ANY',
  HttpRaw = 'HTTP-RAW',
  HttpFormatted = 'HTTP-FORMATTED',
  Ping = 'PING',
  Port = 'PORT',
  Pingdom = 'PINGDOM-CHECK',
  GitHubChecks = 'GITHUB-CHECKS',
  GitHubIssues = 'GITHUB-ISSUES',
  GitLab = 'GITLAB-BUILD',
  Travis = 'TRAVISCI-BUILD',
  Jenkins = 'JENKINS-BUILD',
  AzureDevOpsBuild = 'AZUREDEVOPS-BUILD',
  AzureDevOpsRelease = 'AZUREDEVOPS-RELEASE',

  Empty = 'EMPTY',
  Group = 'GROUP',
}

export default TileType
