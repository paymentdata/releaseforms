# objective:

### to implement a service which can subscribes to _some stream_ of **release events** and, upon release, publish some **release document** which is structured to comply with [PCIDSS](https://pcicompliance.stanford.edu/sites/g/files/sbiybj7706/f/16._change_control_policy_0.pdf)

# releaseforms

generating a document on every release can be cumbersome for an organization any way you look at it.

*especially* if the process is manual, and requires personnel.

-----

`releaseforms` is an effort to enable a pipeline of specially trained gophers who themselves can construct associated documentation, 

providing stakeholders and management with mostly filled out forms which only need to be reviewed and signed off.

## how to use: [w/ provided defaults and minimal setup]

#### Configure GitHub repo/org webhook:

1 navigate to _Settings > Webhooks_ on GitHub

2 select **add webhook**

3 **PayloadURL**: `yourhost:8081/webhooks1`, **Content-Type**: `application/json`, **Secret**: `MyGitHubSuperSecretSecrect...?`

4 `go get github.com/paymentdata/releaseforms` and run w/ the sample `.env` in the ${PWD}

5 `./releaseforms`

6 trigger webhook

7 iterate as needed
