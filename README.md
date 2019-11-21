# objective:

to implement a service which subscribes to **release events** and publishes associated **release documents** which are structured to comply with [PCIDSS](https://pcicompliance.stanford.edu/sites/g/files/sbiybj7706/f/16._change_control_policy_0.pdf)

# releaseforms

_problem statement:_ 

generating a document on every release can be cumbersome for an organization any way you look at it.

*especially* if the process is manual, and requires personnel.

-----

`releaseforms` **shall** implement a pipeline of specially trained gophers who themselves can construct associated documentation, 

providing stakeholders and management with mostly filled out forms which only need to be reviewed and signed off.

# current state of affairs/disclaimer things:

_be warned, this is NOT generating compliant forms right now_

_be warned, this MAY NEVER generate compliant forms_

*this project is very much in it's infancy, and therefore carries a risk of never delivering on the objective.*


## currently depends on:

- [athenapds](https://github.com/arachnys/athenapdf)
 (kudos to the folks working on arachnys projects!! 👏👏👏)

-----

now the fun stuff...

## how to use: [w/ provided defaults and minimal setup]

#### Configure GitHub repo/org webhook:

1 navigate to _Settings > Webhooks_ on GitHub

2 select **add webhook**

3 **PayloadURL**: `yourhost:8081/webhooks1`, **Content-Type**: `application/json`, **Secret**: `MyGitHubSuperSecretSecrect...?`

4 `go get github.com/paymentdata/releaseforms` and run w/ the sample `.env` in the ${PWD}

5 `./releaseforms`

6 trigger webhook

7 iterate as needed
