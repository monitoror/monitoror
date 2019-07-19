:warning: 
* previousStatus can be UNKNOWN
* estimatedDuration can be 0
* author is optional

```json
{
  "category": "BUILD",
  "type": "JENKINS-BUILD||TRAVISCI-BUILD||...",
  "status": "SUCCESS||FAILED||WARNING||ABORTED",
  "previousStatus": "SUCCESS||FAILED||WARNING||UNKNOWN",
  "label": "Test-Job",
  "startedAt": "2019-07-17T13:41:04+02:00",
  "finishedAt": "2019-07-17T16:41:40+02:00",
  "author": {
    "name": "Captain Teemo",
    "avatarUrl": "https://www.gravatar.com/avatar/022a1cd3de7bcf2a3b7cb4253078ed65?d=blank"
  }
}
```

```json
{
  "category": "BUILD",
  "type": "JENKINS-BUILD||TRAVISCI-BUILD||...",
  "status": "QUEUE",
  "previousStatus": "SUCCESS||FAILED||WARNING||UNKNOWN",
  "label": "Test-Job",
  "startedAt": "2019-07-17T13:41:04+02:00",
  "author": {
    "name": "Captain Teemo",
    "avatarUrl": "https://www.gravatar.com/avatar/022a1cd3de7bcf2a3b7cb4253078ed65?d=blank"
  }
}
```

```json
{
  "category": "BUILD",
  "type": "JENKINS-BUILD||TRAVISCI-BUILD||...",
  "status": "RUNNING",
  "previousStatus": "SUCCESS||FAILED||WARNING||UNKNOWN",
  "label": "Test-Job",
  "startedAt": "2019-07-17T13:41:04+02:00",
  "duration": 6,
  "estimatedDuration": 100,
  "author": {
    "name": "Captain Teemo",
    "avatarUrl": "https://www.gravatar.com/avatar/022a1cd3de7bcf2a3b7cb4253078ed65?d=blank"
  }
}
```

```json
{
  "category": "BUILD",
  "type": "JENKINS-BUILD||TRAVISCI-BUILD||...",
  "status": "DISABLED",
  "label": "Test-Job"
}
```

```json
{
  "category": "BUILD",
  "type": "JENKINS-BUILD||TRAVISCI-BUILD||...",
  "status": "WARNING",
  "label": "Test-Job",
  "message": "error message ..."
}
```
