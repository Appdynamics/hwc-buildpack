# AppDynamics Integrated HWC Buildpack

A fork of Cloud Foundry HWC buildpack version 2.3.14 for deploying .NET full framework applications with AppDynamics monitoring.

## Using this buildpack

Cloudfoundry/PCF does not support specifying a buildpack from a git url for HWC like you can for other buildpacks so you will have to follow the instructions below for building the buildpack.

The purpose of this buildpack is to support users until PCF/CF releases a version of PAS for Windows with .profile.d support and updates the Cloudfoundry HWC buildpack to be compatible with multi-buildpack. Then we will release a "supply buildpack" that can always be used against any version of the CF HWC buildpack and this integrated buildpack can be used to support legacy versions of PAS for Windows.

### Dependencies
- [Golang Windows](https://golang.org/dl/)
- [Ginkgo](https://onsi.github.io/ginkgo/)
- [Hostable Web Core](https://github.com/cloudfoundry-incubator/hwc)

### Building the Buildpack

To build this buildpack, run the following command from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.

   ```bash
   source .envrc
   ```
   To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Install buildpack-packager

    ```bash
    ./scripts/install_tools.sh
    ```

1. Build the buildpack

    ```bash
    buildpack-packager [ --cached | --uncached ]
    ```

1. Use in Cloud Foundry

   Upload the buildpack to your Cloud Foundry and optionally specify it by name

    ```bash
    cf create-buildpack [BUILDPACK_NAME] [BUILDPACK_ZIP_FILE_PATH] 1
    cf push my_app -b [BUILDPACK_NAME] -s windows2012R2
    ```

### Changing App/Tier/Node naming

Currently by default, App/Tier/Node uses your application's Space/Space/AppName. You can change this by setting environment variables either on the command line or in Apps Manager:
APPDYNAMICS_AGENT_APPLICATION_NAME, APPDYNAMICS_AGENT_TIER_NAME, APPDYNAMICS_AGENT_NODE_NAME
```bash
cf set-env your_application_name APPDYNAMICS_AGENT_NODE_NAME your_node_name
```

### Testing

Buildpacks use the [Cutlass](https://github.com/cloudfoundry/libbuildpack/tree/master/cutlass) framework for running integration tests.

To test this buildpack, run the following command from the buildpack's directory (must be run against a Cloud Foundry deployment with Windows cells):

1. Source the .envrc file in the buildpack directory.

   ```bash
   source .envrc
   ```
   To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Run unit tests

    ```bash
    ./scripts/unit.sh
    ```

1. Run integration tests

    ```bash
    ./scripts/integration.sh
    ```

### Active Development

The project backlog for the original buildpack is on [Pivotal Tracker](https://www.pivotaltracker.com/n/projects/1042066).
