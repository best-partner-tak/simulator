<!--

NOTICE: THIS FILE IS AUTOGENERATED FROM docs/README.template.md

This file is evaled by a quickly cobbled together bash script to replace the variables.

Backticks are imterpreted by bash so use <code> for inline code and <pre> for code blocks.

If you need to include bsah code snippets you will need to change how the templating works.

-->

# Simulator

A distributed systems and infrastructure simulator for attacking and debugging
Kubernetes: <code>simulator</code> creates a kuberntes cluster for you in your AWS
account; runs scenarios which misconfigure it and/or leave it vulnerable to
compromise and trains you in mitigating against these vulnerabilities.

## Before you start

## Usage

The quickest way to get up and running is to simply:

<pre>
source <(curl https://raw.githubusercontent.com/controlplaneio/simulator-standalone/master/kubesim)
</pre>

or (inside your local version of this repository)

<pre>
make run
</pre>

This will drop you into a bash shell in a launch container.  You will have a
program on the <code>PATH</code> named <code>simulator</code> to interact with.

### The <code>kubesim</code> script

<code>kubesim</code> is a small script written in BASH for getting users up and running with simulator as fast as possible. It pulls the latest version of the simulator container and sets up some options for running the image. It can be installed with the following steps:

<pre>
# cURL the script from GitHub
curl -o kubesim https://raw.githubusercontent.com/controlplaneio/simulator-standalone/master/kubesim
# Make it executeable
chmod a+x kubesim
# Place the script on your path
cp kubesim /usr/local/bin
</pre>

Feel free to modify and/or extend the script if you wish.

### [Simulator CLI Usage](./docs/cli.md)

### End-to-end usage

<pre>
# Create a remote state bucket for terraform
simulator init
# Create the infra if it isn't there
simulator infra create
# List available scenarios
simulator scenario list
# Launch a scenario (sets up your cluster)
simulator scenario launch node-shock-tactics
# Attack the cluster
simulator ssh attack
# ... Complete the scenario :)
# Destroy your cluster when you are done
simulator infra destroy
</pre>

Current list of scenarios:

<pre>
1.5682037678385196e+09	INFO	cmd/scenario.go:25	Available scenarios:
1.568203767838556e+09	INFO	cmd/scenario.go:27	ID: container-ambush, Name: Container Ambush
1.5682037678385653e+09	INFO	cmd/scenario.go:27	ID: rbac-flanking-maneuver, Name: RBAC Flanking Maneuver
1.568203767838575e+09	INFO	cmd/scenario.go:27	ID: secret-tank-desant, Name: Secret Tank Desant
1.5682037678385787e+09	INFO	cmd/scenario.go:27	ID: network-siege, Name: Network Siege
1.5682037678385825e+09	INFO	cmd/scenario.go:27	ID: node-shock-tactics, Name: Node Shock Tactics
1.5682037678385894e+09	INFO	cmd/scenario.go:27	ID: network-feint, Name: Network Feint
</pre>

## How It All Works

### AWS Configuration

Simulator uses terraform to provision its infrastructure.  Terraform in turn
honours [all the rules for configuring the AWS
CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html#cli-environment).
Please follow these.

You can provide your credentials via the <code>AWS_ACCESS_KEY_ID</code> and
<code>AWS_SECRET_ACCESS_KEY</code>, environment variables, representing your AWS Access Key
and AWS Secret Key, respectively. Note that setting your AWS credentials using
either these (or legacy) environment variables will override the use of
<code>AWS_SHARED_CREDENTIALS_FILE</code> and <code>AWS_PROFILE</code>. The <code>AWS_DEFAULT_REGION</code> and
<code>AWS_SESSION_TOKEN</code> environment variables are also used, if applicable.

https://www.terraform.io/docs/backends/types/s3.html

**All the <code>AWS_*</code> configuration environment variables you have set will be propagated into the container**

#### Troubleshooting AWS

- If you get a timeout when running <code>simulator infra create</code> after about 10 minutes, the region you are using
is probably running slowly.  You must run <code>simulator infra destroy</code> and then retry <code>simulator infra
create</code>
- <code>AWS_REGION</code> vs <code>AWS_DEFAULT_REGION</code> - There have been
[some issues](https://github.com/aws/aws-sdk-go/issues/2103) with the
[Go AWS client region configuration](https://github.com/aws/aws-sdk-go#configuring-aws-region)
- [Multi-account](https://www.terraform.io/docs/backends/types/s3.html#multi-account-aws-architecture)

### SSH

Simulator whether run in the launch container or on the host machine will generate its own SSH RSA key pair.  It will
configure the cluster to allow access only with this keypair and automates writing SSH config and keyscanning the
bastion on your behalf using custom SSH config and known_hosts files.  This keeps all simulator-related SSH config
separate from any other configs you may have. All simulator-related SSH files are written to <code>~/.ssh</code> and
are files starting <code>cp_simulator_</code>

**If you delete any of the files then simulator will recreate them and reconfigure the infrastructure as necessary on the
next run**

## Development Workflow

### Git hooks

To ensure all Git hooks are in place run the following

<pre>
make setup-dev
</pre>

Development targets are specified in the [Makefile](./Makefile).

<pre>
setup-dev             Initialise simulation tree with git hooks
reset                 Clean up files left over by simulator
validate-requirements  Verify all requirements are met
previous-tag:        
release-tag:         
gpg-preflight:       
run                   Run the simulator - the build stage of the container runs all the cli tests
docker-build          Builds the launch container
docker-test           Run the tests
dep                   Install dependencies for other targets
build                 Run golang build for the CLI program
is-in-launch-container  checks you are running in the launch container
test                  run all tests except goss tests
test-acceptance       Run tcl acceptance tests for the CLI program
test-smoke            Run expect smoke test to check happy path works end-to-end
test-unit             Run golang unit tests for the CLI program
test-cleanup          cleans up automated test artefacts if e.g. you ctrl-c abort a test run
test-coverage         Run golang unit tests with coverage and opens a browser with the results
docs                  Generate documentation
release: gpg-preflight previous-tag release-tag docker-test docker-build build 
</pre>

### Git commits

We follow [the conventional commit specification](https://www.conventionalcommits.org/en/v1.0.0-beta.4/).
Please ensure your commit messages adhere to this spec.  If you submit a pull request
without following them be aware we will squash and merge your work with a message
that does.

## Architecture

### [Launching a scenario](./docs/launch.md)

### *TODO* [Validating a scenario](./docs/validation.md)

### *TODO* [Evaluating  scenario progress](./docs/evaluation.md)

### Components

* [Simulator CLI tool](./cmd) - Runs in the launch container and orchestrates everything
* [Launch container](./Dockerfile) - Isolates the scripts from the host
* [Terraform Scripts for infrastructure provisioning](./terraform) - AWS infra
* [Perturb.sh](./simulation-scripts/perturb.sh) - sets up a scenario on a cluster
* [Scenarios](./simulation-scripts/scenario)
* [Attack container](./attack) - Runs on the bastion providing all the tools needed to attack a
cluster in the given cloud provider

### Simulator API Documentation

* [Scenario](./docs/api/scenario.md)
* [Simulator](./docs/api/simulator.md)
* [Util](./docs/api/util.md)

## Project Roadmap

There is a [roadmap](./docs/roadmap.md) outlining current and planned work.
