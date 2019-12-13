`releaseforms` **shall** implement a pipeline of specially trained gophers who themselves can construct associated documentation, providing stakeholders and management with release documentation that satisfied control-objective 5 requirements.

## current state of affairs/disclaimer things:

*this project is very much in it's infancy, and therefore carries a risk of never delivering on the objective.*

-----

now the fun stuff...

commandline utility to generate a list of change items which represents changes from one past commit and one commit (either most recent or more recent)

1. Δ/delta

    - generating a list of merged PRs from last released commit to latest commit being released, depends on local/PWD being the git repo for the related release.

2. context

    - for each introduced Δ, generate a change item which collectively represent the given release changelog.


# Usage example:

`./delta $lastDeploymentSHA | ./contextaggregator`, where lastDeploymentSHA is set to the prior releases Commit.

## how to use: [w/ provided defaults and minimal setup]

### run this once for setup:
``` bash
git clone https://github.com/paymentdata/releaseforms
cd releaseforms/
echo -e "REPO=releaseforms\nORG=paymentdata" > .env
go test ./...
```

### for each release: (where the previously released commit for the .git is set as $lastDeploymentSHA)
``` bash
lastDeploymentSHA=d95539a
./delta $lastDeploymentSHA | ./contextaggregator
```

### PDF Generation currently depends on:

- [athenapds](https://github.com/arachnys/athenapdf)
 (kudos to the folks working on arachnys projects!! 👏👏👏)

run via
**needs athenapdf instance running**: `docker run -p 8080:8080 --rm arachnysdocker/athenapdf-service`

