# This file's contents to be executed on startup.

echo "Configuring AppDynamics..."

# Set env vars for agent and profiler
$env:COR_ENABLE_PROFILING = 1
$env:COR_PROFILER = "{39AEABC1-56A5-405F-B8E7-C3668490DB4A}"
$env:COR_PROFILER_PATH_64 = "$env:HOME\.appdynamics\AppDynamics.Profiler_x64.dll"

if ([System.IO.File]::Exists($env:COR_PROFILER_PATH_64)) {
	echo "Found AppDynamics profiler at $env:COR_PROFILER_PATH_64"
} else {
	echo "Error: Could not find AppDynamics profiler at $env:COR_PROFILER_PATH_64"
}

# If the config.json exists, that means the user placed it in their .appdynamics directory, in 
# which case we should use that config instead of the env vars. Otherwise, set the env settings
# and rename the config.json.default to config.json for the agent to work (even though it ignores
# the contents of config.json)
if (-not [System.IO.File]::Exists("$env:HOME\.appdynamics\AppDynamicsConfig.json")) {
    
    Rename-Item -Path "$env:HOME\.appdynamics\AppDynamicsConfig.json.default" -NewName "$env:HOME\.appdynamics\AppDynamicsConfig.json"
    echo "Using controller settings from AppDynamics Tile"

    $config = $env:VCAP_SERVICES | ConvertFrom-Json
    $credentials = $config.appdynamics.credentials

    ${env:appdynamics.controller.hostName} = $credentials.'host-name'
    ${env:appdynamics.controller.port} = $credentials.port
    ${env:appdynamics.agent.accountAccessKey} = $credentials.'account-access-key'
    ${env:appdynamics.agent.accountName} = $credentials.'account-name'
    ${env:appdynamics.controller.ssl.enabled} = $credentials.'ssl-enabled'
    
} else {
	echo "Using settings from .appdynamics/AppDynamicsConfig.json"
}

$vcap_application = $env:VCAP_APPLICATION | ConvertFrom-Json
$appname = $vcap_application.application_name
$spacename = $vcap_application.space_name

# Set default app/tier/node name if user has not specified otherwise
if (-not (Test-Path env:APPDYNAMICS_AGENT_APPLICATION_NAME)) { $env:APPDYNAMICS_AGENT_APPLICATION_NAME = $spacename }
if (-not (Test-Path env:APPDYNAMICS_AGENT_TIER_NAME)) { $env:APPDYNAMICS_AGENT_TIER_NAME = $spacename }
if (-not (Test-Path env:APPDYNAMICS_AGENT_NODE_NAME)) { $env:APPDYNAMICS_AGENT_NODE_NAME = $appname }


# Run HWC with the env vars we just set
.cloudfoundry\hwc.exe
