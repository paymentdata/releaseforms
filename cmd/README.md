pkg cmd: 
commandline utility to generate a list of change items which represents changes from one past commit and one commit (either most recent or more recent)

1. Δ/delta

    - generating a list of merged PRs from last released commit to latest commit being released.

2. context

    - for each introduced Δ, generate a change item which collectively represent the given release changelog.


# Usage example:

`./delta $lastDeploymentSHA | ./contextaggregator`, where lastDeploymentSHA is set to the prior releases Commit.


# Full(~copy+paste) local example: 
_(while this is only on branch `cmd-paradigm` at least)_

**needs athenapdf instance running**: `docker run -p 8080:8080 --rm arachnysdocker/athenapdf-service`

```bash
git clone https://github.com/paymentdata/releaseforms
cd releaseforms/
git checkout remotes/origin/cmd-paradigm
go build github.com/paymentdata/releaseforms/cmd/contextaggregator
go build github.com/paymentdata/releaseforms/cmd/delta
echo -e "REPO=releaseforms\nORG=paymentdata" > .env
lastDeploymentSHA=d95539a
./delta $lastDeploymentSHA | ./contextaggregator
ls -lathr
```