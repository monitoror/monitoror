export enum TileType {
  HttpStatus = 'HTTP-STATUS',
  HttpRaw = 'HTTP-RAW',
  HttpFormatted = 'HTTP-FORMATTED',
  Ping = 'PING',
  Port = 'PORT',
  PingdomCheck = 'PINGDOM-CHECK',
  PingdomTransactionCheck = 'PINGDOM-TRANSACTION-CHECK',
  GitHubChecks = 'GITHUB-CHECKS',
  GitHubPullRequest = 'GITHUB-PULLREQUEST',
  GitHubCount = 'GITHUB-COUNT',
  GitLabPipeline = 'GITLAB-PIPELINE',
  GitLabMergeRequest = 'GITLAB-MERGEREQUEST',
  GitLabIssues = 'GITLAB-COUNT-ISSUES',
  TravisCiBuild = 'TRAVISCI-BUILD',
  JenkinsBuild = 'JENKINS-BUILD',
  AzureDevOpsBuild = 'AZUREDEVOPS-BUILD',
  AzureDevOpsRelease = 'AZUREDEVOPS-RELEASE',

  Empty = 'EMPTY',
  Group = 'GROUP',
}

export default TileType
